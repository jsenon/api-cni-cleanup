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

// Perform cleanning of cnifiles

package cleanner

import (
	"context"
	"io/ioutil"

	"github.com/rs/zerolog/log"

	k "github.com/jsenon/api-cni-cleanup/pkg/kubernetes"
	"go.opencensus.io/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var nbrfile int64

// Cleanner will clean cni folder by deleting file if pod don't exist
func Cleanner(ctx context.Context, api string, cnifiles string) error { // nolinter : gocyclo
	ctx, span := trace.StartSpan(ctx, "(*cniserver).Cleanner")
	defer span.End()

	var client *kubernetes.Clientset
	var err error

	log.Debug().Msgf("You have selected api: %s", api)

	switch api := api; api {
	case "internal":
		client, err = k.K8sInternal(ctx)
		if err != nil {
			log.Error().Msgf("Error Call client connection to k8s internal ", err.Error())
			return err
		}
	case "external":
		client, err = k.K8SExternal(ctx)
		if err != nil {
			log.Error().Msgf("Error Call client connection to k8s external ", err.Error())
			return err
		}
	default:
		log.Fatal().Msg("Error definition api type")
		return err
	}

	// List Pods interface
	log.Debug().Msgf("Debug Client:", client)
	pods, err := client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatal().Msgf("Error listing pod: %s", err.Error())
		return err
	}

	// log.Debug().Msgf("Debug", pods.Items)
	for _, n := range pods.Items {
		log.Debug().Msgf("PodName: %s", n.Name)
		log.Debug().Msgf("NodeName: %s", n.Spec.NodeName)
		log.Debug().Msgf("PodIP: %s", n.Status.PodIP)

		// TODO: Retrieve files on the node
		// Compare with n.Status.PodIP
		// Erase if file exist but n.Status.IP does not
		files, err := ioutil.ReadDir(cnifiles)
		if err != nil {
			log.Fatal().Msgf("Failed to read folder: %v", err)
			return err
		}
		nbrfile = 0
		for _, f := range files {
			// If a cni file named with ip pod exist and a pod have this IP
			if f.Name() == n.Status.PodIP {
				log.Debug().Msg("Pod Exist is running, don't erase file")
			} else {
				// A file exist but no Pod hve this IP
				log.Debug().Msg("Pod don't run, erase file")
			}
		}
	}
	return nil
}
