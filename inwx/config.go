package inwx

import (
	"log"
	"strings"
	"time"

	"github.com/andrexus/goinwx"
	"github.com/pquerna/otp/totp"
)

type Config struct {
	Username string
	Password string
	TAN      string
	TOTPKey  string
	Sandbox  bool
}

func (c *Config) Client() (*goinwx.Client, error) {
	clientOpts := &goinwx.ClientOptions{Sandbox: c.Sandbox}
	client := goinwx.NewClient(c.Username, c.Password, clientOpts)

	defer func() {
		if err := client.Account.Logout(); err != nil {
			log.Printf("[ERROR] Failed to logout: %v", err)
		}
	}()

	log.Printf("[INFO] Trying to login with provided credentials")
	err := client.Account.Login()
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Login successful")

	if c.TAN != "" {
		log.Printf("[INFO] TAN for 2-factor auth is configured. Trying to unlock account")
		if unlockErr := client.Account.Unlock(strings.Replace(c.TAN, " ", "", -1)); unlockErr != nil {
			return nil, unlockErr
		}
	} else if c.TOTPKey != "" {
		log.Printf("[INFO] TOTP key for 2-factor auth is configured. Trying to unlock account")
		tan, totpErr := totp.GenerateCode(c.TOTPKey, time.Now())
		if totpErr != nil {
			return nil, totpErr
		}
		if unlockErr := client.Account.Unlock(tan); unlockErr != nil {
			return nil, unlockErr
		}
	}

	log.Printf("[INFO] INWX client configured for URL: %s", client.BaseURL.String())

	return client, nil
}
