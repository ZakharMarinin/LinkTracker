package metrics

import (
	"context"
	"scrapper/internal/storage"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	dbLinksTotal prometheus.Gauge
}

func NewMetrics() *Metrics {
	dbLinksTotal := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "scrapper",
		Subsystem: "db",
		Name:      "db_links_total",
		Help:      "Total number of links in database",
	})

	return &Metrics{
		dbLinksTotal: dbLinksTotal,
	}
}

func (m *Metrics) DBLinksTotal(ctx context.Context, storage *storage.PostgresStorage) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			count, err := storage.UpdateMetric(ctx)
			if err == nil {
				m.dbLinksTotal.Set(float64(count))
			}

			time.Sleep(5 * time.Minute)
		}
	}()
}
