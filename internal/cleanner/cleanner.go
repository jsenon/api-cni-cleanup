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

package cleanner

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Cleanner(api string) {
	var kubeconfig *string
	var client *kubernetes.Clientset

	fmt.Println("You have selected api: ", api)
	// Internal k8s api
	if api == "internal" {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Error().Msgf("Error config in cluster api kubernetes: ", err.Error())
		}
		client, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Error().Msgf("Error creation clientset kubernetes: ", err.Error())
		}
	}
	// External k8s api based on .kube/config
	if api == "external" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		flag.Parse()
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Error().Msgf("Error config external cluster api kubernetes: ", err.Error())
		}

		// create the clientset
		client, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Error().Msgf("Error creation clientset kubernetes: ", err.Error())
		}
	}
	log.Debug().Msgf("Debug", client)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
