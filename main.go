package main

import (
	"context"
	"fmt"
	"olx/db"
	"olx/olxapi"
	"olx/service"

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
	//apiURL = `https://www.olx.ua/api/v1/offers?offset=0&limit=40&category_id=1757&region_id=21&city_id=121&district_id=115&owner_type=private`
)

func main() {

	var log = logrus.New()

	db, err := db.New("db.json", log)
	if err != nil {
		panic(err)
	}

	client, err := olxapi.New("proxies.txt", apiURL, log)
	if err != nil {
		panic(err)
	}
	// proxyURL, err := url.Parse("socks5://51.83.140.70:8181")
	// if err != nil {
	// 	// c.logger.Error("error ocured in olx client, parsing url", err)
	// 	panic(err)
	// }

	// transport := &http.Transport{
	// 	Proxy: http.ProxyURL(proxyURL),
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true,
	// 		SessionTicketsDisabled: true,
	// 	},
	// }

	// client := &http.Client{
	// 	Transport: transport,
	// }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	// if err != nil {
	// 	// c.logger.Error("cant reate request for olx client", err)
	// 	// return nil, err
	// 	panic(err)
	// }
	// req, err := http.NewRequest("GET", apiURL, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")
	// req.Header.Set("Accept-Encoding", "gzip, deflate")
	// req.Header.Set("Accept", "application/json")
	// req.Header.Set("Host", "www.olx.ua")

	log.Info("olx notifier started")
	serv := service.New(client, db, log, 5, 5)
	msgs := serv.Run(ctx)

	// msgs, err := client.Get(ctx)
	// res, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// var data models.QueryResult

	// bs, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(bs))
	// err = json.Unmarshal(bs, &data)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(data.Data[0].ID)

	// bot := telebot.New("https://api.telegram.org/bot5592177705:AAGJRj16pF3PP6sgAPTnWtOVXemR6vkmZoE/sendMessage", -1001527237007, log)
	// bot.Listen(ctx, msgs)

	for v := range msgs {
		fmt.Printf("new msg\n")
		// if apps != nil {
		// 	fmt.Println("error in msg, continue")
		// 	continue
		// }
		for _, m := range v {
			fmt.Printf("time: %v, url: %v\n", m.LastRefreshTime, m.URL)
		}

		fmt.Print("\n\n\n")
	}

	// bs, err := ioutil.ReadFile("res.json")
	// if err != nil {
	// 	panic(err)
	// }

	// var v models.QueryResult
	// err = json.Unmarshal(bs, &v)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(v.Data[0].ID)
}
