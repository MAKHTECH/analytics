CREATE TABLE IF NOT EXISTS events (
                                      id String,
                                      timestamp DateTime,
                                      eventType LowCardinality(String),
    userId String,
    durationMs Int64,
    properties Map(String, String)
    ) ENGINE = MergeTree()
    ORDER BY timestamp;