package monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/stretchr/testify/require"
)

func setupService(cfg *Config) *monitorService {
	ms := &monitorService{
		config:           cfg,
		getCPUPercent:    cpu.Percent,          // Use retl CPU percent function in production
		getMemUsage:      runtime.ReadMemStats, // Use real MemStats function in production
		getGoroutinesNum: runtime.NumGoroutine,
	}
	return ms
}

// Test case for profile count limit breach
func Test_ProfileLimitBreach(t *testing.T) {
	mockCPUUsage := func(_ time.Duration, _ bool) ([]float64, error) {
		return []float64{90.0}, nil // Simulate high CPU usage
	}

	cfg := DefaultMonitorConfig
	cfg.profilePerDay = 2
	cfg.monitoringInterval = time.Second * 2
	cfg.cpuProfilingDuration = time.Second
	cfg.traceDuration = time.Second
	cfg.cpuThreshold = 20.0
	cfg.profileDir = os.TempDir() + "/profile_limit"
	mockMemUsage := func(stats *runtime.MemStats) {
		stats.Alloc = cfg.memThreshold + 10
	}
	//defer os.RemoveAll(cfg.profileDir)

	ms := setupService(&cfg)
	ms.getCPUPercent = mockCPUUsage
	ms.getMemUsage = mockMemUsage

	// Start the monitoring service
	ms.Start()

	// Sleep to let the service run diagnostics
	time.Sleep(ms.config.monitoringInterval * 7)

	require.Equal(t, cfg.profilePerDay, ms.profileCount) // Profile count should not exceed the daily limit
	for i := 0; i < ms.profileCount; i++ {
		postfix := "_" + strconv.Itoa(i+1)
		require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, cpuDumpFile+postfix), "cpu profile doesn't exist")
	}
	postfix := "_" + strconv.Itoa(ms.profileCount+1)
	require.NoFileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, cpuDumpFile+postfix), "cpu profile exists")

	ms.Stop()
}

// Test case for date change resetting profile count
func Test_DateChangeResetsProfileCount(t *testing.T) {
	mockCPUUsage := func(_ time.Duration, _ bool) ([]float64, error) {
		return []float64{75.0}, nil // Simulate high CPU usage
	}

	mockMemUsage := func(stats *runtime.MemStats) {
		stats.Alloc = 2 * 1024 * 1024 * 1024 // Memory below threshold
	}

	cfg := DefaultMonitorConfig
	cfg.monitoringInterval = time.Second * 2
	cfg.profileDir = os.TempDir() + "/profile_date"
	//defer os.RemoveAll(cfg.profileDir)
	ms := setupService(&cfg)
	ms.getCPUPercent = mockCPUUsage
	ms.getMemUsage = mockMemUsage

	// Simulate reaching the profile count limit
	ms.profileCount = cfg.profilePerDay
	ms.lastProfileDate = time.Now().Add(-24 * time.Hour).Format("2006-01-02") // Simulate date change

	// Start the monitoring service
	ms.Start()

	// Sleep to let the service run diagnostics
	time.Sleep(ms.config.monitoringInterval * 2)

	require.Equal(t, 0, ms.profileCount) // Profile count should reset after the date change

	ms.Stop()
}

// Test case stop start service handling
func Test_StopService_SystemState_NotBeingChecked(_ *testing.T) {
	cfg := DefaultMonitorConfig
	ms := setupService(&cfg)
	// Start the monitoring service
	ms.Start()
	ms.Stop()
}

// Test case for error handling in CPU and memory fetching
func Test_ErrorHandling(t *testing.T) {
	mockCPUUsage := func(_ time.Duration, _ bool) ([]float64, error) {
		return nil, fmt.Errorf("error fetching CPU usage") // Simulate error in CPU usage fetching
	}

	mockMemUsage := func(stats *runtime.MemStats) {
		// No error in memory fetching
		stats.Alloc = 2 * 1024 * 1024 * 1024
	}

	cfg := DefaultMonitorConfig
	cfg.monitoringInterval = time.Second * 2
	cfg.cpuProfilingDuration = time.Second
	cfg.traceDuration = time.Second
	cfg.profileDir = os.TempDir() + "/profile_error"
	//defer os.RemoveAll(cfg.profileDir)
	ms := setupService(&cfg)
	ms.getCPUPercent = mockCPUUsage
	ms.getMemUsage = mockMemUsage

	// Start the monitoring service
	ms.Start()

	// Sleep to let the service run diagnostics
	time.Sleep(ms.config.monitoringInterval * 2)

	// Check that profileCount is still 0 since there was an error in fetching CPU
	require.Equal(t, 0, ms.profileCount)
	postfix := "_" + strconv.Itoa(ms.profileCount+1)
	require.NoFileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, cpuDumpFile+postfix), "cpu profile exists")

	ms.Stop()
}

