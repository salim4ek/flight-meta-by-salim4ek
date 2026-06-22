package sources

import (
	"context"
	"time"
)

// LegQuote is the cheapest price for one city-pair on a date. Times are
// scheduled by the combiner, so a quote just carries duration, price and a link.
type LegQuote struct {
	Carrier      string
	FlightNumber string
	DurationMin  int
	PriceMinor   int64
	Currency     string
	DeepLink     string
}

// LegSource gives per-leg cheapest prices. The combiner stitches legs across
// hubs into self-transfer routes. A real price feed (Travelpayouts) plugs in
// here in place of mockleg.
type LegSource interface {
	Name() string
	CheapestLeg(ctx context.Context, from, to string, date time.Time, q Query) (LegQuote, bool, error)
}
