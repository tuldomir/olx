package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/sirupsen/logrus"
)

// ErrEmptyDB .
var ErrEmptyDB = errors.New("db is empty")

// FileDb .
type FileDb struct {
	appartments []*Appartment
	path        string
	logger      *logrus.Logger
}

// Open .
func Open(path string, log *logrus.Logger) (*FileDb, error) {
	log.Info("initiating db")
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("error ocured in db %v", err)
		return nil, err
	}

	var appartments []*Appartment
	if err := json.Unmarshal(bs, &appartments); err != nil {
		log.Errorf("error ocured in db decoding appartments %v", err)
		return nil, err
	}

	log.Info("db initiated")
	return &FileDb{appartments: appartments, path: path, logger: log}, nil
}

// Save .
func (db *FileDb) Save(ctx context.Context, sl []*Appartment) error {
	db.appartments = sl

	fmt.Println("got slice len: ", len(sl))

	db.logger.Infof("saving data to file %v", db.path)
	sort.Slice(sl, func(i, j int) bool {
		return sl[i].LastRefreshTime.After(sl[j].LastRefreshTime)
	})

	bs, err := json.Marshal(sl)
	if err != nil {
		db.logger.Errorf("error ocured in db encoding appartments %v", err)
		return err
	}

	if err := ioutil.WriteFile(db.path, bs, 0666); err != nil {
		db.logger.Errorf("error ocured in db cant save file %v", err)
		return err
	}

	db.logger.Infof("data saved successfuly")

	return nil
}

// GetLast .
func (db *FileDb) GetLast(ctx context.Context) (*Appartment, error) {
	db.logger.Info("getting last appartment served")
	if len(db.appartments) == 0 {
		db.logger.Errorf("no data in db %v", ErrEmptyDB)
		return nil, ErrEmptyDB
	}

	db.logger.Infof("last appartment retrieved %v, %v", db.appartments[0].ID, db.appartments[0].LastRefreshTime)
	return db.appartments[0], nil
}

// // Close .
// func (db *FileDb) Close() error {
// 	db.logger.Info("closing db")
// 	return db.file.Close()
// }
