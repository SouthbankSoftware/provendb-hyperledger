/*
 * @Author: guiguan
 * @Date:   2020-08-12T15:02:07+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-13T17:51:07+10:00
 */

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/SouthbankSoftware/provendb-hyperledger/pkg/hyperledger"
	"github.com/SouthbankSoftware/provendb-tree/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

const (
	// global names

	name            = "hyperledger"
	defaultHostPort = "0.0.0.0:10016"

	// local names, default values and viper keys

	viperKeyHostPort = "host-port"
	viperKeyEnv      = "env"
	viperKeyLogLevel = "log-level"
)

var (
	// version is set automatically in CI
	cmdRoot = &cobra.Command{
		Use:   name,
		Short: "ProvenDB Hyperledger",
		RunE: func(cmd *cobra.Command, args []string) error {
			env := viper.GetString(viperKeyEnv)
			err := log.SetSharedFactory(name, env)
			if err != nil {
				return err
			}

			if l := viper.GetString(viperKeyLogLevel); l != "" {
				lvl := zapcore.DebugLevel

				err := lvl.UnmarshalText([]byte(l))
				if err != nil {
					return err
				}

				log.SF().Level().SetLevel(lvl)
			}

			return hyperledger.NewService(&hyperledger.ServiceConfig{
				HostPort: viper.GetString(viperKeyHostPort),
			}).Run()
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cmdRoot.PersistentFlags().StringP(viperKeyHostPort, "p", defaultHostPort,
		"specify the server hostPort")
	viper.BindPFlag(viperKeyHostPort, cmdRoot.PersistentFlags().Lookup(viperKeyHostPort))

	cmdRoot.PersistentFlags().StringP(viperKeyEnv, "e", "dev",
		fmt.Sprintf("specify the server running environment: %s", strings.Join(log.RunningEnvStrs, ", ")))
	viper.BindPFlag(viperKeyEnv, cmdRoot.PersistentFlags().Lookup(viperKeyEnv))

	cmdRoot.Flags().String(viperKeyLogLevel, "",
		"specify the log level. Leave empty to automatically config for \"--env\"")
	viper.BindPFlag(viperKeyLogLevel, cmdRoot.Flags().Lookup(viperKeyLogLevel))
}
