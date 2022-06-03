package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	anyppartment = 1757
	sale         = 1758
	rentlong     = 1760
	rentshort    = 9
)

const (
	// TODO split and make dynamic url with params
	apiURL = `https://www.olx.ua/api/v1/offers?offset=0&limit=40&category_id=1757&region_id=21&city_id=121&district_id=115&owner_type=private&currency=UAH&sort_by=created_at:desc&filter_refiners=spell_checker&facets=[{%22field%22:%22district%22,%22fetchLabel%22:true,%22fetchUrl%22:true,%22limit%22:10}]`
)

func main() {

	var log = logrus.New()

	db, err := Open("db.json", log)
	if err != nil {
		panic(err)
	}

	client, err := New("socks5://72.210.252.137:4145", log)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	fmt.Println("olx scanner started")
	serv := NewServ(client, db, log)
	msgs := serv.Serve(ctx, apiURL)

	// "text" : "Click to Open [URL](http://example.com)",
	// "parse_mode" : "markdown

	// msgs := make(chan *Msgs)
	// go func() {
	// 	msgs <- &Msgs{Msgs: []*TeleMsg{{URL: "https://some url", LastRefreshTime: time.Now()}}}
	// }()

	bot := NewBot("https://api.telegram.org/bot5592177705:AAGJRj16pF3PP6sgAPTnWtOVXemR6vkmZoE/sendMessage", log)
	bot.Listen(ctx, msgs)

	// for v := range msgs {
	// 	fmt.Printf("new msg\n")
	// 	if v.Err != nil {
	// 		fmt.Println("error in msg, continue")
	// 		continue
	// 	}
	// 	for _, m := range v.Msgs {
	// 		fmt.Printf("time: %v, url: %v\n", m.LastRefreshTime, m.URL)
	// 	}

	// 	fmt.Print("\n\n\n")
	// }
}
