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

package cmd

import (
	"os"
	"runtime"

	"github.com/jsenon/api-cni-cleanup/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var loglevel bool
var jaegerurl string
var url string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch VPN Server",
	Long: `Launch VPN Server 
           which manage config file generation on VPN Servers
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

		Start()
	},
}

func init() {
	serveCmd.PersistentFlags().StringVar(&url, "url", "http://localhost:8000", "url:port for dynamoDB")
	serveCmd.PersistentFlags().StringVar(&jaegerurl, "jaeger", "http://localhost:14268", "Set jaegger collector endpoint")
	serveCmd.PersistentFlags().BoolVar(&loglevel, "debug", false, "Set log level to Debug")
	rootCmd.AddCommand(serveCmd)
}

// Start the server
func Start() {
}
