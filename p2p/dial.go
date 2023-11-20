// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package p2p

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	mrand "math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/autonity/autonity/common/mclock"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/p2p/netutil"
)

const (
	// This is the amount of time spent waiting in between redialing a certain node. The
	// limit is a bit higher than inboundThrottleTime to prevent failing dials in small
	// private networks.
	dialHistoryExpiration = inboundThrottleTime + 5*time.Second

	// Config for the "Looking for peers" message.
	dialStatsLogInterval = 10 * time.Second // printed at most this often
	dialStatsPeerLimit   = 3                // but not if more than this many dialed peers

	// Endpoint resolution is throttled with bounded backoff.
	initialResolveDelay = 60 * time.Second
	maxResolveDelay     = time.Hour
)

// NodeDialer is used to connect to nodes in the network, typically by using
// an underlying net.Dialer but also using net.Pipe in tests.
type NodeDialer interface {
	//TODO: instead of adding an extra param, think about using a separte dialer object parallel to tcpdialer
	Dial(context.Context, *enode.Node, bool) (net.Conn, error)
}

type nodeResolver interface {
	Resolve(*enode.Node) *enode.Node
}

// tcpDialer implements NodeDialer using real TCP connections.
type tcpDialer struct {
	d *net.Dialer
}

func (t tcpDialer) Dial(ctx context.Context, dest *enode.Node, dialConsensus bool) (net.Conn, error) {
	addr := nodeAddr(dest).String()
	if dialConsensus {
		addr = consensusNodeAddr(dest).String()
	}
	return t.d.DialContext(ctx, "tcp", addr)
}

func nodeAddr(n *enode.Node) net.Addr {
	return &net.TCPAddr{IP: n.IP(), Port: n.TCP()}
}

func consensusNodeAddr(n *enode.Node) net.Addr {
	return &net.TCPAddr{IP: n.CnsIP(), Port: n.CnsTCP()}
}

// checkDial errors:
var (
	errSelf             = errors.New("is self")
	errAlreadyDialing   = errors.New("already dialing")
	errAlreadyConnected = errors.New("already connected")
	errRecentlyDialed   = errors.New("recently dialed")
	errNetRestrict      = errors.New("not contained in netrestrict list")
	errNoPort           = errors.New("node does not provide TCP port")
)

// dialer creates outbound connections and submits them into Server.
// Two types of peer connections can be created:
//
//   - static dials are pre-configured connections. The dialer attempts
//     keep these nodes connected at all times.
//
//   - dynamic dials are created from node discovery results. The dialer
//     continuously reads candidate nodes from its input iterator and attempts
//     to create peer connections to nodes arriving through the iterator.
type dialScheduler struct {
	dialConfig
	setupFunc      dialSetupFunc
	wg             sync.WaitGroup
	cancel         context.CancelFunc
	ctx            context.Context
	nodesIn        chan *enode.Node
	doneCh         chan *dialTask
	addStaticCh    chan *enode.Node
	addConsensusCh chan *enode.Node
	remStaticCh    chan *enode.Node
	addPeerCh      chan *conn
	remPeerCh      chan *conn

	// Everything below here belongs to loop and
	// should only be accessed by code on the loop goroutine.
	dialing   map[string]*dialTask // active tasks
	peers     map[string]struct{}  // all connected peers
	dialPeers int                  // current number of dialed peers

	// The static map tracks all static dial tasks. The subset of usable static dial tasks
	// (i.e. those passing checkDial) is kept in staticPool. The scheduler prefers
	// launching random static tasks from the pool over launching dynamic dials from the
	// iterator.
	static     map[string]*dialTask
	staticPool []*dialTask

	// The dial history keeps recently dialed nodes. Members of history are not dialed.
	history          expHeap
	historyTimer     mclock.Timer
	historyTimerTime mclock.AbsTime

	// for logStats
	lastStatsLog     mclock.AbsTime
	doneSinceLastLog int
}

type dialSetupFunc func(net.Conn, connFlag, *enode.Node) error

