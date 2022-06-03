package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

// OlxClient .
type OlxClient struct {
	client *http.Client
	logger *logrus.Logger
}

// New .
func New(proxy string, log *logrus.Logger) (*OlxClient, error) {
	log.Info("initiating new olx client, proxy:", proxy)
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Error("error ocured in olx client, parsing url", err)
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true,
			SessionTicketsDisabled: true,
		},
	}

	log.Info("client succesfuly initiated")

	return &OlxClient{client: &http.Client{
		Transport: transport,
	}, logger: log}, nil
}

// Get .
func (c *OlxClient) Get(ctx context.Context, url string) (*QueryResult, error) {
	c.logger.Info("olx client got new request with url", url)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("cant reate request for olx client", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Host", "www.olx.ua")

	var data QueryResult
	err = c.httpDo(ctx, req, func(res *http.Response, err error) error {
		defer res.Body.Close()
		if err != nil {
			c.logger.Error("cant make http request to remote api form olx client", err)
			return err
		}

		c.logger.Info("request to olx api was made succesfuly")
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			c.logger.Error("cant decode body response in olx client")
			return err
		}

		c.logger.Info("decoded response successfuly")

		return nil

	})

	c.logger.Info("returninq query result from olx client")
	return &data, nil
}

func (c *OlxClient) httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	errch := make(chan error)
	defer close(errch)

	go func() {
		errch <- f(c.client.Do(req))
	}()

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		<-errch
		return ctx.Err()
	}
}