// Test case for CPU threshold breach
func Test_CPUThresholdBreach(t *testing.T) {
	mockCPUUsage := func(_ time.Duration, _ bool) ([]float64, error) {
		return []float64{85.0}, nil // Simulate CPU usage exceeding threshold
	}

	mockMemUsage := func(stats *runtime.MemStats) {
		stats.Alloc = 2 * 1024 * 1024 * 1024 // Memory below threshold
	}

	cfg := DefaultMonitorConfig
	cfg.monitoringInterval = time.Second * 2
	cfg.cpuProfilingDuration = time.Second
	cfg.traceDuration = time.Second
	cfg.profileDir = os.TempDir() + "/profile_cpu_resource"
	//defer os.RemoveAll(cfg.profileDir)

	ms := setupService(&cfg)
	ms.getCPUPercent = mockCPUUsage
	ms.getMemUsage = mockMemUsage

	// Start the monitoring service
	cpuThreshold, _ := mockCPUUsage(time.Second, false)
	memThreshold := ms.config.memThreshold
	grThreshold := ms.config.numGoroutines
	ms.Start()

	// Sleep to let the service collect diagnostics
	time.Sleep(ms.config.monitoringInterval * 10)

	require.Equal(t, ms.config.cpuThreshold, cpuThreshold[0]*CPUScaleFactor, "cpu threshold is not as expected") // CPU threshold should be updated
	require.Equal(t, ms.config.memThreshold, memThreshold, "mem threshold is not as expected")                   // CPU threshold should be updated
	require.Equal(t, ms.config.numGoroutines, grThreshold, "goroutine threshold is not as expected")             // CPU threshold should be updated
	postfix := "_" + strconv.Itoa(ms.profileCount)
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, cpuDumpFile+postfix), "cpu profile doesn't exist")
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, memDumpFile+postfix), "mem profile doesn't exist")
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, goroutineDumpFile+postfix), "goroutines trace doesn't exist")
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, traceFile+postfix), "go trace doesn't exist")
	require.Equal(t, 1, ms.profileCount) // Only 1 profile should be collected

	ms.Stop()
}

// test memory threshold breadh
func Test_MemoryThresholdBreach(t *testing.T) {
	mockCPUUsage := func(_ time.Duration, _ bool) ([]float64, error) {
		return []float64{79.0}, nil // Simulate CPU usage exceeding threshold
	}

	mockMemUsage := func(stats *runtime.MemStats) {
		stats.Alloc = 7 * 1024 * 1024 * 1024 // Memory below threshold
	}

	cfg := DefaultMonitorConfig
	cfg.monitoringInterval = time.Second * 2
	cfg.cpuProfilingDuration = time.Second
	cfg.traceDuration = time.Second
	cfg.profileDir = os.TempDir() + "/profile_mem_resource"
	//defer os.RemoveAll(cfg.profileDir)

	ms := setupService(&cfg)
	ms.getCPUPercent = mockCPUUsage
	ms.getMemUsage = mockMemUsage

	// Start the monitoring service
	cpuThreshold := ms.config.cpuThreshold
	m := &runtime.MemStats{}
	mockMemUsage(m)
	memThreshold := m.Alloc
	grThreshold := ms.config.numGoroutines
	ms.Start()

	// Sleep to let the service collect diagnostics
	time.Sleep(ms.config.monitoringInterval * 10)

	require.Equal(t, ms.config.cpuThreshold, cpuThreshold, "cpu threshold is not as expected")                                 // CPU threshold should be updated
	require.Equal(t, ms.config.memThreshold, uint64(float64(memThreshold)*MemScaleFactor), "mem threshold is not as expected") // CPU threshold should be updated
	require.Equal(t, ms.config.numGoroutines, grThreshold, "goroutine threshold is not as expected")                           // CPU threshold should be updated
	postfix := "_" + strconv.Itoa(ms.profileCount)
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, cpuDumpFile+postfix), "cpu profile doesn't exist")
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, memDumpFile+postfix), "mem profile doesn't exist")
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, goroutineDumpFile+postfix), "goroutines trace doesn't exist")
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, traceFile+postfix), "go trace doesn't exist")
	require.Equal(t, 1, ms.profileCount) // Only 1 profile should be collected

	ms.Stop()
}

// run goroutine threshold breach
func Test_GoroutineThresholdBreach(t *testing.T) {
	mockCPUUsage := func(_ time.Duration, _ bool) ([]float64, error) {
		return []float64{79.0}, nil // Simulate CPU usage exceeding threshold
	}

	mockMemUsage := func(stats *runtime.MemStats) {
		stats.Alloc = 2 * 1024 * 1024 * 1024 // Memory below threshold
	}

	cfg := DefaultMonitorConfig
	mockGRNum := func() int {
		return 7000
	}
	cfg.monitoringInterval = time.Second * 2
	cfg.cpuProfilingDuration = time.Second
	cfg.traceDuration = time.Second
	cfg.profileDir = os.TempDir() + "/profile_gr_resource"
	//defer os.RemoveAll(cfg.profileDir)

	ms := setupService(&cfg)
	ms.getCPUPercent = mockCPUUsage
	ms.getMemUsage = mockMemUsage
	ms.getGoroutinesNum = mockGRNum

	// Start the monitoring service
	cpuThreshold := ms.config.cpuThreshold
	memThreshold := ms.config.memThreshold
	grThreshold := mockGRNum()
	ms.Start()

	// Sleep to let the service collect diagnostics
	time.Sleep(ms.config.monitoringInterval * 10)

	require.Equal(t, ms.config.cpuThreshold, cpuThreshold, "cpu threshold is not as expected")                                          // CPU threshold should be updated
	require.Equal(t, ms.config.memThreshold, memThreshold, "mem threshold is not as expected")                                          // CPU threshold should be updated
	require.Equal(t, ms.config.numGoroutines, int(float64(grThreshold)*GoroutineScaleFactor), "goroutine threshold is not as expected") // CPU threshold should be updated
	postfix := "_" + strconv.Itoa(ms.profileCount)
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, cpuDumpFile+postfix), "cpu profile doesn't exist")
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, memDumpFile+postfix), "mem profile doesn't exist")
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, goroutineDumpFile+postfix), "goroutines trace doesn't exist")
	require.FileExists(t, filepath.Join(cfg.profileDir, ms.lastProfileDate, traceFile+postfix), "go trace doesn't exist")
	require.Equal(t, 1, ms.profileCount) // Only 1 profile should be collected

	ms.Stop()
}
