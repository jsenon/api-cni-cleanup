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

// Package kubernetes Manage connection to kubernetes
package kubernetes

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"go.opencensus.io/trace"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sInternal Connect to Internal k8s Cluster
func K8sInternal(ctx context.Context) (client *kubernetes.Clientset, err error) {
	_, span := trace.StartSpan(ctx, "(*cniserver).K8sInternal")
	defer span.End()
	config, err := rest.InClusterConfig()
	log.Debug().Msg("Received config object k8s")
	if err != nil {
		log.Error().Msgf("Error config in cluster api kubernetes: ", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		return nil, err
	}
	client, err = kubernetes.NewForConfig(config)
	log.Debug().Msg("Received client object k8s")
	if err != nil {
		log.Error().Msgf("Error creation clientset kubernetes: ", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		return nil, err
	}
	return client, nil
}

// K8SExternal Connect to External k8s Cluster
func K8SExternal(ctx context.Context) (client *kubernetes.Clientset, err error) {
	ctx, span := trace.StartSpan(ctx, "(*cniserver).K8SExternal")
	defer span.End()
	kubeconfig := filepath.Join(homeDir(ctx), ".kube", "config")
	log.Debug().Msgf("Kubeconfig: %s", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Error().Msgf("Error config external cluster api kubernetes: ", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		return nil, err
	}
	log.Debug().Msg("Received config object k8s")

	// create the clientset
	client, err = kubernetes.NewForConfig(config)
	log.Debug().Msg("Received config Clientset")
	if err != nil {
		log.Error().Msgf("Error creation clientset kubernetes: ", err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		return nil, err
	}
	return client, nil
}

// homeDir set home directory
func homeDir(ctx context.Context) string {
	_, span := trace.StartSpan(ctx, "(*cniserver).homeDir")
	defer span.End()
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
