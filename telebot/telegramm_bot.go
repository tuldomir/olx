package telebot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"olx/models"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Bot .
type Bot struct {
	client *http.Client
	url    string
	chatID int
	logger *logrus.Logger
}

// New .
func New(url string, chatid int, log *logrus.Logger) *Bot {
	client := http.DefaultClient
	return &Bot{
		client: client,
		url:    url,
		chatID: chatid,
		logger: log,
	}
}

// Listen .
func (b *Bot) Listen(ctx context.Context, aps chan []*models.Appartment) {

	for m := range aps {
		select {
		case <-ctx.Done():
			return
		default:
			str := convertMsg(m)
			fmt.Println("msg str: ", str)

			msg := &BotMsgReq{ChatID: int64(b.chatID), Text: str, ParseMode: "markdown"}
			if err := b.SendMessage(ctx, msg); err != nil {
				b.logger.Errorf("cant send message to bot %v", err)
			}
		}
	}
}

// BotMsgReq .
type BotMsgReq struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func convertMsg(aps []*models.Appartment) string {
	var builder strings.Builder
	for _, v := range aps {
		t := v.LastRefreshTime.Format(time.RFC822)
		builder.WriteString(fmt.Sprintf("%v : [URL](%v)", t, v.URL))
	}

	builder.WriteString("\n")
	return builder.String()
}

// SendMessage .
func (b *Bot) SendMessage(ctx context.Context, msg *BotMsgReq) error {

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, b.url, bytes.NewBuffer(payload))
	if err != nil {
		b.logger.Error("cant create request for bot client", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	err = b.httpDo(ctx, req, func(res *http.Response, err error) error {
		defer res.Body.Close()
		if err != nil {
			return err
		}

		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		fmt.Printf("res bot msg: %v\n", string(bs))
		return nil
	})

	return err

}

func (b *Bot) httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	errch := make(chan error)

	go func() {
		errch <- f(b.client.Do(req))
	}()

	select {
	case <-ctx.Done():
		<-errch
		return ctx.Err()
	case err := <-errch:
		return err
	}
}
