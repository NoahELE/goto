package main

import (
	"errors"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type KeyUrlPair struct {
	gorm.Model
	Key string
	Url string
}

type UrlStore struct {
	db *gorm.DB
}

const MAX_RETRY = 10

func NewUrlStore() *UrlStore {
	// new sqlite connection
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		// panic when failed to connect to db
		panic(err)
	}

	// fit db to the model
	db.AutoMigrate(&KeyUrlPair{})

	return &UrlStore{db}
}

func (s *UrlStore) set(key, url string) error {
	// attempt to insert a record, if failed return the error
	res := s.db.Create(&KeyUrlPair{Key: key, Url: url})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *UrlStore) Get(key string) (string, error) {
	// select a record by key
	var ku KeyUrlPair
	res := s.db.First(&ku, "key = ?", key)

	// if not found, return the error
	if res.Error != nil {
		return "", res.Error
	}
	return ku.Url, nil
}

func (s *UrlStore) Put(url string) (string, error) {
	// select a record with url
	var ku KeyUrlPair
	res := s.db.First(&ku, "url = ?", url)

	// match res.Error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// if not found, gen a key and create a record
		l := 5
		retry := 0
		for {
			// gen a nanoid as the key
			key, err := gonanoid.New(l)
			if err != nil {
				return "", err
			}

			// match the set result
			if err := s.set(key, url); err != nil {
				// if unsuccessful, increase the length of nanoid by 1
				l++

				// increase counter by 1
				retry++
				// return the error when reaching max retry
				if retry >= MAX_RETRY {
					return "", fmt.Errorf(
						"exceeds MAX_RETRY (%d): db error is %w",
						MAX_RETRY,
						err,
					)
				}
			} else {
				// if the set operation is successful, return key
				return key, nil
			}
		}
	} else if res.Error == nil {
		// if found, return that record
		return ku.Key, nil
	} else {
		// return the error if not nil
		return "", res.Error
	}
}
