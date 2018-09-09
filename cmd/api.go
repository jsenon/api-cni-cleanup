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

// API Mode

package cmd

import (
	"context"
	"os"
	"runtime"

	"github.com/jsenon/api-cni-cleanup/config"
	"github.com/jsenon/api-cni-cleanup/internal/calc"
	"github.com/jsenon/api-cni-cleanup/pkg/exporter"
	"github.com/jsenon/api-cni-cleanup/pkg/rest"
	"go.opencensus.io/trace"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var apiCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch CNI Cleanner Server",
	Long: `Launch CNI Cleanner API Server 
           which manage CNI oprhane file and promtheus metrics
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
		log.Debug().Msgf("Jaeger Remote URL %s", viper.GetString("jaegerurl"))

		Start()
	},
}

func init() {

	rootCmd.AddCommand(apiCmd)
}

// Start the server
func Start() {
	ctx := context.Background()
	remotejaegurl := viper.GetString("jaegerurl")
	if remotejaegurl != "" {
		log.Debug().Msg("Jaeger endpoint has been defined")
		jaegerexporter.NewExporterCollector()
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	err := nbrfiles.StatsFiles(ctx, viper.GetString("cnifiles"))
	if err != nil {
		log.Error().Msgf("Error when retrieving stats for cni folder", err.Error())
	}
	rest.ServeRest()
}
