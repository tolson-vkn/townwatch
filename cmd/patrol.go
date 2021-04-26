package cmd

import (
	"sync"

	"github.com/tolson-vkn/townwatch/common/receivers"
	"github.com/tolson-vkn/townwatch/common/watcher"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	noValidate bool
	noStartup  bool
)

var patrolCmd = &cobra.Command{
	Use:   "patrol",
	Short: "Start log watch server",
	Long: `Start the server up and watch logs
defined in the configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("And now my watch begins.")
		r, err := receivers.InitReceiver()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("Error creating Receiver.")
		}

		watchers, err := watcher.InitWatchers()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("Error creating Watchers.")
		}

		if !noValidate {
			for _, w := range watchers {
				if len(w.Examples) > 0 {
					err := watcher.InspectExampleLines(w.Regex, w.Examples)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"watcher": w.Name,
							"err":     err,
						}).Fatal("Error in config inspection.")
					}
				} else {
					logrus.Infof("Watcher %s has no examples.", w.Name)
				}
			}
		}

		if !noStartup {
			logrus.Info("Sending startup notification.")
			err := r.Notify("Townwatch Startup", "Daemon starting...")
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Fatal("Error in startup receiver.")
			}
		}

		for {
			w := &sync.WaitGroup{}
			for i := 0; i < len(watchers); i++ {
				w.Add(1)
				go watchers[i].WatchAndReport(r, w)
			}
			w.Wait()
		}

	},
}

func init() {
	// This does not not hurt me.
	patrolCmd.Flags().BoolVar(&noStartup, "no-startup", false, "Startup notification.")
	patrolCmd.Flags().BoolVar(&noValidate, "no-validate", false, "Don't validate examples in config.")
}
