package config

import "time"

// Order is a struct for Order datastore entities.
type Order struct {
	ID        int64     `jsonapi:"primary,order" datastore:"-"`
	CreatedAt time.Time `jsonapi:"attr,created_at"`
	Delivered bool      `jsonapi:"attr,delivered"`
}
