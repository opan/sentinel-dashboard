package model

import "time"

type SentinelMaster struct {
	ID              int
	SentinelID      int    `json:"sentinel_id" db:"sentinel_id"`
	IP              string `redis:"ip"`
	Port            int    `redis:"port"`
	MasterName      string `redis:"name" json:"name" mapstructure:"name" db:"name"`
	Quorum          int    `redis:"quorum"`
	FailoverTimeout int64  `redis:"failover_timeout" json:"failover_timeout" mapstructure:"failover-timeout"`
	ParallelSyncs   int    `redis:"parallel_syncs" json:"parallel_syncs" mapstructure:"parallel-syncs"`
	NumSlaves       int    `redis:"num_slaves" json:"num_slaves" mapstructure:"num-slaves"`
	Options         *string
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
