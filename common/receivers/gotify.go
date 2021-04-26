package receivers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Gotify struct {
	token string
	url   string
}

func (g *Gotify) Notify(title, message string) error {
	_, err := http.PostForm(
		g.url,
		url.Values{"message": {message}, "title": {title}},
	)
	if err != nil {
		return err
	}

	logrus.Info("Sent gotify notification.")
	return nil
}

func createGotify(config map[string]interface{}) (*Gotify, error) {
	var url string
	var ssl, exist bool
	var port int

	// Probably should type check these but it should always be a string? No?
	// and probably could just use some struct syntax magic but whatever.
	if _, exist = config["ssl"].(bool); exist == false {
		ssl = true
	}

	if config["ssl"].(bool) == false {
		logrus.Warn("SSL is not on enable with [reciever.ssl: true], continuing.")
		ssl = false
	} else {
		ssl = true
	}

	if _, exist = config["server"].(string); exist == false {
		return nil, errors.New("Missing [server] in config.")
	}

	// Port might be given as string or int so do some special things.
	if _, exist = config["port"]; exist == false {
		if ssl {
			port = 443
			logrus.Warnf("[port] not in config ssl is enabled, using port: %d.\n", port)
		} else {
			port = 80
			logrus.Warnf("[port] not in config ssl is disabled, using port: %d.\n", port)
		}
	} else {
		switch p := config["port"].(type) {
		case int:
			port = config["port"].(int)
		case string:
			var err error
			port, err = strconv.Atoi(p)
			if err != nil {
				return nil, errors.New("Couldn't parse port from [port]")
			}
		default:
			return nil, errors.New("Cannot parse type of [port].")
		}
	}

	if ssl {
		url = fmt.Sprintf("https://%s:%d/message?token=%s", config["server"], port, config["token"])
	} else {
		url = fmt.Sprintf("http://%s:%d/message?token=%s", config["server"], port, config["token"])
	}

	gotifyReceiver := &Gotify{
		config["token"].(string),
		url,
	}

	return gotifyReceiver, nil
}
