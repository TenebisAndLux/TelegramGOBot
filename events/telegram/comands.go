package telegram

import (
	"TelegramGOBot/clients/telegram"
	"TelegramGOBot/lib/e"
	"TelegramGOBot/storage"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("Got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savaPage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savaPage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("Can't do command: save page", err) }()

	sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExist {
		return sendMsg(msgAlreadyExists)
		//return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(page); err != nil {
		return err
	}
	if err := sendMsg(msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("Can't do command: can't send random page", err) }()
	sendMsg := NewMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(username)
	{
		if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
			return sendMsg("An empty folder. You haven't posted any links yet")
		}
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	sendMsg := NewMessageSender(chatID, p.tg)
	return sendMsg(msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	sendMsg := NewMessageSender(chatID, p.tg)
	return sendMsg(msgHello)
}

func NewMessageSender(ChatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(ChatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
