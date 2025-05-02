package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

// Client представляет клиент для ClickHouse
type Client struct {
	conn clickhouse.Conn
}

// Config содержит настройки подключения к ClickHouse
type Config struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Timeout  time.Duration
}

// NewClient создает новый клиент ClickHouse
func NewClient(database, username, password, host string, port int) (*Client, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	opts := &clickhouse.Options{
		Addr: []string{address},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		DialTimeout: 5 * time.Second,
	}

	conn, err := clickhouse.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к ClickHouse: %w", err)
	}

	ctx := context.Background()
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("не удалось выполнить ping ClickHouse: %w", err)
	}

	return &Client{conn: conn}, nil
}

// Close закрывает соединение с базой данных
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateTablesIfNotExist создает необходимые таблицы, если они не существуют
func (c *Client) CreateTablesIfNotExist(ctx context.Context) error {
	panic("implement me")
}
