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
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"

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
		runtime.Goexit()
	}
	writeJSONResponse(w, http.StatusOK, data)
}

// Health will provide the information about state of the service.
func Health(w http.ResponseWriter, _ *http.Request) {
	data, err := json.Marshal(healthCheckResponse{Status: "UP"})
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		runtime.Goexit()
	}
	log.Debug().Msgf("Debug Marshall health", data)

	writeJSONResponse(w, http.StatusOK, data)
}

// writeJsonResponse will convert response to json
func writeJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		log.Error().Msgf("Error %s", err.Error())
		runtime.Goexit()
	}
}
