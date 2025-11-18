package pw

import (
	"errors"
	"os"

	"github.com/playwright-community/playwright-go"
)

func NewProxy() (*playwright.Proxy, error) {

	var server, username, password string
	var err error
	var ok bool

	if server, ok = os.LookupEnv("PROXY_SRVR"); !ok {
		err = errors.Join(err, errors.New("missing environment variable for server"))
	}
	if username, ok = os.LookupEnv("PROXY_USER"); !ok {
		err = errors.Join(err, errors.New("missing environment variable for username"))
	}
	if password, ok = os.LookupEnv("PROXY_PASS"); !ok {
		err = errors.Join(err, errors.New("missing environment variable for password"))
	}

	if err != nil {
		return nil, err
	}

	return &playwright.Proxy{
		Server:   server,
		Username: &username,
		Password: &password,
	}, nil
}
