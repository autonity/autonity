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

	"github.com/shirou/gopsutil/v4/cpu"

	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/p2p"
)

const (
	goroutineDumpFile = "goroutines.txt"
	cpuDumpFile       = "cpu.profile"
	memDumpFile       = "mem.profile"
	traceFile         = "trace.out"

	CPUScaleFactor       = 1.05
	MemScaleFactor       = 1.1
	GoroutineScaleFactor = 1.2
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

var DefaultMonitorConfig = Config{
	cpuThreshold:         80,
	numGoroutines:        6000,
	memThreshold:         6 * 1024 * 1024 * 1024,
	profilePerDay:        3,
	monitoringInterval:   time.Second * 60,
	cpuProfilingDuration: time.Second * 20,
	traceDuration:        time.Second * 5,
	profileDir:           "profiles",
}

type monitorService struct {
	ctx              context.Context
	cancel           context.CancelFunc
	config           *Config
	lastProfileDate  string
	profileCount     int
	wg               sync.WaitGroup
	getCPUPercent    func(interval time.Duration, perCpu bool) ([]float64, error)
	getMemUsage      func(stats *runtime.MemStats)
	getGoroutinesNum func() int
}

func New(stack *node.Node, cfg *Config) {
	cfg.profileDir = filepath.Join(stack.InstanceDir(), cfg.profileDir)
	ms := &monitorService{
		config:           cfg,
		wg:               sync.WaitGroup{},
		getCPUPercent:    cpu.Percent,
		getMemUsage:      runtime.ReadMemStats,
		getGoroutinesNum: runtime.NumGoroutine,
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

func (ms *monitorService) collectCPUDump(profileDir, postfix string) {
	// cpu profiling
	cpuDump := filepath.Join(profileDir, cpuDumpFile+postfix)
	f, err := os.Create(cpuDump)
	if err != nil {
		log.Error("Couldn't create file to write cpu profile", "error", err)
		return
	}
	defer f.Close()
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Error("Couldn't start cpu profiling", "error", err)
		return
	}
	time.Sleep(ms.config.cpuProfilingDuration)
	pprof.StopCPUProfile()
	log.Info("dumped CPU profile", "path", cpuDump)
}

func (ms *monitorService) collectHeapDump(profileDir, postfix string) {
	// mem profiling
	memDump := filepath.Join(profileDir, memDumpFile+postfix)
	f, err := os.Create(memDump)
	if err != nil {
		log.Error("Couldn't create file to write mem profile", "error", err)
		return
	}
	defer f.Close()
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		log.Error("Couldn't write mem profile", "error", err)
	}
	log.Info("dumped heap profile", "path", memDump)
}

func (ms *monitorService) collectGoRoutines(profileDir, postfix string) {
	// goroutines stack trace
	goroutines := filepath.Join(profileDir, goroutineDumpFile+postfix)
	f, err := os.Create(goroutines)
	if err != nil {
		log.Error("Couldn't create file to write goroutines", "error", err)
		return
	}
	defer f.Close()
	err = pprof.Lookup("goroutine").WriteTo(f, 2)
	if err != nil {
		log.Error("Couldn't write goroutines", "error", err)
	}
	log.Info("dumped goroutines", "path", goroutines)
}

func (ms *monitorService) collectGoTrace(profileDir, postfix string) {
	// go tracing
	traceDump := filepath.Join(profileDir, traceFile+postfix)
	f, err := os.Create(traceDump)
	if err != nil {
		log.Error("Couldn't create file to write trace", "error", err)
		return
	}
	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		log.Error("Couldn't start go trace", "error", err)
		return
	}
	time.Sleep(ms.config.traceDuration)
	trace.Stop()
	log.Info("dumped go trace", "path", traceDump)
}

func (ms *monitorService) collectDiagnostics(currentDate string) {
	profileDir := filepath.Join(ms.config.profileDir, currentDate)
	err := os.MkdirAll(profileDir, 0774)
	if err != nil && !os.IsExist(err) {
		log.Error("Error creating profile directory")
		return
	}
	postfix := "_" + strconv.Itoa(ms.profileCount+1)
	ms.collectCPUDump(profileDir, postfix)
	ms.collectHeapDump(profileDir, postfix)
	ms.collectGoRoutines(profileDir, postfix)
	ms.collectGoTrace(profileDir, postfix)
}

func (ms *monitorService) checkSystemState() {
	currentDate := time.Now().Format("2006-01-02")

	if currentDate != ms.lastProfileDate {
		ms.profileCount = 0
		ms.lastProfileDate = currentDate
	}

	if ms.profileCount >= ms.config.profilePerDay {
		return
	}

	cpuUsage, err := ms.getCPUPercent(time.Second, false)
	if err != nil {
		log.Error("fetching cpu usage", "error", err)
		return
	}

	thresholdBreach := false
	m := &runtime.MemStats{}
	ms.getMemUsage(m)
	currentMem := m.Alloc
	if currentMem > ms.config.memThreshold {
		log.Info("system memory is beyond threshold, dumping profiles", "current memmory", currentMem, "threshold", ms.config.memThreshold)
		ms.config.memThreshold = uint64(float64(currentMem) * MemScaleFactor)
		thresholdBreach = true
	}
	currentGoroutines := ms.getGoroutinesNum()
	if currentGoroutines > ms.config.numGoroutines {
		log.Info("number of running goroutines are beyond threshold, dumping profiles", "current goroutines", currentGoroutines, "threshold", ms.config.numGoroutines)
		ms.config.numGoroutines = int(float64(currentGoroutines) * GoroutineScaleFactor)
		thresholdBreach = true
	}
	currentCPU := cpuUsage[0]
	if currentCPU > ms.config.cpuThreshold {
		log.Info("system cpu usage is beyond threshold, dumping profiles", "current cpu", currentCPU, "threshold", ms.config.cpuThreshold)
		ms.config.cpuThreshold = currentCPU * CPUScaleFactor
		thresholdBreach = true
	}

	if thresholdBreach {
		ms.collectDiagnostics(currentDate)
		ms.profileCount++
	}
}
