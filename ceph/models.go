package ceph

type OSD struct {
	HostName      string   `json:"host name"`
	ID            uint64   `json:"id"`
	KbAvailable   uint64   `json:"kb available"`
	KbUsed        uint64   `json:"kb used"`
	ReadByteRate  uint64   `json:"read byte rate"`
	ReadOpsRate   uint64   `json:"read ops rate"`
	State         []string `json:"state"`
	WriteByteRate uint64   `json:"write byte rate"`
	WriteOpsRate  uint64   `json:"write ops rate"`
}

type Mon struct {
	Rank       uint64 `json:"rank"`
	Name       string `json:"name"`
	Addr       string `json:"addr"`
	PublicAddr string `json:"public_addr"`
	Priority   uint64 `json:"priority"`
	Weight     uint64 `json:"weight"`
}

type Mgr struct {
	ID string
}

type Pool struct {
	PoolID             int    `json:"pool_id"`
	PoolName           string `json:"pool_name"`
	CreateTime         string `json:"create_time"`
	Size               uint64 `json:"size"`
	MinSize            uint64 `json:"min_size"`
	CrushRule          int    `json:"crush_rule"`
	PgAutoscaleMode    string `json:"pg_autoscale_mode"`
	PgNum              uint64 `json:"pg_num"`
	LastChange         string `json:"last_change"`
	QuotaMaxBytes      uint64 `json:"quota_max_bytes"`
	QuotaMaxObjects    uint64 `json:"quota_max_objects"`
	ErasureCodeProfile string `json:"erasure_code_profile"`
}