type dialConfig struct {
	self           enode.ID         // our own ID
	maxDialPeers   int              // maximum number of dialed peers
	maxActiveDials int              // maximum number of active dials
	netRestrict    *netutil.Netlist // IP netrestrict list, disabled if nil
	resolver       nodeResolver
	dialer         NodeDialer
	log            log.Logger
	clock          mclock.Clock
	rand           *mrand.Rand
	trusted        *sync.Map
}

func (cfg dialConfig) withDefaults() dialConfig {
	if cfg.maxActiveDials == 0 {
		cfg.maxActiveDials = defaultMaxPendingPeers
	}
	if cfg.log == nil {
		cfg.log = log.Root()
	}
	if cfg.clock == nil {
		cfg.clock = mclock.System{}
	}
	if cfg.rand == nil {
		seedb := make([]byte, 8)
		crand.Read(seedb)
		seed := int64(binary.BigEndian.Uint64(seedb))
		cfg.rand = mrand.New(mrand.NewSource(seed))
	}
	return cfg
}

func newDialScheduler(config dialConfig, it enode.Iterator, setupFunc dialSetupFunc) *dialScheduler {
	d := &dialScheduler{
		dialConfig:  config.withDefaults(),
		setupFunc:   setupFunc,
		dialing:     make(map[string]*dialTask),
		static:      make(map[string]*dialTask),
		peers:       make(map[string]struct{}),
		doneCh:      make(chan *dialTask),
		nodesIn:     make(chan *enode.Node),
		addStaticCh: make(chan *enode.Node),
		remStaticCh: make(chan *enode.Node),
		addPeerCh:   make(chan *conn),
		remPeerCh:   make(chan *conn),
	}
	d.lastStatsLog = d.clock.Now()
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.wg.Add(2)
	go d.readNodes(it)
	go d.loop(it)
	return d
}

// stop shuts down the dialer, canceling all current dial tasks.
func (d *dialScheduler) stop() {
	d.cancel()
	d.wg.Wait()
}

// addStatic adds a static dial candidate.
func (d *dialScheduler) addConsensusChannel(n *enode.Node) {
	select {
	case d.addConsensusCh <- n:
	case <-d.ctx.Done():
	}
}

// addStatic adds a static dial candidate.
func (d *dialScheduler) addStatic(n *enode.Node) {
	select {
	case d.addStaticCh <- n:
	case <-d.ctx.Done():
	}
}

// removeStatic removes a static dial candidate.
func (d *dialScheduler) removeStatic(n *enode.Node) {
	select {
	case d.remStaticCh <- n:
	case <-d.ctx.Done():
	}
}

// peerAdded updates the peer set.
func (d *dialScheduler) peerAdded(c *conn) {
	select {
	case d.addPeerCh <- c:
	case <-d.ctx.Done():
	}
}

// peerRemoved updates the peer set.
func (d *dialScheduler) peerRemoved(c *conn) {
	select {
	case d.remPeerCh <- c:
	case <-d.ctx.Done():
	}
}

