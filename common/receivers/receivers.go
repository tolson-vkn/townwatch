package receivers

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var recieverTypes = []string{"smtp", "pushbullet", "gotify", "stdout"}

type Receiver interface {
	Notify(title, message string) error
}

// Probably should use viper unmarshel but I can't be bothered right now.
// Not quite sure how to do it with all the recievers having different k/vs
func InitReceiver() (Receiver, error) {
	var r Receiver
	var recType string
	var err error

	v := viper.GetViper()
	recMap := v.GetStringMap("receiver")

	if len(recMap) == 0 {
		logrus.Warn("Couldn't find a receiver, using [stdout]")
		recType = "stdout"
	}

	recType, ok := recMap["type"].(string)
	if !ok {
		return nil, errors.New("Cannot parse reciever type.")
	}

	var found = false
	for _, rec := range recieverTypes {
		if recMap["type"] == rec {
			found = true
		}
	}
	if !found {
		return nil, errors.New("Receiver type key does not exist.")
	}

	switch recType {
	case "smtp":
		r, err = createSMTP(recMap)
		if err != nil {
			return nil, err
		}
	case "pushbullet":
		r, err = createPushBullet(recMap)
		if err != nil {
			return nil, err
		}
	case "gotify":
		r, err = createGotify(recMap)
		if err != nil {
			return nil, err
		}
	case "stdout":
		r = &Stdout{}
	}

	return r, nil
}
