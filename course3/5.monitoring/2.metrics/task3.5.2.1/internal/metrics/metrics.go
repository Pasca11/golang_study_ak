package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	LoginCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "Login_counter",
		Help: "Counts /login uses",
	})
	RegisterCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "Register_counter",
		Help: "Counts /register uses",
	})

	LoginHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "Login_histogram",
		Help:    "Histogram of Login",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})
	RegisterHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "Register_histogram",
		Help:    "Histogram of Register",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})
)

func RegisterAll() {
	prometheus.MustRegister(LoginCounter)
	prometheus.MustRegister(RegisterCounter)
	prometheus.MustRegister(LoginHistogram)
	prometheus.MustRegister(RegisterHistogram)
}
