package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel string

	rootCmd = &cobra.Command{
		Use:   "townwatch",
		Short: "A log watching tool which alerts",
		Long:  `A log watching tool which alerts`,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/townwatch/townwatch.yaml)")

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(patrolCmd)
	rootCmd.AddCommand(checkCmd)
}

func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

func initConfig() {
	if cfgFile == "" {
		viper.SetConfigName("townwatch")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/townwatch/")
		viper.AddConfigPath("$HOME/.townwatch")
	} else {
		viper.SetConfigFile(cfgFile)
	}

	err := viper.ReadInConfig()
	if err != nil {
		e, ok := err.(viper.ConfigParseError)
		if ok {
			logrus.Fatalf("error parsing config file: %v", e)
		}
	} else {
		logrus.Debugf("Using config file: %v", viper.ConfigFileUsed())
	}
}
