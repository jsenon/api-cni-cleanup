package jaegerexporter

import (
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

// NewExporterCollector register a new Opencensus to Jaeger exporter
func NewExporterCollector() {
	// Register the Jaeger exporter to be able to retrieve
	// the collected spans.
	addressjaeger := viper.GetString("jaegerurl")
	log.Debug().Msgf("In NewExporterCollector func, jaeger set to: %s", addressjaeger)
	exporter, err := jaeger.NewExporter(jaeger.Options{
		Endpoint:    addressjaeger,
		ServiceName: "api-cni-cleanup",
	},
	)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		runtime.Goexit()
	}
	trace.RegisterExporter(exporter)
}
