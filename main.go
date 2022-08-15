package main

import (
	"context"
	"fmt"
	"olx/db"
	"olx/olxapi"
	"olx/service"
	"olx/telebot"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	anyppartment = 1757
	sale         = 1758
	rentlong     = 1760
	rentshort    = 9
)

const (
	// olxAPIURL = `https://www.olx.ua/api/v1/offers?offset=0&limit=40&category_id=1757&region_id=21&city_id=121&district_id=115&owner_type=private&currency=UAH&sort_by=created_at:desc&filter_refiners=spell_checker&facets=[{%22field%22:%22district%22,%22fetchLabel%22:true,%22fetchUrl%22:true,%22limit%22:10}]`

	olxAPIURL = "https://www.olx.ua/api/v1"
	// olx params
	offset     = 0
	limit      = 40
	categoryID = 1757
	regionID   = 21
	cityID     = 121
	districtID = 115

	telegrammAPIURL = "https://api.telegram.org"
)

func main() {
	botURL := os.Getenv("BOT_URL")
	chatIDStr := os.Getenv("BOT_ID")

	if botURL == "" || chatIDStr == "" {
		panic("bot url and bot id not specified in env")
	}

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		panic(err)
	}

	olxURL := fmt.Sprintf("%v/offers?offset=%v&limit=%v&category_id=%v&region_id=%v&city_id=%v&district_id=%v&owner_type=private&currency=USD&sort_by=created_at:desc&filter_refiners=spell",
		olxAPIURL, offset, limit, categoryID, regionID, cityID, districtID)

	var log = logrus.New()

	db, err := db.New("db.json", log)
	if err != nil {
		panic(err)
	}

	client, err := olxapi.New("proxies.txt", log)
	if err != nil {
		panic(err)
	}

	poller := service.New(client, 5, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("olx notifier started")
	serv := service.NewServer(poller, db, olxURL, log)
	msgs := serv.Run(ctx)

	bot := telebot.New(fmt.Sprintf("%v/%v/sendMessage", telegrammAPIURL, botURL), chatID, log)
	bot.Listen(ctx, msgs)
}
