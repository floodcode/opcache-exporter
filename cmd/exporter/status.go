package main

// OPcacheStatus contains information about OPcache
type OPcacheStatus struct {
	OPcacheEnabled       bool                 `json:"opcache_enabled"`
	CacheFull            bool                 `json:"cache_full"`
	RestartPending       bool                 `json:"restart_pending"`
	RestartInProgress    bool                 `json:"restart_in_progress"`
	MemoryUsage          MemoryUsage          `json:"memory_usage"`
	InternedStringsUsage InternedStringsUsage `json:"interned_strings_usage"`
	OPcacheStatistics    OPcacheStatistics    `json:"opcache_statistics"`
}

// MemoryUsage contains information about OPcache memory usage
type MemoryUsage struct {
	UsedMemory              int64   `json:"used_memory"`
	FreeMemory              int64   `json:"free_memory"`
	WastedMemory            int64   `json:"wasted_memory"`
	CurrentWastedPercentage float64 `json:"current_wasted_percentage"`
}

// InternedStringsUsage contains information about OPcache interned strings usage
type InternedStringsUsage struct {
	BufferSize     int64 `json:"buffer_size"`
	UsedMemory     int64 `json:"used_memory"`
	FreeMemory     int64 `json:"free_memory"`
	NumerOfStrings int64 `json:"number_of_strings"`
}

// OPcacheStatistics contains information about OPcache statistics
type OPcacheStatistics struct {
	NumCachedScripts   int64   `json:"num_cached_scripts"`
	NumCachedKeys      int64   `json:"num_cached_keys"`
	MaxCachedKeys      int64   `json:"max_cached_keys"`
	Hits               int64   `json:"hits"`
	StartTime          int64   `json:"start_time"`
	LastRestartTime    int64   `json:"last_restart_time"`
	OOMRestarts        int64   `json:"oom_restarts"`
	HashRestarts       int64   `json:"hash_restarts"`
	ManualRestarts     int64   `json:"manual_restarts"`
	Misses             int64   `json:"misses"`
	BlacklistMisses    int64   `json:"blacklist_misses"`
	BlacklistMissRatio float64 `json:"blacklist_miss_ratio"`
	OPcacheHitRate     float64 `json:"opcache_hit_rate"`
}
