package storage

import (
	"TelegramGOBot/lib/e"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"time"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (p *Page, err error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("No saved page")

type Page struct {
	URL      string
	UserName string
	Created  time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("cat't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("cat't calculate hash", err)
	}
	return fmt.Sprint("%x", h.Sum(nil)), nil
}
