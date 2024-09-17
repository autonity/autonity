package monitor

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"

	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/p2p"
)

const (
	goroutineDumpFile = "goroutines.txt"
	cpuDumpFile       = "cpu.profile"
	memDumpFile       = "mem.profile"
	traceFile         = "trace.out"
)

type Config struct {
	cpuThreshold         float64
	numGoroutines        int
	memThreshold         uint64
	profilePerDay        int
	monitoringInterval   time.Duration
	cpuProfilingDuration time.Duration
	traceDuration        time.Duration
	profileDir           string
}

var DefaultMonitorConfig = &Config{
	cpuThreshold:         80,
	numGoroutines:        1000,
	memThreshold:         4 * 1024 * 1024,
	profilePerDay:        5,
	monitoringInterval:   time.Second * 60,
	cpuProfilingDuration: time.Second * 20,
	traceDuration:        time.Second * 5,
	profileDir:           "profiles",
}

type monitorService struct {
	ctx            context.Context
	cancel         context.CancelFunc
	config         *Config
	lastProfileDay string
	profileCount   int
	wg             sync.WaitGroup
	getCpuPercent  func(interval time.Duration, perCpu bool) ([]float64, error)
	getMemUsage    func(stats *runtime.MemStats)
}

func New(stack *node.Node, cfg *Config) {
	ms := &monitorService{
		config:        cfg,
		wg:            sync.WaitGroup{},
		getCpuPercent: cpu.Percent,
		getMemUsage:   runtime.ReadMemStats,
	}
	stack.RegisterLifecycle(ms)
}

func (ms *monitorService) Start() error {
	// setup context
	ctx, cancel := context.WithCancel(context.Background())
	ms.ctx = ctx
	ms.cancel = cancel
	ms.wg.Add(1)

	go func() {
		defer ms.wg.Done()
		for {
			select {
			case <-time.After(ms.config.monitoringInterval):
				ms.checkSystemState()
			case <-ms.ctx.Done():
				log.Info("Stopping monitoring system")
				// if any of these is running
				pprof.StopCPUProfile()
				trace.Stop()
				return
			}
		}
	}()
	return nil
}

func (ms *monitorService) Stop() error {
	ms.cancel()
	ms.wg.Wait()
	return nil
}

func (ms *monitorService) Protocols() []p2p.Protocol {
	return nil
}

func (ms *monitorService) updateThresholds() {
	//TODO: reset thresholds, when should we reset thresholds or resetting is not at all required
	ms.config.cpuThreshold = ms.config.cpuThreshold * 1.2
	ms.config.memThreshold = uint64(float64(ms.config.memThreshold) * 1.2)
	ms.config.numGoroutines = int(float64(ms.config.numGoroutines) * 1.2)
}

func (ms *monitorService) collectDiagnostics(currentDate string) {
	profileDir := filepath.Join(ms.config.profileDir, currentDate)
	err := os.MkdirAll(profileDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error("Error creating profile directory")
		return
	}
	postfix := "_" + strconv.Itoa(ms.profileCount+1)

	// cpu profiling
	cpuDump := filepath.Join(profileDir, cpuDumpFile+postfix)
	f, err := os.Create(cpuDump)
	if err != nil {
		log.Error("Couldn't create file to write cpu profile", "error", err)
		return
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		f.Close()
		log.Error("Couldn't start cpu profiling", "error", err)
		return
	}
	time.Sleep(ms.config.cpuProfilingDuration)
	pprof.StopCPUProfile()
	f.Close()

	// mem profiling
	memDump := filepath.Join(profileDir, memDumpFile+postfix)
	f, err = os.Create(memDump)
	if err != nil {
		log.Error("Couldn't create file to write mem profile", "error", err)
		return
	}
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		f.Close()
		log.Error("Couldn't write mem profile", "error", err)
		return
	}
	f.Close()

	// goroutines stack trace
	goroutines := filepath.Join(profileDir, goroutineDumpFile+postfix)
	f, err = os.Create(goroutines)
	if err != nil {
		log.Error("Couldn't create file to write goroutines", "error", err)
		return
	}
	err = pprof.Lookup("goroutine").WriteTo(f, 0)
	if err != nil {
		f.Close()
		log.Error("Couldn't write goroutines", "error", err)
		return
	}
	f.Close()

	// go tracing
	traceDump := filepath.Join(profileDir, traceFile+postfix)
	f, err = os.Create(traceDump)
	if err != nil {
		log.Error("Couldn't create file to write goroutines", "error", err)
		return
	}
	err = trace.Start(f)
	if err != nil {
		f.Close()
		log.Error("Couldn't start go trace", "error", err)
		return
	}
	time.Sleep(ms.config.traceDuration)
	trace.Stop()
	f.Close()
}

func (ms *monitorService) checkSystemState() {
	currentDate := time.Now().Format("2006-01-02")

	if currentDate != ms.lastProfileDay {
		ms.profileCount = 0
		ms.lastProfileDay = currentDate
	}

	if ms.profileCount >= ms.config.profilePerDay {
		return
	}

	cpuUsage, err := ms.getCpuPercent(time.Second, false)
	if err != nil {
		log.Error("fetching cpu usage", "error", err)
		return
	}

	m := &runtime.MemStats{}
	ms.getMemUsage(m)
	if m.Alloc > ms.config.memThreshold ||
		runtime.NumGoroutine() > ms.config.numGoroutines || cpuUsage[0] > ms.config.cpuThreshold {
		ms.collectDiagnostics(currentDate)
		ms.updateThresholds()
		ms.profileCount++
	}
}