// loop is the main loop of the dialer.
func (d *dialScheduler) loop(it enode.Iterator) {
	var (
		nodesCh    chan *enode.Node
		historyExp = make(chan struct{}, 1)
	)

loop:
	for {
		// Launch new dials if slots are available.
		slots := d.freeDialSlots()
		slots -= d.startStaticDials(slots)
		if slots > 0 {
			nodesCh = d.nodesIn
		} else {
			nodesCh = nil
		}
		d.rearmHistoryTimer(historyExp)
		d.logStats()

		select {
		case node := <-nodesCh:
			flag := dynDialedConn
			for i := 0; i <= 1; i++ {
				id := node.ID().String() + ":" + strconv.Itoa(i)
				if err := d.checkDial(node, id); err != nil {
					d.log.Trace("Discarding dial candidate", "id", node.ID(), "ip", node.IP(), "reason", err)
				} else {
					d.startDial(newDialTask(node, flag))
				}
				// 2nd dial task is for consensusConn
				flag |= consensusConn
			}

		case task := <-d.doneCh:
			isConsensusConn := (task.flags >> consensusConn) & 1
			dialKey := task.dest.ID().String() + ":" + strconv.Itoa(isConsensusConn)
			delete(d.dialing, dialKey)
			d.updateStaticPool(dialKey)
			d.doneSinceLastLog++

		case c := <-d.addPeerCh:
			// Exclude trusted nodes here -
			// dialPeers is used to limit the number of dials based on maxPeers config, trusted nodes should be tried
			// regardless of the maxPeers config
			if (c.is(dynDialedConn) || c.is(staticDialedConn)) && !c.is(trustedConn) {
				d.dialPeers++
			}
			isConsensusConn := "0"
			if c.is(consensusConn) {
				isConsensusConn = "1"
			}
			id := c.node.ID().String() + ":" + isConsensusConn
			d.peers[id] = struct{}{}
			// Remove from static pool because the node is now connected.
			task := d.static[id]
			if task != nil && task.staticPoolIndex >= 0 {
				d.removeFromStaticPool(task.staticPoolIndex)
			}
			// TODO: cancel dials to connected peers

		case c := <-d.remPeerCh:
			if (c.is(dynDialedConn) || c.is(staticDialedConn)) && !c.is(trustedConn) {
				d.dialPeers--
			}
			isConsensusConn := "0"
			if c.is(consensusConn) {
				isConsensusConn = "1"
			}
			id := c.node.ID().String() + ":" + isConsensusConn

			delete(d.peers, id)
			d.updateStaticPool(id)

		case node := <-d.addConsensusCh:
			id := node.ID().String() + ":1"
			_, exists := d.static[id]
			d.log.Trace("Adding static node", "id", id, "ip", node.IP(), "added", !exists)
			if exists {
				continue loop
			}
			task := newDialTask(node, staticDialedConn|consensusConn)
			d.static[id] = task
			if d.checkDial(node, id) == nil {
				d.addToStaticPool(task)
			}

		case node := <-d.addStaticCh:
			id := node.ID().String() + ":0"
			_, exists := d.static[id]
			d.log.Trace("Adding static node", "id", id, "ip", node.IP(), "added", !exists)
			if exists {
				continue loop
			}
			task := newDialTask(node, staticDialedConn)
			d.static[id] = task
			if d.checkDial(node, id) == nil {
				d.addToStaticPool(task)
			}

		case node := <-d.remStaticCh:
			for i := 0; i <= 1; i++ {
				id := node.ID().String() + ":" + strconv.Itoa(i)
				task := d.static[id]
				d.log.Trace("Removing static node", "id", id, "ok", task != nil)
				if task != nil {
					delete(d.static, id)
					if task.staticPoolIndex >= 0 {
						//TODO: fix
						d.removeFromStaticPool(task.staticPoolIndex)
					}
				}
			}

		case <-historyExp:
			d.expireHistory()

		case <-d.ctx.Done():
			it.Close()
			break loop
		}
	}

	d.stopHistoryTimer(historyExp)
	for range d.dialing {
		<-d.doneCh
	}
	d.wg.Done()
}

// readNodes runs in its own goroutine and delivers nodes from
// the input iterator to the nodesIn channel.
func (d *dialScheduler) readNodes(it enode.Iterator) {
	defer d.wg.Done()

	for it.Next() {
		select {
		case d.nodesIn <- it.Node():
		case <-d.ctx.Done():
		}
	}
}

// logStats prints dialer statistics to the log. The message is suppressed when enough
// peers are connected because users should only see it while their client is starting up
// or comes back online.
func (d *dialScheduler) logStats() {
	now := d.clock.Now()
	if d.lastStatsLog.Add(dialStatsLogInterval) > now {
		return
	}
	if d.dialPeers < dialStatsPeerLimit && d.dialPeers < d.maxDialPeers {
		d.log.Info("Looking for peers", "peercount", len(d.peers), "tried", d.doneSinceLastLog, "static", len(d.static))
	}
	d.doneSinceLastLog = 0
	d.lastStatsLog = now
}

// rearmHistoryTimer configures d.historyTimer to fire when the
// next item in d.history expires.
func (d *dialScheduler) rearmHistoryTimer(ch chan struct{}) {
	if len(d.history) == 0 || d.historyTimerTime == d.history.nextExpiry() {
		return
	}
	d.stopHistoryTimer(ch)
	d.historyTimerTime = d.history.nextExpiry()
	timeout := time.Duration(d.historyTimerTime - d.clock.Now())
	d.historyTimer = d.clock.AfterFunc(timeout, func() { ch <- struct{}{} })
}

