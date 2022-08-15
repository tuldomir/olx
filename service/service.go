package service

import (
	"context"
	"errors"
	"fmt"
	"olx/db"
	"olx/models"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

// PollClient .
type PollClient interface {
	Get(context.Context, string) (*models.QueryResult, error)
	InitProxy() error
}

// IPoller .
type IPoller interface {
	Poll(context.Context, string) chan *models.QueryResult
}

// Poller .
type Poller struct {
	client   PollClient
	interval int
	logger   *logrus.Logger
}

// New .
func New(client PollClient, interval int, log *logrus.Logger) *Poller {
	log.Info("initiating new server")
	return &Poller{
		client:   client,
		interval: interval,
		logger:   log,
	}
}

// Poll .
func (p *Poller) Poll(ctx context.Context, url string) chan *models.QueryResult {

	type result struct {
		res *models.QueryResult
		err error
	}

	var (
		t       = time.Duration(time.Duration(p.interval) * time.Second)
		queries = make(chan *models.QueryResult)
		results = make(chan *result)
		guard   = make(chan struct{})
		count   = 0

		f = func(res *models.QueryResult, err error) *result {
			return &result{res: res, err: err}
		}
	)

	go func() {
		for range guard {
			select {
			case <-ctx.Done():
				return
			case results <- f(p.client.Get(ctx, url)):
			}
		}
	}()

	go func() {
		defer close(queries)
		defer close(results)
		defer close(guard)

		for {
			select {
			case <-ctx.Done():
				return
			case r := <-results:
				if r.err != nil {
					p.logger.Errorf("request not succeded err %v", r.err)
					count++

					if count > 5 {
						fmt.Println("too many errors try another proxy")
						if err := p.client.InitProxy(); err != nil {
							return
						}

						count = 0
					}

					guard <- struct{}{}
					continue
				}

				select {
				case <-ctx.Done():
					return
				case queries <- r.res:
				}

			case <-time.After(t):
				if count != 0 {
					continue
				}

				guard <- struct{}{}
			}
		}
	}()

	return queries
}

// Server .
type Server struct {
	poller IPoller
	db     *db.FileDb
	apiurl string
	logger *logrus.Logger
}

// NewServer .
func NewServer(poller IPoller, db *db.FileDb, url string, log *logrus.Logger) *Server {
	return &Server{
		poller: poller,
		db:     db,
		apiurl: url,
		logger: log,
	}
}

// Run .
func (s *Server) Run(ctx context.Context) chan []*models.Appartment {
	s.logger.Info("serving messages from server")

	msgs := make(chan []*models.Appartment)

	go func() {
		defer close(msgs)

		for v := range s.poller.Poll(ctx, s.apiurl) {
			select {
			case <-ctx.Done():
				return
			default:
				aps, err := s.processData(ctx, v)
				if err != nil {
					s.logger.Errorf("error while processing queryres, %v", err)
					continue
				}

				if err := s.db.Save(ctx, aps); err != nil {
					s.logger.Errorf("error cant save to db, %v", err)
					continue
				}

				select {
				case <-ctx.Done():
				case msgs <- aps:
				}
			}
		}
	}()

	return msgs
}

func (s *Server) processData(ctx context.Context, res *models.QueryResult) ([]*models.Appartment, error) {
	s.logger.Info("olx server processing data")
	var empty bool
	last, err := s.db.GetLast(ctx)
	if err != nil {
		s.logger.Error("olx server error ocucured ", err)
		switch {
		case errors.Is(err, db.ErrEmptyDB):
			s.logger.Info("error was empty error")
			empty = true
			break
		default:
			return nil, err
		}
	}

	organic := extractOrganic(res)
	if empty {
		return organic, nil
	}

	s.logger.Info("olx server db is not empty getting fresh")
	organic = sortFresh(organic, last.LastRefreshTime)
	return organic, nil
}

func sortFresh(sl []*models.Appartment, last time.Time) []*models.Appartment {

	sort.Slice(sl, func(i, j int) bool {
		return sl[i].LastRefreshTime.After(sl[j].LastRefreshTime)
	})

	var idx int
	for i, v := range sl {
		if v.LastRefreshTime.After(last) {
			continue
		}
		idx = i
		break
	}

	return sl[:idx]
}

func extractOrganic(res *models.QueryResult) []*models.Appartment {
	organic := make([]*models.Appartment, len(res.Metadata.Source.Organic))

	for i, v := range res.Metadata.Source.Organic {
		organic[i] = res.Data[v]
	}

	return organic
}
