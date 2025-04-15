package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Количество всех HTTP-запросов",
		},
		[]string{"path", "method", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Время выполнения HTTP-запросов",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	PickupPointsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pickup_points_created_total",
			Help: "Количество созданных ПВЗ",
		},
	)

	ReceptionsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "receptions_created_total",
			Help: "Количество созданных приёмок заказов",
		},
	)

	ProductsAdded = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "products_added_total",
			Help: "Количество добавленных товаров",
		},
	)
)

func Init() {
	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		PickupPointsCreated,
		ReceptionsCreated,
		ProductsAdded,
	)
}
