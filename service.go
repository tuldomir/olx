package main

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

// OlxService .
type OlxService struct {
	client *OlxClient
	db     *FileDb
	logger *logrus.Logger
}

// Result .
type Result struct {
	Data *QueryResult
	Err  error
}

// NewServ .
func NewServ(client *OlxClient, db *FileDb, log *logrus.Logger) *OlxService {
	log.Info("initiating new server")
	return &OlxService{
		client: client,
		db:     db,
		logger: log,
	}
}

// Serve .
func (s *OlxService) Serve(ctx context.Context, url string) chan *Msgs {
	msgs := make(chan *Msgs)

	s.logger.Info("serving messages from server")

	go func() {
		defer close(msgs)

		for {
			select {
			case <-ctx.Done():
				s.logger.Info("context was canceled, stopping server")
				return
			case <-time.After(5 * time.Minute):
				s.logger.Info("olx server 5 min is out making a request to olx api")
				res := s.doQuery(ctx, url, 5, 300)
				if res.Err != nil {
					s.logger.Error("olx server error occured, %v, maybe try to change proxy\n", res.Err)
					continue
					// TODO change client proxy
				}
				s.logger.Info("olx server query was successyly done")
				msg := s.processData(ctx, res)

				if msg.Err != nil {
					s.logger.Errorf("olx server got msg with error inside %v", msg.Err)
					continue
				}

				if len(msg.Msgs) < 1 {
					s.logger.Info("olx server: msg is empty nothing to update")
					continue
				}

				msgs <- msg
				if err := s.db.Save(ctx, res.Data.Data); err != nil {
					s.logger.Error("cant save to db", err)
				}
			}
		}
	}()

	return msgs
}

func (s *OlxService) doQuery(ctx context.Context, url string, tries, timeout int) *Result {
	qres := make(chan *Result)
	defer close(qres)
	s.logger.Infof("olx server making query to client tries: %v, timeout: %v", tries, timeout)

	queryFunc := func(ctx context.Context) {
		go func() {

			ctx, done := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer done()

			res, err := s.client.Get(ctx, url)
			if err != nil {
				// TODO here check err && try reconnect && maybe change client proxy
				s.logger.Error("olx server error ocured ", err)
			}

			s.logger.Info("olx server got response from client")
			qres <- &Result{Data: res, Err: err}
		}()

	}

	var res *Result

	ok := true
	for i := 0; i < tries && ok; i++ {
		queryFunc(ctx)
		select {
		case res = <-qres:
			if res.Err != nil {
				s.logger.Errorf("olx server error occured while making request, %v, next try %v", res.Err, i)
				continue
			}
			s.logger.Infof("olx server result is succesful breaking tries on try %v", i)
			ok = false

		case <-ctx.Done():
			return &Result{Err: ctx.Err()}
		}
	}
	s.logger.Info("olx server returning result from query")

	return res
}

// Msgs .
type Msgs struct {
	Msgs []*TeleMsg
	Err  error
}

func (s *OlxService) processData(ctx context.Context, res *Result) *Msgs {
	s.logger.Info("olx server processing data")
	var empty bool
	last, err := s.db.GetLast(ctx)
	if err != nil {
		s.logger.Error("olx server error ocucured ", err)
		switch {
		case errors.Is(err, ErrEmptyDB):
			s.logger.Info("error was empty error")
			empty = true
			break
		default:
			return &Msgs{Err: err}
		}
	}

	organic := s.extractOrganic(ctx, res.Data)

	if !empty {
		s.logger.Info("olx server db is not empty getting fresh")
		organic = s.sortFresh(ctx, organic, last.LastRefreshTime)
	}

	teleMsg := make([]*TeleMsg, len(organic))
	for i, v := range organic {
		teleMsg[i] = &TeleMsg{URL: v.URL, LastRefreshTime: v.LastRefreshTime}
	}

	return &Msgs{Msgs: teleMsg, Err: nil}

}

func (s *OlxService) sortFresh(ctx context.Context, sl []*Appartment, last time.Time) []*Appartment {

	s.logger.Infof("sorting fresh appartments by time, %v", last)
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

	s.logger.Infof("length of slice is : %v", len(sl[:idx]))

	return sl[:idx]
}

func (s *OlxService) extractOrganic(ctx context.Context, res *QueryResult) []*Appartment {
	s.logger.Info("extracting only organic appartments")
	organic := make([]*Appartment, len(res.Metadata.Source.Organic))

	for i, v := range res.Metadata.Source.Organic {
		organic[i] = res.Data[v]
	}

	return organic
}
