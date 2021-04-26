package watcher

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sync"

	"github.com/tolson-vkn/townwatch/common/receivers"

	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Watcher struct {
	Name     string   `mapstructure:"name"`
	Regex    string   `mapstructure:"regex"`
	Path     string   `mapstructure:"path"`
	Title    string   `mapstructure:"title"`
	Message  string   `mapstructure:"message"`
	Examples []string `mapstructure:"examples"`
}

func InitWatchers() ([]Watcher, error) {
	var watchers []Watcher

	v := viper.GetViper()

	err := v.UnmarshalKey("watchers", &watchers)
	if err != nil {
		logrus.Fatal(err)
	}

	// Loop over watcher values if something is value-zeroed. Error out.
	// Except examples, those are just for testing.
	for _, w := range watchers {
		structIterator := reflect.ValueOf(w)
		for i := 0; i < structIterator.NumField(); i++ {
			field := structIterator.Type().Field(i).Name

			// This is the only optional watcher config.
			if field == "Examples" {

				// Test the examples.
				err := InspectExampleLines(w.Regex, w.Examples)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Watcher %s has error. %s", w.Name, err))
				}
				continue
			}
			val := structIterator.Field(i).Interface()

			// Struct value is value-zeroed.
			if reflect.DeepEqual(val, reflect.Zero(structIterator.Field(i).Type()).Interface()) {
				return nil, errors.New(fmt.Sprintf("Missing field value for watch key: %s.", field))
			}
		}
	}

	return watchers, nil
}

// Should probably channel errors instead of log fatal.
func (watcher *Watcher) WatchAndReport(receiver receivers.Receiver, w *sync.WaitGroup) {
	defer w.Done()

	logrus.Infof("Starting watcher: %s", watcher.Name)
	logrus.WithFields(logrus.Fields{
		"name":     watcher.Name,
		"path":     watcher.Path,
		"title":    watcher.Title,
		"message":  watcher.Message,
		"examples": watcher.Examples,
		"regex":    watcher.Regex,
	}).Debug("Watcher fields.")

	regex := regexp.MustCompile(watcher.Regex)

	// Does the log exist?
	if _, err := os.Stat(watcher.Path); err != nil {
		// return errors.New(fmt.Sprintf("No log file at path: %s", watcher.Path))
		logrus.WithFields(logrus.Fields{
			"err": errors.New(fmt.Sprintf("No log file at path: %s", watcher.Path)),
		}).Fatalf("%s had an error.", watcher.Name)
	}

	seek_eof := tail.SeekInfo{
		Offset: 0,
		Whence: os.SEEK_END,
	}

	t, err := tail.TailFile(watcher.Path, tail.Config{
		Location: &seek_eof,
		ReOpen:   true,
		Follow:   true,
		Logger:   logrus.New(),
	})
	if err != nil {
		// return err
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Fatalf("%s had an error.", watcher.Name)
	}

	captures := make(map[string]string)
	for line := range t.Lines {
		captures = inspectLine(regex, line.Text)
		// Matched something
		if captures != nil {
			title, message, err := watcher.parseTemplates(captures)
			if err != nil {
				// return err
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Fatalf("%s had an error.", watcher.Name)
			}

			logrus.WithFields(logrus.Fields{
				"line": line.Text,
			}).Info("Line was captured!")

			receiver.Notify(title, message)
			if err != nil {
				// return err
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Fatalf("%s had an error.", watcher.Name)
			}
		}
	}
	return
}
