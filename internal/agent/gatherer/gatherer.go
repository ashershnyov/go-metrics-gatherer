package gatherer

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type Gauges map[string]float64
type Counters map[string]int64
type fetcher func() error

// Gatherer gathers and stores metrics.
type Gatherer struct {
	gaugesMu   *sync.RWMutex
	countersMu *sync.RWMutex
	fetchers   []fetcher
	gauges     Gauges
	counters   Counters
}

// New creates a Gatherer with a RandomValue Gauge and PollCount Counter initialized.
func New() *Gatherer {
	g := &Gatherer{
		gaugesMu:   &sync.RWMutex{},
		countersMu: &sync.RWMutex{},
		gauges:     Gauges{"RandomValue": 0},
		counters:   Counters{"PollCount": 0},
	}

	g.fetchers = []fetcher{
		g.gatherMemStats,
		g.gatherPSUtilGauges,
	}

	return g
}

// GatherAndUpdateLoop updates metrics once every time set interval of time passes.
func (g *Gatherer) GatherAndUpdateLoop(interval time.Duration) {
	errChan := make(chan error, len(g.fetchers))
	pollTicker := time.NewTicker(interval)
	go func() {
		for range pollTicker.C {
			var wg sync.WaitGroup
			wg.Add(len(g.fetchers))
			for _, f := range g.fetchers {
				go func() {
					err := f()
					if err != nil {
						errChan <- err
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}()

	for i := 0; i < len(g.fetchers); i++ {
		if err := <-errChan; err != nil {
			log.Printf("Failed gathering metrics with error: %s", err.Error())
		}
	}
}

// GetGauges returns all Gauges.
func (g *Gatherer) GetGauges() Gauges {
	g.gaugesMu.RLock()
	defer g.gaugesMu.RUnlock()
	return g.gauges
}

// GetGauges returns all Counters.
func (g *Gatherer) GetCounters() Counters {
	g.countersMu.RLock()
	defer g.countersMu.RUnlock()
	return g.counters
}

func (g *Gatherer) gatherCounters(_ *runtime.MemStats) {
	g.countersMu.Lock()
	defer g.countersMu.Unlock()
	g.counters["PollCount"]++
}

func (g *Gatherer) gatherGauges(stats *runtime.MemStats) {
	g.gaugesMu.Lock()
	defer g.gaugesMu.Unlock()
	g.gauges["Alloc"] = float64(stats.Alloc)
	g.gauges["BuckHashSys"] = float64(stats.BuckHashSys)
	g.gauges["Frees"] = float64(stats.Frees)
	g.gauges["GCCPUFraction"] = stats.GCCPUFraction
	g.gauges["GCSys"] = float64(stats.GCSys)
	g.gauges["HeapAlloc"] = float64(stats.HeapAlloc)
	g.gauges["HeapIdle"] = float64(stats.HeapIdle)
	g.gauges["HeapInuse"] = float64(stats.HeapInuse)
	g.gauges["HeapObjects"] = float64(stats.HeapObjects)
	g.gauges["HeapReleased"] = float64(stats.HeapReleased)
	g.gauges["HeapSys"] = float64(stats.HeapSys)
	g.gauges["LastGC"] = float64(stats.LastGC)
	g.gauges["Lookups"] = float64(stats.Lookups)
	g.gauges["MCacheInuse"] = float64(stats.MCacheInuse)
	g.gauges["MCacheSys"] = float64(stats.MCacheSys)
	g.gauges["MSpanInuse"] = float64(stats.MSpanInuse)
	g.gauges["MSpanSys"] = float64(stats.MSpanSys)
	g.gauges["Mallocs"] = float64(stats.Mallocs)
	g.gauges["NextGC"] = float64(stats.NextGC)
	g.gauges["NumForcedGC"] = float64(stats.NumForcedGC)
	g.gauges["NumGC"] = float64(stats.NumGC)
	g.gauges["OtherSys"] = float64(stats.OtherSys)
	g.gauges["PauseTotalNs"] = float64(stats.PauseTotalNs)
	g.gauges["StackInuse"] = float64(stats.StackInuse)
	g.gauges["StackSys"] = float64(stats.StackSys)
	g.gauges["Sys"] = float64(stats.Sys)
	g.gauges["TotalAlloc"] = float64(stats.TotalAlloc)
	g.gauges["RandomValue"] = rand.Float64()
}

// gatherMemStats collects and updates runtime metrics as well as PollCount and RandomValue.
func (g *Gatherer) gatherMemStats() error {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	g.gatherGauges(&stats)
	g.gatherCounters(&stats)
	return nil
}

// gatherPSUtilGauges gathers PSUtil gauges: Total Memory, Free Memroy and CPU Utilization.
func (g *Gatherer) gatherPSUtilGauges() error {
	m, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed getting virtual memory util: %w", err)
	}
	g.gauges["TotalMemory"] = float64(m.Total)
	g.gauges["FreeMemory"] = float64(m.Free)

	cpuCount, err := cpu.Counts(false)
	if err != nil {
		return fmt.Errorf("failed getting cpu count: %w", err)
	}

	cpuUsage, err := cpu.Percent(0, true)
	if err != nil {
		return fmt.Errorf("failed getting cpu utilizations: %w", err)
	}

	g.gaugesMu.Lock()
	defer g.gaugesMu.Unlock()
	for i := 0; i < cpuCount; i++ {
		utilKey := fmt.Sprintf("CPUutilization%v", i+1)
		g.gauges[utilKey] = cpuUsage[i]
	}
	return nil
}
