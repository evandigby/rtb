package redis

import (
	"github.com/evandigby/rtb"
)

type RedisPacer struct {
	da     NoDbDataAccess
	banker rtb.Banker
}

func (p *RedisPacer) CanBid(account int64) bool {
	return true
}

func NewRedisPacer(da NoDbDataAccess, banker rtb.Banker) rtb.Pacer {
	p := new(RedisPacer)

	p.da = da
	p.banker = banker

	return p
}
