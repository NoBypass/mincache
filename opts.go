package mincache

import "time"

type options struct {
	setAt    time.Time
	expireAt time.Time
	expireIn time.Duration
	expires  bool
}

func (o *options) expired() bool {
	return o.expires && ((o.setAt == time.Time{} && o.expireAt.Before(time.Now())) || (o.setAt != time.Time{} && o.setAt.Add(o.expireIn).Before(time.Now())))
}

func (o *options) expiration() time.Time {
	if o.expireAt.After(o.setAt) {
		return o.expireAt
	}
	return o.setAt.Add(o.expireIn)
}

func apply(opts []Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Option is used to configure when an item expires.
// Possible options are:
//
//	ExpireAt(time.Time) - set the expiration time
//	ExpireIn(time.Duration) - set the expiration time from now
type Option func(*options)

func ExpireAt(t time.Time) Option {
	return func(o *options) {
		o.expires = true
		o.expireAt = t
	}
}

func ExpireIn(d time.Duration) Option {
	return func(o *options) {
		o.setAt = time.Now()
		o.expires = true
		o.expireIn = d
	}
}
