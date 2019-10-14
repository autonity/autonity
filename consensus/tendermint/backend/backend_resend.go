package backend

import (
	"context"
	"errors"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"time"
)

func (sb *Backend) sendToPeer(ctx context.Context, addr common.Address, hash common.Hash, payload []byte, p consensus.Peer) chan error {
	ms, ok := sb.recentMessages.Get(addr)
	errCh := make(chan error, 1)
	var m *lru.ARCCache
	if ok {
		m, _ = ms.(*lru.ARCCache)
		if _, k := m.Get(hash); k {
			// This peer had this event, skip it
			sb.logger.Trace("inner sender loop. message sent earlier", "peer", addr.Hex(), "msg", hash.Hex())
			errCh <- nil
			return errCh
		}
	} else {
		m, _ = lru.NewARC(inmemoryMessages)
	}

	go func(ctx context.Context, p consensus.Peer, m *lru.ARCCache) {
		ticker := time.NewTicker(retryInterval * time.Millisecond)
		defer ticker.Stop()

		var err error
		var try int

	SenderLoop:
		for {
			select {
			case <-sb.stopped:
				err = peerError{errors.New("error while sending tendermintMsg message to the peer(backend is closed)"), addr}
				break SenderLoop
			case <-ticker.C:
				try++

				if err = p.Send(tendermintMsg, payload); err != nil {
					err = peerError{errors.New("error while sending tendermintMsg message to the peer: " + err.Error()), addr}

					sb.logger.Warn("inner sender loop. error",
						"try", try,
						"peer", addr.Hex(),
						"msg", hash.Hex(),
						"err", err.Error(),
						"cacheLen", m.Len(),
					)
				} else {
					err = nil

					sb.logger.Warn("inner sender loop. success",
						"try", try,
						"peer", addr.Hex(),
						"msg", hash.Hex(),
						"cacheLen", m.Len(),
					)
					break SenderLoop
				}
			case <-ctx.Done():
				err = peerError{errors.New("error while sending tendermintMsg message to the peer(context done): " + ctx.Err().Error()), addr}
				break SenderLoop
			}
		}

		if err == nil {
			m.Add(hash, true)
			sb.recentMessages.Add(addr, m)
		}

		errCh <- err
	}(ctx, p, m)

	return errCh
}

func (sb *Backend) ReSend(ctx context.Context, numberOfWorkers int) {
	wg := sync.WaitGroup{}

	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			wg.Done()
			sb.workerSendLoop(ctx)
		}(ctx)
	}

	// we want to be sure that all workers started
	wg.Wait()
}

func (sb *Backend) workerSendLoop(ctx context.Context) {
	for {
		select {
		case <-sb.stopped:
			return
		case msgToPeers := <-sb.resend:
			sb.trySend(ctx, msgToPeers)
		case <-ctx.Done():
			return
		}
	}
}

func (sb *Backend) sendToResendCh(ctx context.Context, m messageToPeers) {
	select {
	case <-sb.stopped:
		return
	case <-ctx.Done():
		return
	case sb.resend <- m:
		// fixme we'll loose message in case of stop or restart. Do we need persistence?
		//sent to channel
	}
}

func (sb *Backend) trySend(ctx context.Context, msgToPeers messageToPeers) {
	if int(time.Since(msgToPeers.startTime).Seconds()) > TTL {
		sb.logger.Trace("worker loop. messages TTL expired", "messages", msgToPeers)
		return
	}

	if !delayBeforeResendPassed(msgToPeers) {
		// send messages to the channel to further tries
		sb.sendToResendCh(ctx, msgToPeers)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, TTL*time.Second)
	defer cancel()

	errPeers := sb.sendToConnectedPeers(ctx, msgToPeers)

	if int(time.Since(msgToPeers.startTime).Seconds()) > TTL {
		sb.logger.Trace("worker loop. TTL expired", "msg", msgToPeers)
		return
	}

	if len(errPeers) > 0 {
		// send messages to the channel to further tries for error and not connected at the current time peers
		msg := messageToPeers{
			message{
				msgToPeers.msg.hash,
				msgToPeers.msg.payload,
			},
			errPeers,
			msgToPeers.startTime,
			time.Now(),
		}

		sb.sendToResendCh(ctx, msg)
	}
}

func (sb *Backend) sendToConnectedPeers(ctx context.Context, msgToPeers messageToPeers) []common.Address {
	connectedPeers := sb.getPeers(msgToPeers)
	var errConnectedPeers = make([]common.Address, 0, len(connectedPeers))

	if sb.broadcaster == nil || len(connectedPeers) == 0 {
		return errConnectedPeers
	}

	sb.logger.Trace("worker loop. resend to connected peers", "msg", msgToPeers.msg.hash.String(), "peers", peersToString(getPeerKeys(connectedPeers)))
	errChs := make([]chan error, len(connectedPeers))

	// send to connected peers and collect errors
	i := 0
	for addr, p := range connectedPeers {
		errChs[i] = sb.sendToPeer(ctx, addr, msgToPeers.msg.hash, msgToPeers.msg.payload, p)
		i++
	}

	wg := sync.WaitGroup{}
	wg.Add(len(connectedPeers))

	errConnectedCh := make(chan common.Address, len(connectedPeers))

	// collect peers that haven't received the message
	for _, errCh := range errChs {
		go func(errCh chan error) {
			err := <-errCh
			if err != nil {
				pe, ok := err.(peerError)
				if ok {
					errConnectedCh <- pe.addr

					sb.logger.Info(pe.Error(), "peer", pe.addr)
				}
			}

			close(errCh)
			wg.Done()
		}(errCh)
	}

	wg.Wait()
	close(errConnectedCh)

	for addr := range errConnectedCh {
		errConnectedPeers = append(errConnectedPeers, addr)
	}

	return errConnectedPeers
}

func (sb *Backend) getPeers(msgToPeers messageToPeers) (connectedPeers map[common.Address]consensus.Peer) {
	m := make(map[common.Address]struct{})
	for _, p := range msgToPeers.peers {
		m[p] = struct{}{}
	}

	connectedPeers = sb.broadcaster.FindPeers(m)
	return
}

func getPeerKeys(psMap map[common.Address]consensus.Peer) []common.Address {
	ps := make([]common.Address, 0, len(psMap))
	for k := range psMap {
		ps = append(ps, k)
	}
	return ps
}

type message struct {
	hash    common.Hash
	payload []byte
}

type peerError struct {
	error
	addr common.Address
}

func delayBeforeResendPassed(msgToPeers messageToPeers) bool {
	return !(time.Since(msgToPeers.lastTry).Truncate(time.Millisecond).Nanoseconds()/int64(time.Millisecond) < retryInterval &&
		time.Since(msgToPeers.startTime).Truncate(time.Millisecond).Nanoseconds()/int64(time.Millisecond) > retryInterval)
}
