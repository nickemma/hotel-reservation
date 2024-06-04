package db

import "github.com/nickemma/hotel-reservation/types"

type UserStore interface {
	GetUserByID(id string) (*types.User, error)
}

type MongoUserStore struct{}
