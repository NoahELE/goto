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
	db    *gorm.DB
	cache chan *KeyUrlPair
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

	// init cache
	cache := make(chan *KeyUrlPair)
	store := &UrlStore{db, cache}
	go store.saveLoop()
	return store
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
	// check if the url is already in db
	var ku KeyUrlPair
	res := s.db.First(&ku, "url = ?", url)

	// match res.Error
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// if not found, gen a key and push it to cache
		l := 5
		retry := 0
		for {
			// generate a nanoid as the key
			key, err := gonanoid.New(l)
			if err != nil {
				return "", err
			}

			// test if the key already exists
			if _, err := s.Get(key); err != nil {
				// increase length of generated nanoid
				l++
				// return error if exceeds max retry
				retry++
				if retry >= MAX_RETRY {
					return "", fmt.Errorf(
						"exceeds MAX_RETRY (%d): db error is %w",
						MAX_RETRY,
						err,
					)
				}
			} else {
				s.cache <- &KeyUrlPair{Key: key, Url: url}
				return key, nil
			}
		}
	} else if res.Error == nil {
		// if found, return that record
		return ku.Key, nil
	} else {
		// return other errors
		return "", res.Error
	}
}

func (s *UrlStore) saveLoop() {
	ku := <-s.cache
	s.db.Create(ku)
}
