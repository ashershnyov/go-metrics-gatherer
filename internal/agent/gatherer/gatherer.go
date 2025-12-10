package gatherer

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type Gauges map[string]float64
type Counters map[string]int64

// Gatherer gathers and stores metrics.
type Gatherer struct {
	gauges   Gauges
	counters Counters
}

// New creates a Gatherer with a RandomValue Gauge and PollCount Counter initialized.
func New() *Gatherer {
	return &Gatherer{
		gauges:   Gauges{"RandomValue": 0},
		counters: Counters{"PollCount": 0},
	}
}

// GatherAndUpdateLoop updates metrics once every time set interval of time passes.
func (g *Gatherer) GatherAndUpdateLoop(interval time.Duration) {
	pollTicker := time.NewTicker(interval)
	for {
		select {
		case <-pollTicker.C:
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				g.gatherMemStats()
				wg.Done()
			}()
			go func() {
				g.gatherPSUtilGauges()
				wg.Done()
			}()
		}
	}

}

// GetGauges returns all Gauges.
func (g *Gatherer) GetGauges() Gauges {
	return g.gauges
}

// GetGauges returns all Counters.
func (g *Gatherer) GetCounters() Counters {
	return g.counters
}

func (g *Gatherer) gatherCounters(_ *runtime.MemStats) {
	g.counters["PollCount"]++
}

func (g *Gatherer) gatherGauges(stats *runtime.MemStats) {
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
func (g *Gatherer) gatherMemStats() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	g.gatherGauges(&stats)
	g.gatherCounters(&stats)
}

// gatherPSUtilGauges gathers PSUtil gauges: Total Memory, Free Memroy and CPU Utilization.
func (g *Gatherer) gatherPSUtilGauges() {
	m, _ := mem.VirtualMemory()
	g.gauges["TotalMemory"] = float64(m.Total)
	g.gauges["FreeMemory"] = float64(m.Free)
	cpuCount, _ := cpu.Counts(false)
	cpuUsage, _ := cpu.Percent(0, true)
	for i := 0; i < cpuCount; i++ {
		utilKey := fmt.Sprintf("CPUutilization%v", i+1)
		g.gauges[utilKey] = cpuUsage[i]
	}
}