// stopHistoryTimer stops the timer and drains the channel it sends on.
func (d *dialScheduler) stopHistoryTimer(ch chan struct{}) {
	if d.historyTimer != nil && !d.historyTimer.Stop() {
		<-ch
	}
}

// expireHistory removes expired items from d.history.
func (d *dialScheduler) expireHistory() {
	d.historyTimer.Stop()
	d.historyTimer = nil
	d.historyTimerTime = 0
	d.history.expire(d.clock.Now(), func(hkey string) {
		d.updateStaticPool(hkey)
	})
}

// freeDialSlots returns the number of free dial slots. The result can be negative
// when peers are connected while their task is still running.
func (d *dialScheduler) freeDialSlots() int {
	slots := (d.maxDialPeers - d.dialPeers) * 2
	if slots > d.maxActiveDials {
		slots = d.maxActiveDials
	}
	free := slots - len(d.dialing)
	return free
}

// checkDial returns an error if node n should not be dialed.
func (d *dialScheduler) checkDial(n *enode.Node, connID string) error {
	if n.ID() == d.self {
		return errSelf
	}
	if n.IP() != nil && n.TCP() == 0 {
		// This check can trigger if a non-TCP node is found
		// by discovery. If there is no IP, the node is a static
		// node and the actual endpoint will be resolved later in dialTask.
		return errNoPort
	}
	//	connID := n.ID().String() + ":" + connType
	if _, ok := d.dialing[connID]; ok {
		return errAlreadyDialing
	}
	if _, ok := d.peers[connID]; ok {
		return errAlreadyConnected
	}
	if d.netRestrict != nil && !d.netRestrict.Contains(n.IP()) {
		if d.trusted == nil {
			return errNetRestrict
		} else if !isTrustedIP(d.trusted, n.IP()) {
			return errNetRestrict
		}
	}
	if d.history.contains(connID) {
		return errRecentlyDialed
	}
	return nil
}

// startStaticDials starts n static dial tasks.
func (d *dialScheduler) startStaticDials(n int) (started int) {
	for started = 0; started < n && len(d.staticPool) > 0; started++ {
		idx := d.rand.Intn(len(d.staticPool))
		task := d.staticPool[idx]
		d.startDial(task)
		d.removeFromStaticPool(idx)
	}
	return started
}

// updateStaticPool attempts to move the given static dial back into staticPool.
func (d *dialScheduler) updateStaticPool(id string) {
	if task, ok := d.static[id]; ok {
		task.lock.RLock()
		defer task.lock.RUnlock()
		if task.staticPoolIndex < 0 && d.checkDial(task.dest, id) == nil {
			d.addToStaticPool(task)
		}
	}
}

func (d *dialScheduler) addToStaticPool(task *dialTask) {
	if task.staticPoolIndex >= 0 {
		panic("attempt to add task to staticPool twice")
	}
	d.staticPool = append(d.staticPool, task)
	task.staticPoolIndex = len(d.staticPool) - 1
}

// removeFromStaticPool removes the task at idx from staticPool. It does that by moving the
// current last element of the pool to idx and then shortening the pool by one.
func (d *dialScheduler) removeFromStaticPool(idx int) {
	task := d.staticPool[idx]
	end := len(d.staticPool) - 1
	d.staticPool[idx] = d.staticPool[end]
	d.staticPool[idx].staticPoolIndex = idx
	d.staticPool[end] = nil
	d.staticPool = d.staticPool[:end]
	task.staticPoolIndex = -1
}

// startDial runs the given dial task in a separate goroutine.
func (d *dialScheduler) startDial(task *dialTask) {
	d.log.Trace("Starting p2p dial", "id", task.dest.ID(), "ip", task.dest.IP(), "flag", task.flags)
	isConsensusConn := (task.flags >> consensusConn) & 1
	hkey := task.dest.ID().String() + ":" + strconv.Itoa(isConsensusConn)
	d.history.add(hkey, d.clock.Now().Add(dialHistoryExpiration))
	d.dialing[hkey] = task
	go func() {
		task.run(d)
		d.doneCh <- task
	}()
}

