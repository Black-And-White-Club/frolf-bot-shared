// observability/prometheus/database/metrics.go
package databasemetrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// DatabaseMetrics defines metrics specific to database operations
type DatabaseMetrics interface {
	// Record query duration
	RecordQueryDuration(seconds float64)

	// Record query result
	RecordQueryResult(operation string, success bool)

	// Record connection pool status
	RecordConnectionPoolStatus(open, idle, used int)

	// Record transaction operations
	RecordTransaction(operation string, success bool, durationSeconds float64)

	// Record query type execution
	RecordQueryType(queryType string)
}

// databaseMetrics implements DatabaseMetrics
type databaseMetrics struct {
	queryDurationHistogram       *prometheus.Histogram
	queryResultCounter           *prometheus.CounterVec
	connectionPoolGauge          *prometheus.GaugeVec
	transactionCounter           *prometheus.CounterVec
	transactionDurationHistogram *prometheus.Histogram
	queryTypeCounter             *prometheus.CounterVec
}

// NewDatabaseMetrics creates a new DatabaseMetrics implementation
func NewDatabaseMetrics(registry *prometheus.Registry, prefix string) DatabaseMetrics {
	queryDurationHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "db",
			Name:      "query_duration_seconds",
			Help:      "Time taken to execute database queries.",
			Buckets:   prometheus.DefBuckets,
		},
	)

	queryResultCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "db",
			Name:      "query_results_total",
			Help:      "Total number of database query results, partitioned by operation and success.",
		},
		[]string{"operation", "success"},
	)

	connectionPoolGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prefix,
			Subsystem: "db",
			Name:      "connection_pool_status",
			Help:      "Current status of database connection pool.",
		},
		[]string{"state"}, // open, idle, used
	)

	transactionCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "db",
			Name:      "transactions_total",
			Help:      "Total number of database transactions, partitioned by operation and success.",
		},
		[]string{"operation", "success"},
	)

	transactionDurationHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "db",
			Name:      "transaction_duration_seconds",
			Help:      "Time taken to execute database transactions.",
			Buckets:   prometheus.DefBuckets,
		},
	)

	queryTypeCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "db",
			Name:      "query_types_total",
			Help:      "Total number of queries by type, partitioned by query type.",
		},
		[]string{"query_type"}, // select, insert, update, delete, etc.
	)

	registry.MustRegister(
		queryDurationHistogram,
		queryResultCounter,
		connectionPoolGauge,
		transactionCounter,
		transactionDurationHistogram,
		queryTypeCounter,
	)

	return &databaseMetrics{
		queryDurationHistogram:       &queryDurationHistogram,
		queryResultCounter:           queryResultCounter,
		connectionPoolGauge:          connectionPoolGauge,
		transactionCounter:           transactionCounter,
		transactionDurationHistogram: &transactionDurationHistogram,
		queryTypeCounter:             queryTypeCounter,
	}
}

// RecordQueryDuration records the time taken for a database query
func (m *databaseMetrics) RecordQueryDuration(seconds float64) {
	(*m.queryDurationHistogram).Observe(seconds)
}

// RecordQueryResult records a database query result
func (m *databaseMetrics) RecordQueryResult(operation string, success bool) {
	m.queryResultCounter.WithLabelValues(operation, fmt.Sprintf("%t", success)).Inc()
}

// RecordConnectionPoolStatus records the current status of the connection pool
func (m *databaseMetrics) RecordConnectionPoolStatus(open, idle, used int) {
	m.connectionPoolGauge.WithLabelValues("open").Set(float64(open))
	m.connectionPoolGauge.WithLabelValues("idle").Set(float64(idle))
	m.connectionPoolGauge.WithLabelValues("used").Set(float64(used))
}

// RecordTransaction records a database transaction operation
func (m *databaseMetrics) RecordTransaction(operation string, success bool, durationSeconds float64) {
	m.transactionCounter.WithLabelValues(operation, fmt.Sprintf("%t", success)).Inc()
	(*m.transactionDurationHistogram).Observe(durationSeconds)
}

// RecordQueryType records execution of a specific query type
func (m *databaseMetrics) RecordQueryType(queryType string) {
	m.queryTypeCounter.WithLabelValues(queryType).Inc()
}
