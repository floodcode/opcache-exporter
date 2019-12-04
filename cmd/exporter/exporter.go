package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	fcgiclient "github.com/tomasen/fcgi_client"
)

const (
	namespace = "opcache"
)

var (
	enabledDesc           = newMetric("enabled", "Is OPcache enabled.")
	cacheFullDesc         = newMetric("cache_full", "Is OPcache full.")
	restartPendingDesc    = newMetric("restart_pending", "Is restart pending.")
	restartInProgressDesc = newMetric("restart_in_progress", "Is restart in progress.")

	memoryUsageUsedMemoryDesc              = newMetric("memory_usage_used_memory", "OPcache used memory.")
	memoryUsageFreeMemoryDesc              = newMetric("memory_usage_free_memory", "OPcache free memory.")
	memoryUsageWastedMemoryDesc            = newMetric("memory_usage_wasted_memory", "OPcache wasted memory.")
	memoryUsageCurrentWastedPercentageDesc = newMetric("memory_usage_current_wasted_percentage", "OPcache current wasted percentage.")

	internedStringsUsageBufferSizeDesc     = newMetric("interned_strings_usage_buffer_size", "OPcache interned string buffer size.")
	internedStringsUsageUsedMemoryDesc     = newMetric("interned_strings_usage_used_memory", "OPcache interned string used memory.")
	internedStringsUsageUsedFreeMemory     = newMetric("interned_strings_usage_free_memory", "OPcache interned string free memory.")
	internedStringsUsageUsedNumerOfStrings = newMetric("interned_strings_usage_number_of_strings", "OPcache interned string number of strings.")

	statisticsNumCachedScripts   = newMetric("statistics_num_cached_scripts", "OPcache statistics, number of cached scripts.")
	statisticsNumCachedKeys      = newMetric("statistics_num_cached_keys", "OPcache statistics, number of cached keys.")
	statisticsMaxCachedKeys      = newMetric("statistics_max_cached_keys", "OPcache statistics, max cached keys.")
	statisticsHits               = newMetric("statistics_hits", "OPcache statistics, hits.")
	statisticsStartTime          = newMetric("statistics_start_time", "OPcache statistics, start time.")
	statisticsLastRestartTime    = newMetric("statistics_last_restart_time", "OPcache statistics, last restart time")
	statisticsOOMRestarts        = newMetric("statistics_oom_restarts", "OPcache statistics, oom restarts")
	statisticsHashRestarts       = newMetric("statistics_hash_restarts", "OPcache statistics, hash restarts")
	statisticsManualRestarts     = newMetric("statistics_manual_restarts", "OPcache statistics, manual restarts")
	statisticsMisses             = newMetric("statistics_misses", "OPcache statistics, misses")
	statisticsBlacklistMisses    = newMetric("statistics_blacklist_misses", "OPcache statistics, blacklist misses")
	statisticsBlacklistMissRatio = newMetric("statistics_blacklist_miss_ratio", "OPcache statistics, blacklist miss ratio")
	statisticsHitRate            = newMetric("statistics_hit_rate", "OPcache statistics, opcache hit rate")
)

func newMetric(metricName, metricDesc string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "", metricName), metricDesc, nil, nil)
}

func boolMetric(value bool) float64 {
	return map[bool]float64{true: 1, false: 0}[value]
}

func intMetric(value int64) float64 {
	return float64(value)
}

// Exporter collects OPcache status from the given FastCGI URI and exports them using
// the prometheus metrics package.
type Exporter struct {
	mutex sync.RWMutex

	uri        string
	scriptPath string
}

// NewExporter returns an initialized Exporter.
func NewExporter(uri, scriptPath string) (*Exporter, error) {
	exporter := &Exporter{
		uri:        uri,
		scriptPath: scriptPath,
	}

	return exporter, nil
}

// Describe describes all the metrics ever exported by the OPcache exporter.
// Implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- enabledDesc
	ch <- cacheFullDesc
	ch <- restartPendingDesc
	ch <- restartInProgressDesc
	ch <- memoryUsageUsedMemoryDesc
	ch <- memoryUsageFreeMemoryDesc
	ch <- memoryUsageWastedMemoryDesc
	ch <- memoryUsageCurrentWastedPercentageDesc
	ch <- internedStringsUsageBufferSizeDesc
	ch <- internedStringsUsageUsedMemoryDesc
	ch <- internedStringsUsageUsedFreeMemory
	ch <- internedStringsUsageUsedNumerOfStrings
	ch <- statisticsNumCachedScripts
	ch <- statisticsNumCachedKeys
	ch <- statisticsMaxCachedKeys
	ch <- statisticsHits
	ch <- statisticsStartTime
	ch <- statisticsLastRestartTime
	ch <- statisticsOOMRestarts
	ch <- statisticsHashRestarts
	ch <- statisticsManualRestarts
	ch <- statisticsMisses
	ch <- statisticsBlacklistMisses
	ch <- statisticsBlacklistMissRatio
	ch <- statisticsHitRate
}

