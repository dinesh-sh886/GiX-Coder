package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

// TestDB holds database connection for tests
type TestDB struct {
	Pool      *pgxpool.Pool
	Config    *PostgresConfig
	Container *postgres.PostgresContainer
}

// PostgresConfig holds PostgreSQL test configuration
type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

// NewTestDB creates a new test database using Testcontainers
func NewTestDB(ctx context.Context, t *testing.T) *TestDB {
	t.Helper()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("gix_coder_test"),
		postgres.WithUsername("gix_coder"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			postgres.WaitForListeningPort("5432/tcp").
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err, "Failed to start postgres container")

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	// Run migrations
	// TODO: Run migrations here when migration tooling is ready

	return &TestDB{
		Pool: pool,
		Config: &PostgresConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "gix_coder_test",
			User:     "gix_coder",
			Password: "test",
		},
		Container: container,
	}
}

// Close closes the test database
func (db *TestDB) Close(ctx context.Context) {
	if db.Pool != nil {
		db.Pool.Close()
	}
	if db.Container != nil {
		db.Container.Terminate(ctx)
	}
}

// TestRedis holds Redis connection for tests
type TestRedis struct {
	Client    *redis.Client
	Container *redis.RedisContainer
}

// NewTestRedis creates a new test Redis using Testcontainers
func NewTestRedis(ctx context.Context, t *testing.T) *TestRedis {
	t.Helper()

	container, err := redis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			redis.WaitForListeningPort("6379/tcp").
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err, "Failed to start redis container")

	client := redis.NewClient(&redis.Options{
		Addr: container.GetHostPort(ctx, "6379/tcp"),
	})

	err := client.Ping(ctx).Err()
	require.NoError(t, err)

	return &TestRedis{
		Client:    client,
		Container: container,
	}
}

// Close closes the test Redis
func (r *TestRedis) Close(ctx context.Context) {
	if r.Client != nil {
		r.Client.Close()
	}
	if r.Container != nil {
		r.Container.Terminate(ctx)
	}
}

// TestSuite provides common test setup
type TestSuite struct {
	T      *testing.T
	Ctx    context.Context
	Cancel context.CancelFunc
	DB     *TestDB
	Redis  *TestRedis
}

// NewTestSuite creates a new test suite with all dependencies
func NewTestSuite(t *testing.T) *TestSuite {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())

	db := NewTestDB(ctx, t)
	redis := NewTestRedis(ctx, t)

	return &TestSuite{
		T:      t,
		Ctx:    ctx,
		Cancel: cancel,
		DB:     db,
		Redis:  redis,
	}
}

// Cleanup cleans up all resources
func (s *TestSuite) Cleanup() {
	s.Cancel()
	if s.DB != nil {
		s.DB.Close(s.Ctx)
	}
	if s.Redis != nil {
		s.Redis.Close(s.Ctx)
	}
}

// AssertNoError asserts that err is nil
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	require.NoError(t, err, msgAndArgs...)
}

// AssertError asserts that err is not nil
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	require.Error(t, err, msgAndArgs...)
}
