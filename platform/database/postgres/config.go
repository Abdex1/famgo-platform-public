
// 2. POSTGRESQL ENTERPRISE PLATFORM

//platform/database/postgres/config.go

package postgres

import "time"

type Config struct {
    DSN                 string
    MaxConns            int32
    MinConns            int32
    MaxConnLifetime     time.Duration
    MaxConnIdleTime     time.Duration
    HealthCheckPeriod   time.Duration
    ConnectTimeout      time.Duration
    QueryTimeout        time.Duration
    EnableTracing       bool
    EnableMetrics       bool
    ApplicationName     string
    ReadReplicaDSN      string
    EnableReadReplicas  bool
}

func DefaultConfig() Config {
    return Config{
        MaxConns:           50,
        MinConns:           5,
        MaxConnLifetime:    30 * time.Minute,
        MaxConnIdleTime:    5 * time.Minute,
        HealthCheckPeriod:  30 * time.Second,
        ConnectTimeout:     15 * time.Second,
        QueryTimeout:       5 * time.Second,
        EnableTracing:      true,
        EnableMetrics:      true,
        ApplicationName:    "mobility-platform",
    }
}