// Collect collects metrics of OPcache stats.
// Implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	status, err := e.getOpcacheStatus()
	if err != nil {
		status = new(OPcacheStatus)
	}

	ch <- prometheus.MustNewConstMetric(enabledDesc, prometheus.GaugeValue, boolMetric(status.OPcacheEnabled))
	ch <- prometheus.MustNewConstMetric(cacheFullDesc, prometheus.GaugeValue, boolMetric(status.CacheFull))
	ch <- prometheus.MustNewConstMetric(restartPendingDesc, prometheus.GaugeValue, boolMetric(status.RestartPending))
	ch <- prometheus.MustNewConstMetric(restartInProgressDesc, prometheus.GaugeValue, boolMetric(status.RestartInProgress))
	ch <- prometheus.MustNewConstMetric(memoryUsageUsedMemoryDesc, prometheus.GaugeValue, intMetric(status.MemoryUsage.UsedMemory))
	ch <- prometheus.MustNewConstMetric(memoryUsageFreeMemoryDesc, prometheus.GaugeValue, intMetric(status.MemoryUsage.FreeMemory))
	ch <- prometheus.MustNewConstMetric(memoryUsageWastedMemoryDesc, prometheus.GaugeValue, intMetric(status.MemoryUsage.WastedMemory))
	ch <- prometheus.MustNewConstMetric(memoryUsageCurrentWastedPercentageDesc, prometheus.GaugeValue, status.MemoryUsage.CurrentWastedPercentage)
	ch <- prometheus.MustNewConstMetric(internedStringsUsageBufferSizeDesc, prometheus.GaugeValue, intMetric(status.InternedStringsUsage.BufferSize))
	ch <- prometheus.MustNewConstMetric(internedStringsUsageUsedMemoryDesc, prometheus.GaugeValue, intMetric(status.InternedStringsUsage.UsedMemory))
	ch <- prometheus.MustNewConstMetric(internedStringsUsageUsedFreeMemory, prometheus.GaugeValue, intMetric(status.InternedStringsUsage.FreeMemory))
	ch <- prometheus.MustNewConstMetric(statisticsNumCachedScripts, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.NumCachedScripts))
	ch <- prometheus.MustNewConstMetric(statisticsNumCachedKeys, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.NumCachedKeys))
	ch <- prometheus.MustNewConstMetric(statisticsMaxCachedKeys, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.MaxCachedKeys))
	ch <- prometheus.MustNewConstMetric(statisticsHits, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.Hits))
	ch <- prometheus.MustNewConstMetric(statisticsStartTime, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.StartTime))
	ch <- prometheus.MustNewConstMetric(statisticsLastRestartTime, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.LastRestartTime))
	ch <- prometheus.MustNewConstMetric(statisticsOOMRestarts, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.OOMRestarts))
	ch <- prometheus.MustNewConstMetric(statisticsHashRestarts, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.HashRestarts))
	ch <- prometheus.MustNewConstMetric(statisticsManualRestarts, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.ManualRestarts))
	ch <- prometheus.MustNewConstMetric(statisticsMisses, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.Misses))
	ch <- prometheus.MustNewConstMetric(statisticsBlacklistMisses, prometheus.GaugeValue, intMetric(status.OPcacheStatistics.BlacklistMisses))
	ch <- prometheus.MustNewConstMetric(statisticsBlacklistMissRatio, prometheus.GaugeValue, status.OPcacheStatistics.BlacklistMissRatio)
	ch <- prometheus.MustNewConstMetric(statisticsHitRate, prometheus.GaugeValue, status.OPcacheStatistics.OPcacheHitRate)
}

func (e *Exporter) getOpcacheStatus() (*OPcacheStatus, error) {
	client, err := fcgiclient.Dial("tcp", e.uri)
	if err != nil {
		return nil, err
	}

	env := map[string]string{
		"SCRIPT_FILENAME": e.scriptPath,
	}

	resp, err := client.Get(env)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	status := new(OPcacheStatus)
	err = json.Unmarshal(content, status)
	if err != nil {
		return nil, errors.New(string(content))
	}

	return status, nil
}
