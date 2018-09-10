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

package restapi

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"

	"go.opencensus.io/tag"

	"github.com/jsenon/api-cni-cleanup/internal/cleanner"
	"github.com/spf13/viper"
	"go.opencensus.io/trace"

	"github.com/rs/zerolog/log"
)

type healthCheckResponse struct {
	Status string `json:"status"`
}

type wellknownResponse struct {
	Servicename        string `json:"Servicename"`
	Servicedescription string `json:"Servicedescription"`
	Version            string `json:"Version"`
	Versionfull        string `json:"Versionfull"`
	Revision           string `json:"Revision"`
	Branch             string `json:"Branch"`
	Builddate          string `json:"Builddate"`
	Swaggerdocurl      string `json:"Swaggerdocurl"`
	Healthzurl         string `json:"Healthzurl"`
	Metricurl          string `json:"Metricurl"`
	Endpoints          string `json:"Endpoints"`
}

// WellKnownFingerHandler will provide the information about the service.
func WellKnownFingerHandler(w http.ResponseWriter, _ *http.Request) {
	ctx, span := trace.StartSpan(context.Background(), "(*cniserver).WellKnownFingerHandler")
	span.Annotate(nil, "Received REST /.well-known")
	defer span.End()
	item := wellknownResponse{
		Servicename:        "api-cni-cleanup",
		Servicedescription: "CNI File Cleanner and Monitoring",
		Version:            "0.1",
		Versionfull:        "v.0.1",
		Revision:           "",
		Branch:             "",
		Builddate:          "",
		Swaggerdocurl:      "",
		Healthzurl:         "/healthz",
		Metricurl:          "",
		Endpoints:          ""}
	data, err := json.Marshal(item)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		runtime.Goexit()
	}
	writeJSONResponse(ctx, w, http.StatusOK, data)
}

// Health will provide the information about state of the service.
func Health(w http.ResponseWriter, _ *http.Request) {
	ctx, span := trace.StartSpan(context.Background(), "(*cniserver).Health")
	span.Annotate(nil, "Received REST /healthz")
	defer span.End()
	data, err := json.Marshal(healthCheckResponse{Status: "UP"})
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		runtime.Goexit()
	}
	log.Debug().Msgf("Debug Marshall health", data)

	writeJSONResponse(ctx, w, http.StatusOK, data)
}

// writeJsonResponse will convert response to json
func writeJSONResponse(ctx context.Context, w http.ResponseWriter, status int, data []byte) {
	_, span := trace.StartSpan(ctx, "(*cniserver).writeJSONResponse")
	span.Annotate(nil, "Write string to JSON")
	defer span.End()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		runtime.Goexit()
	}
}

// Cleanner will launch cleanning
func Cleanner(w http.ResponseWriter, _ *http.Request) {
	ctx, span := trace.StartSpan(context.Background(), "(*cniserver).Cleanner")
	log.Debug().Msg("In func Cleanner")
	span.Annotate(nil, "Received REST /cleanup")
	defer span.End()

	// Test Tag insert in Trace
	Application, err := tag.NewKey("Application")
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
	}
	Version, err := tag.NewKey("Version")
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
	}
	ctx, err = tag.New(ctx,
		tag.Insert(Application, "apicnicleanup"),
		tag.Insert(Version, "v.0.1"),
	)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
	}

	cnifiles := viper.GetString("cnifiles")
	api := viper.GetString("api")
	err = cleanner.Cleanner(ctx, api, cnifiles)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		data, errmarsh := json.Marshal(healthCheckResponse{Status: "Error"})
		if errmarsh != nil {
			log.Error().Msgf("Error %s", errmarsh.Error())
		}
		writeJSONResponse(ctx, w, http.StatusProcessing, data)
	}
	data, err := json.Marshal(healthCheckResponse{Status: "Done"})
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
	}
	writeJSONResponse(ctx, w, http.StatusOK, data)
}

// Used for debug
// func CountFile(w http.ResponseWriter, _ *http.Request) {
// 	ctx := context.Background()
// 	_, span := trace.StartSpan(ctx, "(*api).CountFile")
// 	defer span.End()

// 	nbrfiles.StatsFiles(ctx, viper.GetString("cnifiles"))

// }
