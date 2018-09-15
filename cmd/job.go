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

// CRONJOB MODE

package cmd

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/jsenon/api-cni-cleanup/config"
	k "github.com/jsenon/api-cni-cleanup/pkg/kubernetes"
	"go.opencensus.io/trace"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var urlserver string
var autodiscover bool

// serveCmd represents the serve command
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Job Cleanner",
	Long: `Launch API cleanning 
           on api-cni-cleanup server
           `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.With().Str("Service", config.Service).Logger()
		log.Logger = log.With().Str("Version", config.Version).Logger()

		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if loglevel {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			err := os.Setenv("LOGLEVEL", "debug")
			if err != nil {
				log.Error().Msgf("Error %s", err.Error())
				runtime.Goexit()
			}
		}
		log.Debug().Msg("Log level set to Debug")
		log.Debug().Msgf("Folder to watch: %s", viper.GetString("cnifiles"))
		log.Debug().Msgf("Url of cni api server: %s", viper.GetString("urlserver"))
		Job()
	},
}

func init() {
	rootCmd.AddCommand(jobCmd)
	jobCmd.Flags().StringVar(&urlserver, "urlserver", "http://myserver1:9010/cleanup,http://myserve2:9010/cleanup", "Set URL of cni api server")
	jobCmd.Flags().BoolVar(&autodiscover, "autodiscover", false, "Manage auto discovery of API CNI Server")
	err := viper.BindPFlag("urlserver", jobCmd.Flags().Lookup("urlserver"))
	if err != nil {
		log.Error().Msgf("Error binding urlserve value: %v", err.Error())
	}
	err = viper.BindPFlag("autodiscover", jobCmd.Flags().Lookup("autodiscover"))
	if err != nil {
		log.Error().Msgf("Error binding autodiscover value: %v", err.Error())
	}

}

// Job Contact CNI Server for cleanning
func Job() {
	ctx, span := trace.StartSpan(context.Background(), "(*cniserver).Job")
	defer span.End()
	switch discover := viper.GetBool("autodiscover"); discover {
	case true:
		log.Debug().Msg("Automatic discovery")
		err := urlauto(ctx)
		if err != nil {
			log.Error().Msgf("Error URL Auto func: %v", err.Error())
		}
	case false:
		log.Debug().Msg("Manual discovery")
		err := urlmanu(ctx)
		if err != nil {
			log.Error().Msgf("Error URL Manu func: %v", err.Error())
		}
	}
}

// urlauto will discover api cni cleanup server endpoint
func urlauto(ctx context.Context) error {
	_, span := trace.StartSpan(context.Background(), "(*cniserver).urlauto")
	defer span.End()
	log.Debug().Msg("In func urlauto")
	//call func PodDiscovery

	pods, err := k.PodDiscovery(ctx)
	if err != nil {
		log.Fatal().Msg("Error listing pod")
		return err
	}
	for _, n := range pods.Items {
		if strings.Contains(n.Name, "kube-flannel") {
			log.Debug().Msgf("PodName: %s", n.Name)
			log.Debug().Msgf("PodIP: %s", n.Status.PodIP)
			url := "http://" + n.Status.PodIP + ":9010" + "/cleanup"
			log.Debug().Msgf("Contact url: %s", url)
			client := &http.Client{}
			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				log.Error().Msgf("Error on NewRequest: %s", err.Error())
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Error().Msgf("Error call server %s", err.Error())
			}
			defer resp.Body.Close() //nolint: errcheck

			log.Info().Msgf("Response status: %s", resp.Status)
			log.Debug().Msgf("Response Headers: %s", resp.Header)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error().Msgf("Error read body %s", err.Error())
			}
			log.Info().Msgf("Response Body: %s", string(body))
			err = resp.Body.Close()
			if err != nil {
				log.Error().Msgf("Error close body %s", err.Error())
			}

		}
	}
	return nil
}

// urlmanu will call API CNI Server with manual url declaration
func urlmanu(ctx context.Context) error {
	_, span := trace.StartSpan(context.Background(), "(*cniserver).urlmanu")
	defer span.End()
	log.Debug().Msg("In func urlmanu")
	s := strings.Split(viper.GetString("urlserver"), ",")
	for _, n := range s {
		log.Debug().Msgf("Call server: %s", n)

		// Check if host exist before trying to contact it
		u, err := url.Parse(n)
		if err != nil {
			log.Error().Msgf("Error parsing url %s", err.Error())
		}
		log.Debug().Msgf("Extracted Hostname: %s", u.Hostname())
		_, err = net.LookupIP(u.Hostname())
		if err != nil {
			log.Error().Msgf("Error Lookup: %s", err.Error())
		} else {
			// Contact API Server
			client := &http.Client{}
			req, err := http.NewRequest("POST", n, nil)
			if err != nil {
				log.Error().Msgf("Error on NewRequest: %s", err.Error())
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Error().Msgf("Error call server %s", err.Error())
			}
			defer resp.Body.Close() //nolint: errcheck

			log.Info().Msgf("Response status: %s", resp.Status)
			log.Debug().Msgf("Response Headers: %s", resp.Header)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error().Msgf("Error read body %s", err.Error())
			}
			log.Info().Msgf("Response Body: %s", string(body))
			err = resp.Body.Close()
			if err != nil {
				log.Error().Msgf("Error close body %s", err.Error())
			}
		}
	}
	return nil
}
