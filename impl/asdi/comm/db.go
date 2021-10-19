package main

import (
	"errors"
	"sync"
)

type Database struct {
	Users sync.Map
	Txns  sync.Map
}

func InitDatabase() *Database {
	return &Database{
		Users: sync.Map{},
		Txns:  sync.Map{},
	}
}

func (db *Database) CreateUser(name string, balance int) {
	db.Users.Store(name, balance)
}

func (db *Database) GetUser(name string) (int, error) {
	balance, ok := db.Users.Load(name)
	if !ok {
		return 0, errors.New("Database internal error")
	}
	return balance.(int), nil
}
