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
	"context"
	"net/http"
	"runtime"

	"github.com/rs/zerolog/log"

	"github.com/jsenon/api-cni-cleanup/internal/restapi"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
)

const (
	port = ":9010"
)

// ServeRest start API Rest Server
func ServeRest(ctx context.Context) {
	_, span := trace.StartSpan(ctx, "(*api).ServeRest")
	defer span.End()

	// ctx := context.Background()

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

	// Prometheus Forwarder on /metrics
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "apicnicleanup",
	})
	if err != nil {
		log.Fatal().Msgf("Failed to create Prometheus exporter: %v", err)
	}
	view.RegisterExporter(pe)

	// Register trace
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// API Part
	// Start Muxer
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", restapi.Health)
	mux.HandleFunc("/.well-known", restapi.WellKnownFingerHandler)
	mux.Handle("/metrics", pe)

	// Use for debuging
	// mux.HandleFunc("/file", restapi.CountFile)

	h := &ochttp.Handler{Handler: mux}
	err = view.Register(ochttp.DefaultServerViews...)
	if err != nil {
		log.Fatal().Msg("Failed to register ochttp.DefaultServerViews")
	}
	err = http.ListenAndServe(port, h)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		runtime.Goexit()
	}
}
