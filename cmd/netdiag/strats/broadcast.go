package strats

import (
	"sync"
	"time"
)

type Broadcast struct {
	Peers PeerGetter
}

func (p *Broadcast) Execute(args *ArgSizeCount, reply *ResultSimpleBroadcast) error {
	buff := make([]byte, args.Size)
	if _, err := rand.Read(buff); err != nil {
		return err
	}

	results := make([][]PacketResult, len(p.engine.peers))
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := range p.engine.peers {
		if p.engine.peers[i] == nil {
			continue
		}
		results[i] = make([]PacketResult, args.PacketCount)
		wg.Add(1)
		go func(id int) {
			resultsCh := make([]chan any, args.PacketCount)
			for j := 0; j < args.PacketCount; j++ {
				var err error
				resultsCh[j], err = p.engine.peers[id].sendDataAsync(buff)
				if err != nil {
					log.Error("error sending data async", "err", err)
				}
			}
			timer := time.NewTimer(5 * time.Second)
			for j := 0; j < args.PacketCount; j++ {
				select {
				case ans := <-resultsCh[j]:
					replyTime := ans.(AckDataPacket).Time
					results[id][j] = PacketResult{
						TimeReqReceived: time.Unix(int64(replyTime)/int64(time.Second), int64(replyTime)%int64(time.Second)),
						SyscallDuration: 0,
						Err:             "",
					}

				case <-timer.C:
					results[id][j] = PacketResult{Err: "TIMEOUT"}
					timer.Reset(5 * time.Millisecond)
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	*reply = ResultSimpleBroadcast{
		Size:          args.Size,
		Count:         args.PacketCount,
		PacketResults: results,
		StartTime:     startTime,
	}
	return nil
}
