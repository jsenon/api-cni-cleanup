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
	"github.com/jsenon/api-cni-cleanup/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var loglevel bool
var jaegerurl string
var api string
var cnifiles string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "api-cni-cleanup",
	Short: "CNI Cleanner ",
	Long: `CNI File Cleanner and Monitoring
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().
			Err(err).
			Str("service", config.Service).
			Msgf("Can't exec cmd for %s", config.Service)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&api, "api", "internal", "External or Internal K8S cluster")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.api-cni-cleanup.yaml)")
	rootCmd.PersistentFlags().StringVar(&jaegerurl, "jaegerurl", "", "Set jaegger collector endpoint")
	rootCmd.PersistentFlags().BoolVar(&loglevel, "debug", false, "Set log level to Debug")
	rootCmd.PersistentFlags().StringVar(&cnifiles, "cnifiles", "/var/lib/cni", "Set CNI Folder")
	err := viper.BindPFlag("cnifiles", rootCmd.PersistentFlags().Lookup("cnifiles"))
	if err != nil {
		log.Error().Msgf("Error binding cnifiles value: %v", err.Error())
	}
	err = viper.BindPFlag("api", rootCmd.PersistentFlags().Lookup("api"))
	if err != nil {
		log.Error().Msgf("Error binding api value: %v", err.Error())
	}
	err = viper.BindPFlag("jaegerurl", rootCmd.PersistentFlags().Lookup("jaegerurl"))
	if err != nil {
		log.Error().Msgf("Error binding jaegerurl value: %v", err.Error())
	}
	viper.SetDefault("jaegerurl", "")
	viper.SetDefault("api", "internal")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}
