package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresPool wraps pgxpool.Pool
type PostgresPool struct {
	*pgxpool.Pool
}

// PostgresConfig holds PostgreSQL connection configuration.
type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string
	Schema   string
}

// NewPostgresPool creates a new PostgreSQL connection pool.
func NewPostgresPool(cfg interface {
	GetDatabaseHost() string
	GetDatabasePort() int
	GetDatabaseDatabase() string
	GetDatabaseUser() string
	GetDatabasePassword() string
	GetDatabaseSSLMode() string
	GetDatabaseSchema() string
},
) (*PostgresPool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
		cfg.GetDatabaseUser(),
		cfg.GetDatabasePassword(),
		cfg.GetDatabaseHost(),
		cfg.GetDatabasePort(),
		cfg.GetDatabaseDatabase(),
		cfg.GetDatabaseSSLMode(),
		cfg.GetDatabaseSchema(),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &PostgresPool{pool}, nil
}

// RunMigrations runs database migrations.
func RunMigrations(pool *pgxpool.Pool) error {
	migrations := []string{
		`CREATE SCHEMA IF NOT EXISTS gateway;`,
		`CREATE TABLE IF NOT EXISTS gateway.executions (
			execution_id VARCHAR(128) PRIMARY KEY,
			workflow_id VARCHAR(128) NOT NULL,
			workflow_type VARCHAR(64) NOT NULL,
			status VARCHAR(32) NOT NULL,
			input JSONB,
			output JSONB,
			error_code VARCHAR(64),
			error_message VARCHAR(512),
			started_at BIGINT NOT NULL,
			completed_at BIGINT,
			correlation_id VARCHAR(128),
			tenant_id VARCHAR(128),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_executions_workflow_id ON gateway.executions(workflow_id);`,
		`CREATE INDEX IF NOT EXISTS idx_executions_tenant_id ON gateway.executions(tenant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_executions_status ON gateway.executions(status);`,
		`CREATE INDEX IF NOT EXISTS idx_executions_created_at ON gateway.executions(created_at);`,
		`CREATE TABLE IF NOT EXISTS gateway.idempotency_keys (
			idempotency_key VARCHAR(128) PRIMARY KEY,
			request_hash VARCHAR(64) NOT NULL,
			response_status INT NOT NULL,
			response_body BYTEA,
			response_content_type VARCHAR(128),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_idempotency_expires_at ON gateway.idempotency_keys(expires_at);`,
	}

	ctx := context.Background()
tx, err := pool.Begin(ctx)
if err != nil {
	return fmt.Errorf("begin transaction: %w", err)
}
defer func() {
	_ = tx.Rollback(ctx)
}()

	for _, migration := range migrations {
		if _, err := tx.Exec(ctx, migration); err != nil {
			return fmt.Errorf("execute migration: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// ExecutionRepository handles execution persistence.
type ExecutionRepository struct {
	pool *pgxpool.Pool
}

func NewExecutionRepository(pool *pgxpool.Pool) *ExecutionRepository {
	return &ExecutionRepository{pool: pool}
}

// ExecutionRepositoryInterface defines the interface for execution repositories.
type ExecutionRepositoryInterface interface {
	Create(ctx context.Context, record *ExecutionRecord) error
	Get(ctx context.Context, executionID string) (*ExecutionRecord, error)
	Update(ctx context.Context, record *ExecutionRecord) error
	List(ctx context.Context, workflowID, status, tenantID string, limit, offset int) ([]*ExecutionRecord, error)
}

type ExecutionRecord struct {
	ExecutionID   string
	WorkflowID    string
	WorkflowType  string
	Status        string
	Input         map[string]string
	Output        map[string]string
	ErrorCode     string
	ErrorMessage  string
	StartedAt     int64
	CompletedAt   int64
	CorrelationID string
	TenantID      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (r *ExecutionRepository) Create(ctx context.Context, record *ExecutionRecord) error {
	query := `
		INSERT INTO gateway.executions (
			execution_id, workflow_id, workflow_type, status, input, output,
			error_code, error_message, started_at, completed_at, correlation_id, tenant_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (execution_id) DO UPDATE SET
			status = EXCLUDED.status,
			output = EXCLUDED.output,
			error_code = EXCLUDED.error_code,
			error_message = EXCLUDED.error_message,
			completed_at = EXCLUDED.completed_at,
			updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query,
		record.ExecutionID,
		record.WorkflowID,
		record.WorkflowType,
		record.Status,
		record.Input,
		record.Output,
		record.ErrorCode,
		record.ErrorMessage,
		record.StartedAt,
		record.CompletedAt,
		record.CorrelationID,
		record.TenantID,
	)
	return err
}

func (r *ExecutionRepository) Get(ctx context.Context, executionID string) (*ExecutionRecord, error) {
	query := `
		SELECT execution_id, workflow_id, workflow_type, status, input, output,
			error_code, error_message, started_at, completed_at, correlation_id, tenant_id,
			created_at, updated_at
		FROM gateway.executions
		WHERE execution_id = $1
	`

	row := r.pool.QueryRow(ctx, query, executionID)

	var record ExecutionRecord
	err := row.Scan(
		&record.ExecutionID,
		&record.WorkflowID,
		&record.WorkflowType,
		&record.Status,
		&record.Input,
		&record.Output,
		&record.ErrorCode,
		&record.ErrorMessage,
		&record.StartedAt,
		&record.CompletedAt,
		&record.CorrelationID,
		&record.TenantID,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
	return nil, nil
}
	return &record, err
}

func (r *ExecutionRepository) Update(ctx context.Context, record *ExecutionRecord) error {
	query := `
		UPDATE gateway.executions SET
			status = $2,
			output = $3,
			error_code = $3,
			error_message = $4,
			completed_at = $5,
			updated_at = NOW()
		WHERE execution_id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		record.ExecutionID,
		record.Status,
		record.Output,
		record.ErrorCode,
		record.ErrorMessage,
		record.CompletedAt,
	)
	return err
}

func (r *ExecutionRepository) List(ctx context.Context, workflowID, status, tenantID string, limit, offset int) ([]*ExecutionRecord, error) {
	query := `
		SELECT execution_id, workflow_id, workflow_type, status, input, output,
			error_code, error_message, started_at, completed_at, correlation_id, tenant_id,
			created_at, updated_at
		FROM gateway.executions
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if workflowID != "" {
		query += fmt.Sprintf(" AND workflow_id = $%d", argCount)
		args = append(args, workflowID)
		argCount++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if tenantID != "" {
		query += fmt.Sprintf(" AND tenant_id = $%d", argCount)
		args = append(args, tenantID)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*ExecutionRecord
	for rows.Next() {
		var record ExecutionRecord
		err := rows.Scan(
			&record.ExecutionID,
			&record.WorkflowID,
			&record.WorkflowType,
			&record.Status,
			&record.Input,
			&record.Output,
			&record.ErrorCode,
			&record.ErrorMessage,
			&record.StartedAt,
			&record.CompletedAt,
			&record.CorrelationID,
			&record.TenantID,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, rows.Err()
}

func NewIdempotencyRepository(pool *pgxpool.Pool) IdempotencyRepository {
	return &idempotencyRepository{pool: pool}
}

type IdempotencyRepository interface {
	Get(ctx context.Context, key string) (*IdempotencyResponse, bool, error)
	Set(ctx context.Context, key, requestHash string, resp *IdempotencyResponse, ttl time.Duration) error
	TryLock(ctx context.Context, key string) (bool, error)
	Unlock(ctx context.Context, key string) error
	Cleanup(ctx context.Context) error
}

type IdempotencyResponse struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

type idempotencyRepository struct {
	pool *pgxpool.Pool
}

func (r *idempotencyRepository) Get(ctx context.Context, key string) (*IdempotencyResponse, bool, error) {
	query := `
		SELECT response_status, response_body, response_content_type
		FROM gateway.idempotency_keys
		WHERE idempotency_key = $1 AND expires_at > NOW()
	`

	row := r.pool.QueryRow(ctx, query, key)

	var resp IdempotencyResponse
	err := row.Scan(&resp.StatusCode, &resp.Body, &resp.ContentType)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &resp, true, nil
}

func (r *idempotencyRepository) Set(ctx context.Context, key, requestHash string, resp *IdempotencyResponse, ttl time.Duration) error {
	query := `
		INSERT INTO gateway.idempotency_keys (
			idempotency_key, request_hash, response_status, response_body,
			response_content_type, expires_at
		) VALUES ($1, $2, $3, $4, $5, NOW() + $6)
		ON CONFLICT (idempotency_key) DO UPDATE SET
			response_status = EXCLUDED.response_status,
			response_body = EXCLUDED.response_body,
			response_content_type = EXCLUDED.response_content_type,
			expires_at = EXCLUDED.expires_at
	`

	_, err := r.pool.Exec(ctx, query,
		key,
		requestHash,
		resp.StatusCode,
		resp.Body,
		resp.ContentType,
		ttl,
	)
	return err
}

func (r *idempotencyRepository) Cleanup(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM gateway.idempotency_keys WHERE expires_at < NOW()`)
	return err
}

func (r *idempotencyRepository) TryLock(ctx context.Context, key string) (bool, error) {
	// Use PostgreSQL advisory lock for distributed idempotency locking
	// Hash the key to a 64-bit integer for pg_advisory_xact_lock
	hash := hashKey(key)
	query := `SELECT pg_try_advisory_xact_lock($1)`
	var acquired bool
	err := r.pool.QueryRow(ctx, query, hash).Scan(&acquired)
	if err != nil {
		return false, err
	}
	return acquired, nil
}

func (r *idempotencyRepository) Unlock(ctx context.Context, key string) error {
	// Advisory transaction locks are automatically released at transaction end
	// This is a no-op for xact locks, but kept for interface consistency
	return nil
}

// hashKey converts a string key to a 64-bit integer for PostgreSQL advisory locks
func hashKey(key string) int64 {
	// Use a simple hash function - FNV-1a 64-bit
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)
	hash := uint64(offset64)
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}
	return int64(hash)
}
