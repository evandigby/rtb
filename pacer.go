package rtb

import (
	"time"
)

// Defines an object that can set the pace of campaign spending
type Pacer interface {
	CanBid(account int64) bool
}

// Defines a specific type of pacer that will pace bids of a time period
type TimeSegmentedPacer interface {
	Pacer
	Segment() time.Duration
}
