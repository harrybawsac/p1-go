package models

import "time"

type Reading struct {
    ID            int64     `db:"id"`
    TimestampUTC  time.Time `db:"timestamp_utc"`
    MeterType     string    `db:"meter_type"`
    Value         float64   `db:"value"`
    Unit          string    `db:"unit"`
    Status        string    `db:"status"`
    SourceID      string    `db:"source_id"`
    CorrelationID string    `db:"correlation_id"`
}

