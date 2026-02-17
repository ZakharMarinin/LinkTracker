package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metric struct {
	userMsgTotal *prometheus.CounterVec
}

func NewMetric() *Metric {
	userMsgTotal := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tgbot",
			Name:      "messages_sent_total",
			Help:      "Total number of messages sent to users",
		},
		[]string{"status"})

	return &Metric{
		userMsgTotal: userMsgTotal,
	}
}

func (m *Metric) MessageSent(status string) {
	m.userMsgTotal.WithLabelValues(status).Inc()
}