// A dialTask generated for each node that is dialed.
type dialTask struct {
	staticPoolIndex int
	flags           connFlag
	// These fields are private to the task and should not be
	// accessed by dialScheduler while the task is running.
	// Unfortunately this is not the case, therefore they need protection.

	lock         sync.RWMutex
	dest         *enode.Node
	lastResolved mclock.AbsTime
	resolveDelay time.Duration
}

func newDialTask(dest *enode.Node, flags connFlag) *dialTask {
	return &dialTask{dest: dest, flags: flags, staticPoolIndex: -1}
}

type dialError struct {
	error
}

func (t *dialTask) run(d *dialScheduler) {
	//TODO: think about resolve further, is it okay to assume if a node is resolved it's cnsIP is also set
	if t.needResolve() && !t.resolve(d) {
		return
	}

	err := t.dial(d, t.dest)
	if err != nil {
		// For static nodes, resolve one more time if dialing fails.
		if _, ok := err.(*dialError); ok && t.flags&staticDialedConn != 0 {
			if t.resolve(d) {
				t.dial(d, t.dest)
			}
		}
	}
}

func (t *dialTask) needResolve() bool {
	ip := t.dest.IP()
	isConsensusConn := (t.flags >> consensusConn) & 1
	if isConsensusConn {
		ip = t.dest.CnsIP()
	}
	return t.flags&staticDialedConn != 0 && ip == nil
}

// resolve attempts to find the current endpoint for the destination
// using discovery.
//
// Resolve operations are throttled with backoff to avoid flooding the
// discovery network with useless queries for nodes that don't exist.
// The backoff delay resets when the node is found.
func (t *dialTask) resolve(d *dialScheduler) bool {
	if d.resolver == nil {
		return false
	}
	if t.resolveDelay == 0 {
		t.resolveDelay = initialResolveDelay
	}
	if t.lastResolved > 0 && time.Duration(d.clock.Now()-t.lastResolved) < t.resolveDelay {
		return false
	}
	resolved := d.resolver.Resolve(t.dest)
	t.lastResolved = d.clock.Now()
	if resolved == nil {
		t.resolveDelay *= 2
		if t.resolveDelay > maxResolveDelay {
			t.resolveDelay = maxResolveDelay
		}
		d.log.Debug("Resolving node failed", "id", t.dest.ID(), "newdelay", t.resolveDelay)
		return false
	}
	// The node was found.
	t.resolveDelay = initialResolveDelay
	t.lock.Lock()
	t.dest = resolved
	t.lock.Unlock()
	//TODO: fix log
	d.log.Debug("Resolved node", "id", t.dest.ID(), "addr", &net.TCPAddr{IP: t.dest.IP(), Port: t.dest.TCP()})
	return true
}

// dial performs the actual connection attempt.
func (t *dialTask) dial(d *dialScheduler, dest *enode.Node) error {
	isConsensusConn := (t.flags >> consensusConn) & 1
	addr := &net.TCPAddr{IP: dest.IP(), Port: dest.TCP()}
	if isConsensusConn {
		addr = &net.TCPAddr{IP: dest.CnsIP(), Port: dest.CnsTCP()}
	}
	fd, err := d.dialer.Dial(d.ctx, t.dest, isConsensusConn)
	if err != nil {
		d.log.Trace("Dial error", "id", t.dest.ID(), "addr", addr, "conn", t.flags, "err", cleanupDialErr(err))
		return &dialError{err}
	}
	mfd := newMeteredConn(fd, false, addr)
	return d.setupFunc(mfd, t.flags, dest)
}

func (t *dialTask) String() string {
	id := t.dest.ID()
	return fmt.Sprintf("%v %x %v:%d", t.flags, id[:8], t.dest.IP(), t.dest.TCP())
}

func cleanupDialErr(err error) error {
	if netErr, ok := err.(*net.OpError); ok && netErr.Op == "dial" {
		return netErr.Err
	}
	return err
}
