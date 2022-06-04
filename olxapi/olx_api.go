package olxapi

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"olx/models"
	"os"

	"github.com/sirupsen/logrus"
)

// Client .
type Client struct {
	Client  *http.Client
	logger  *logrus.Logger
	api     string
	proxy   string
	proxies string
}

// New .
func New(proxies string, api string, log *logrus.Logger) (*Client, error) {
	log.Infof("initiating new olx client, proxies : %v", proxies)

	c := &Client{
		Client:  &http.Client{},
		logger:  log,
		api:     api,
		proxies: proxies,
	}

	if err := c.InitProxy(); err != nil {
		return nil, err
	}

	log.Info("client succesfuly initiated")
	return c, nil
}

// InitProxy .
func (c *Client) InitProxy() error {
	proxy, err := c.GetProxy()
	if err != nil {
		return err
	}

	c.logger.Infof("init new proxy: %v", proxy)

	return c.SetProxy(proxy)
}

// SetProxy .
func (c *Client) SetProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		c.logger.Error("error ocured in olx client, parsing url", err)
		return err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true,
			SessionTicketsDisabled: true,
		},
	}

	c.logger.Info("proxy succesfuly initiated")
	c.Client.Transport = transport
	c.proxy = proxy
	return nil
}

// GetProxy .
func (c *Client) GetProxy() (string, error) {
	file, err := os.Open(c.proxies)
	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if c.proxy == "" && line != "" {
			return line, nil
		}

		if c.proxy == line && scanner.Scan() {
			return scanner.Text(), nil
		}
	}

	return "", errors.New("empty proxy list")
}

// MakeRequest .
func (c *Client) MakeRequest(ctx context.Context) (*http.Request, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.api, nil)
	if err != nil {
		c.logger.Error("cant reate request for olx client", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Host", "www.olx.ua")

	return req, err
}

// Get .
func (c *Client) Get(ctx context.Context) (*models.QueryResult, error) {
	c.logger.Info("olx client trying new request to api")

	// ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	// defer cancel()

	req, err := c.MakeRequest(ctx)
	if err != nil {
		return nil, err
	}

	if err != nil {
		c.logger.Errorf("cant make req for api %v", err)
		return nil, err
	}

	var data models.QueryResult
	err = c.httpDo(ctx, req, func(res *http.Response, err error) error {
		// fmt.Println("here")
		defer func() {
			if res != nil {
				res.Body.Close()
			}
		}()

		// res, err := c.Client.Do(req)

		if err != nil {
			c.logger.Error("cant make http request to remote api form olx client", err)
			return err
		}

		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("cant get response from api, status: %v", res.StatusCode)
			c.logger.Error(err)
			return err
		}

		// bs, _ := ioutil.ReadAll(res.Body)
		// fmt.Println(string(bs))
		// err = json.Unmarshal(bs, &data)

		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			// if err != nil {
			c.logger.Error("cant decode body response in olx client")
			return err
		}
		// c.logger.Info("request to olx api was made succesfuly")

		c.logger.Info("decoded response successfuly")
		return nil
	})

	// fmt.Println("res: ", data, err)
	if err != nil {
		return nil, err
	}

	return &data, err
}

func (c *Client) httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	errch := make(chan error)
	defer close(errch)

	go func() {
		errch <- f(c.Client.Do(req))
	}()

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		<-errch
		return ctx.Err()
	}
}