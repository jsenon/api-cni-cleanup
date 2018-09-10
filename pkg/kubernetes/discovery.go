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

// Manage POD Discovery

package kubernetes

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.opencensus.io/trace"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodDiscovery will send IP of API CNI CLeanup POD
func PodDiscovery(ctx context.Context) (pods *v1.PodList, err error) {
	ctx, span := trace.StartSpan(ctx, "(*cniserver).K8SExternal")
	defer span.End()
	log.Debug().Msgf("In Poddiscovery")

	var client *kubernetes.Clientset

	// case internal or external k8s api
	api := viper.GetString("api")
	switch api := api; api {
	case "internal":
		client, err = K8sInternal(ctx)
		if err != nil {
			log.Error().Msgf("Error Call client connection to k8s internal ", err.Error())
			span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
			return nil, err
		}
	case "external":
		client, err = K8SExternal(ctx)
		if err != nil {
			log.Error().Msgf("Error Call client connection to k8s external ", err.Error())
			span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
			return nil, err
		}
	default:
		log.Fatal().Msg("Error definition api type")
		return nil, err
	}
	// List Pods interface
	pods, err = client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		log.Fatal().Msg("Error listing pod")
		return nil, err
	}
	return pods, nil
}
