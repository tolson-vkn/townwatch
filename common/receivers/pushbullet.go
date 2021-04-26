package receivers

import (
	"errors"

	pushbullet "github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/sirupsen/logrus"
)

type PushBullet struct {
	client pushbullet.Pushbullet
}

func (pb *PushBullet) Notify(title, message string) error {
	note := requests.NewNote()
	note.Title = title
	note.Body = message

	// Send the note via Pushbullet.
	if _, err := pb.client.PostPushesNote(note); err != nil {
		return err
	}

	logrus.Info("Sent pushbullet notification.")
	return nil
}

func createPushBullet(config map[string]interface{}) (*PushBullet, error) {
	var exist bool

	// Probably should type check these but it should always be a string? No?
	// and probably could just use some struct syntax magic but whatever.
	if _, exist = config["access_token"].(string); exist == false {
		return nil, errors.New("Missing [access_token] in config.")
	}

	client := pushbullet.New(config["access_token"].(string))

	pbReceiver := &PushBullet{
		*client,
	}

	return pbReceiver, nil
}
