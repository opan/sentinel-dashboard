package model

import "time"

type SentinelMaster struct {
	ID              int       `json:"id,omitempty" db:"id"`
	SentinelID      int       `json:"sentinel_id,omitempty" db:"sentinel_id"`
	IP              string    `redis:"ip" json:"ip" db:"ip"`
	Port            int       `redis:"port" json:"port" db:"port"`
	MasterName      string    `redis:"name" json:"name" mapstructure:"name" db:"name"`
	Quorum          int       `redis:"quorum" json:"quorum" db:"quorum"`
	FailoverTimeout int64     `redis:"failover_timeout" json:"failover_timeout" mapstructure:"failover-timeout"`
	ParallelSyncs   int       `redis:"parallel_syncs" json:"parallel_syncs" mapstructure:"parallel-syncs"`
	NumSlaves       int       `redis:"num_slaves" json:"num_slaves" mapstructure:"num-slaves"`
	Options         *string   `redis:"options" json:"options" db:"options"`
	CreatedAt       time.Time `json:"created_at,omitempty" db:"created_at"`
}
