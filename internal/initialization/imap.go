package initialization

import (
	"fmt"

	"github.com/InstaySystem/is-be/internal/config"
	"github.com/emersion/go-imap/client"
)

func InitIMAP(cfg *config.Config) (*client.Client, error) {
	imapAddr := fmt.Sprintf("%s:%d", cfg.IMAP.Host, cfg.IMAP.Port)
	c, err := client.DialTLS(imapAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("IMAP - %w", err)
	}

	if err = c.Login(cfg.IMAP.User, cfg.IMAP.Password); err != nil {
		return nil, fmt.Errorf("IMAP - %w", err)
	}

	return c, nil
}
