package receivers

import (
	"github.com/sirupsen/logrus"
)

type Stdout struct{}

func (stdout *Stdout) Notify(title, message string) error {
	logrus.WithFields(logrus.Fields{
		"title":   title,
		"message": message,
	}).Info("Notify!")

	return nil
}
