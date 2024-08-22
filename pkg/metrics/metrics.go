package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// DefaultRegisterer and DefaultGatherer are the implementations of the
	// prometheus Registerer and Gatherer interfaces that all metrics operations
	// will use. They are variables so that packages that embed this library can
	// replace them at runtime, instead of having to pass around specific
	// registries.
	DefaultRegisterer prometheus.Registerer = prometheus.DefaultRegisterer
	DefaultGatherer   prometheus.Gatherer   = prometheus.DefaultGatherer
	collectors        []prometheus.Collector
)

// MustRegister adds the collectors to the list of collectors to be registered.
// Deferred registration allows packages to call MustRegister in their init()
// function, while embedding libraries can set the Registerer variable in this
// package before calling Register() to actually register the collectors.
func MustRegister(cs ...prometheus.Collector) {
	collectors = append(collectors, cs...)
}

// Register calls DefaultRegisterer.MustRegister to register and clear the
// accumulated list of collectors.
func Register() {
	DefaultRegisterer.MustRegister(collectors...)
	collectors = nil
}

// Handler returns a metrics handler wired up to the current registerer and gatherer.
// Register is also called, to ensure that any queued collectors are flushed.
func Handler() http.Handler {
	Register()
	return promhttp.InstrumentMetricHandler(
		DefaultRegisterer, promhttp.HandlerFor(DefaultGatherer, promhttp.HandlerOpts{}),
	)
}
