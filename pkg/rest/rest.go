// Copyright © 2018 Julien SENON <julien.senon@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rest

import (
	"net/http"
	"runtime"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/jsenon/api-cni-cleanup/internal/restapi"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/zpages"
)

const (
	port = ":9010"
)

// ServeRest start API Rest Server
func ServeRest() {
	log.Info().Msg("Start Rest z-Page Server")
	go func() {
		mux := http.DefaultServeMux
		zpages.Handle(mux, "/")
		log.Info().Msg("Start debuging Z-Pages Server on http://127.0.0.1:7777")
		err := http.ListenAndServe("127.0.0.1:7777", mux)
		if err != nil {
			log.Error().Msgf("Error %s", err.Error())
			runtime.Goexit()
		}
	}()
	log.Info().Msg("Start Rest Server")
	log.Info().Msg("Listening REST on port" + port)

	// API Part
	// Start Muxer
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", restapi.Health)
	mux.HandleFunc("/.well-known", restapi.WellKnownFingerHandler)

	// Metrics REST on /restmetrics
	mux.Handle("/restmetrics", promhttp.Handler())

	// Prometheus Forwarder on /metrics
	prefix := "vpncentralmanager"
	promexporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: prefix,
	})
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		runtime.Goexit()
	}
	view.RegisterExporter(promexporter)
	mux.Handle("/metrics", promexporter)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		runtime.Goexit()
	}
}
