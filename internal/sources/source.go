// Package sources defines the contract every flight data source implements and
// the query/filter types shared across the search pipeline.
package sources

import (
	"context"
	"time"

	"flightmeta/internal/offer"
)

// StopsMode buckets the user's stop preference. Empty means "any".
type StopsMode string

const (
	StopsAny     StopsMode = ""
	StopsDirect  StopsMode = "direct"   // 0 stops
	StopsOne     StopsMode = "one"      // exactly 1 stop
	StopsOnePlus StopsMode = "one_plus" // 1 or more stops
)

// Filters are applied after sources return their raw offers.
type Filters struct {
	Stops             StopsMode
	IncludeAirlines   []string // if non-empty, keep only offers using these carriers
	ExcludeAirlines   []string
	AllowSelfTransfer bool
	// OnlyVisaFreeTransit drops offers that have any layover where the
	// traveller's passport definitely needs a transit visa (VisaStatus ==
	// "visa_required"). Requires the visa module to have enriched the offers.
	OnlyVisaFreeTransit bool
}

// Query is a single search request in source-agnostic form.
type Query struct {
	Origin      string     // IATA
	Destination string     // IATA
	DepartDate  time.Time  // departure day
	ReturnDate  *time.Time // nil == one-way
	Passengers  int
	Passport    string // ISO country, e.g. RU, KZ — drives transit-visa logic later
	Currency    string
	Filters     Filters
}

// Adapter is implemented by every data source (Kiwi, Travelpayouts, Amadeus,
// mock, ...). Implementations MUST respect ctx cancellation/timeout and never
// block past it — the orchestrator fans out to all sources concurrently and
// applies a per-source deadline.
type Adapter interface {
	Name() string
	Search(ctx context.Context, q Query) ([]offer.Offer, error)
}
