package cmd

import (
	"github.com/tolson-vkn/townwatch/common/receivers"
	"github.com/tolson-vkn/townwatch/common/watcher"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	noExamples bool
	noNotify   bool
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check config file",
	Long: `Check configuration file for errors and send
example alert.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Checking configuration.")
		r, err := receivers.InitReceiver()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("Error creating Receiver.")
		}

		if !noNotify {
			logrus.Info("Sending notify notification.")
			err := r.Notify("Townwatch Notify Test", "Daemon starting...")
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Fatal("Error in notify receiver.")
			}
		} else {
			logrus.Warn("Test notification was skipped.")
		}

		watchers, err := watcher.InitWatchers()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("Error creating Watchers.")
		}

		if !noExamples {
			for _, w := range watchers {
				if len(w.Examples) > 0 {
					err := watcher.InspectExampleLines(w.Regex, w.Examples)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"watcher": w.Name,
							"err":     err,
						}).Fatal("Error in config inspection.")
					} else {
						logrus.Infof("Watcher %s has passed inspection.", w.Name)
					}
				} else {
					logrus.Warnf("Watcher %s has no examples.", w.Name)
				}
			}
		} else {
			logrus.Warn("Example validation was skipped.")
		}
	},
}

func init() {
	// These do not not hurt me.
	checkCmd.Flags().BoolVar(&noNotify, "no-notify", false, "Don't validate receiver with a notificaiton.")
	checkCmd.Flags().BoolVar(&noExamples, "no-examples", false, "Don't validate examples in config.")
}
