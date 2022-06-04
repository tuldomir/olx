package service

import (
	"context"
	"errors"
	"fmt"
	"olx/db"
	"olx/models"
	"olx/olxapi"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

// Brocker .
type Brocker struct {
	client   *olxapi.Client
	db       *db.FileDb
	logger   *logrus.Logger
	interval int
	tries    int
}

// New .
func New(client *olxapi.Client, db *db.FileDb, log *logrus.Logger, interval, tries int) *Brocker {
	log.Info("initiating new server")
	return &Brocker{
		client:   client,
		db:       db,
		logger:   log,
		interval: interval,
		tries:    tries,
	}
}

func (b *Brocker) poll(ctx context.Context, iterval, tries int) chan *models.QueryResult {

	queries := make(chan *models.QueryResult)

	t := time.Duration(time.Duration(iterval) * time.Second)
	count := 0

	go func() {
		defer close(queries)

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(t):
				if count > tries {

					fmt.Println("setting new proxy")
					if err := b.client.InitProxy(); err != nil {
						return
					}

					count = 0
					continue
				}

				res, err := b.client.Get(ctx)
				if err != nil {
					b.logger.Errorf("cant get res from api, call #%v, err: ", count, err)
					count++
					t = time.Duration(0)
					continue
				}

				count = 0
				t = time.Duration(time.Duration(iterval) * time.Minute)
				queries <- res
			}
		}
	}()

	return queries
}

// Run .
func (b *Brocker) Run(ctx context.Context) chan []*models.Appartment {
	b.logger.Info("serving messages from server")

	msgs := make(chan []*models.Appartment)

	ch := b.poll(ctx, b.interval, b.tries)

	go func() {
		defer close(msgs)

		for v := range ch {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Println("data :", v)
				aps, err := b.processData(ctx, v)
				if err != nil {
					b.logger.Errorf("error while processing queryres, %v", err)
					continue
				}

				msgs <- aps
			}
		}
	}()

	return msgs
}

func (b *Brocker) processData(ctx context.Context, res *models.QueryResult) ([]*models.Appartment, error) {
	b.logger.Info("olx server processing data")
	var empty bool
	last, err := b.db.GetLast(ctx)
	if err != nil {
		b.logger.Error("olx server error ocucured ", err)
		switch {
		case errors.Is(err, db.ErrEmptyDB):
			b.logger.Info("error was empty error")
			empty = true
			break
		default:
			return nil, err
		}
	}

	organic := b.extractOrganic(ctx, res)

	if !empty {
		b.logger.Info("olx server db is not empty getting fresh")
		organic = b.sortFresh(ctx, organic, last.LastRefreshTime)
	}

	return organic, nil
}

func (b *Brocker) sortFresh(ctx context.Context, sl []*models.Appartment, last time.Time) []*models.Appartment {

	b.logger.Infof("sorting fresh appartments by time, %v", last)
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

	b.logger.Infof("length of slice is : %v", len(sl[:idx]))
	return sl[:idx]
}

func (b *Brocker) extractOrganic(ctx context.Context, res *models.QueryResult) []*models.Appartment {
	b.logger.Info("extracting only organic appartments")
	organic := make([]*models.Appartment, len(res.Metadata.Source.Organic))

	for i, v := range res.Metadata.Source.Organic {
		organic[i] = res.Data[v]
	}

	return organic
}
