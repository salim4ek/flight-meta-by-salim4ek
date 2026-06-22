// Package offer defines the source-agnostic flight offer model that every data
// source is normalized into. It is intentionally dependency-free to avoid
// import cycles with the source adapters, ranking, visa and connection
// packages that all build on top of it.
package offer

import "time"

// ConnectionType describes whether an itinerary is sold as a single protected
// ticket or stitched together from separate tickets (self-transfer / virtual
// interlining). Self-transfer is the core differentiator of this product, but
// it carries real risk: a missed connection is the traveller's problem because
// the legs are independent contracts of carriage.
type ConnectionType string

const (
	Official     ConnectionType = "official"      // one ticket, connection protected
	SelfTransfer ConnectionType = "self_transfer" // separate tickets, NOT protected
)

// Segment is a single marketed flight leg.
type Segment struct {
	From             string    `json:"from"` // IATA airport code
	To               string    `json:"to"`   // IATA airport code
	DepartUTC        time.Time `json:"departUtc"`
	ArriveUTC        time.Time `json:"arriveUtc"`
	MarketingCarrier string    `json:"marketingCarrier"` // IATA airline code
	FlightNumber     string    `json:"flightNumber"`
	Cabin            string    `json:"cabin,omitempty"`
}

// Layover is the gap between two consecutive segments at a connecting airport.
type Layover struct {
	Airport      string        `json:"airport"`
	Duration     time.Duration `json:"durationNs"`
	SelfTransfer bool          `json:"selfTransfer"` // separate tickets across this point
	// Risk is filled by the connection module (later phase); sources leave it
	// zero-valued. VisaStatus and TransitNote are filled by the visa module
	// based on the traveller's passport.
	Risk        string `json:"risk,omitempty"`        // "", safe, risky, infeasible
	VisaStatus  string `json:"visaStatus,omitempty"`  // no_visa, twov, visa_required, unknown
	TransitNote string `json:"transitNote,omitempty"` // human-readable transit-visa note
}

// Offer is one bookable (or redirectable) travel option.
type Offer struct {
	ID         string         `json:"id"`
	Source     string         `json:"source"` // kiwi, travelpayouts, mock, ...
	Segments   []Segment      `json:"segments"`
	Layovers   []Layover      `json:"layovers,omitempty"`
	Connection ConnectionType `json:"connection"`
	PriceMinor int64          `json:"priceMinor"` // price in minor currency units
	Currency   string         `json:"currency"`
	DeepLink   string         `json:"deepLink"` // partner redirect URL
	Carriers   []string       `json:"carriers"`
	// Unique marks an itinerary mainstream aggregators do not surface as a
	// single option (typically a self-transfer combo via some hub). This is the
	// headline value of the product and is route-agnostic — any hub, not a
	// hardcoded country.
	Unique bool `json:"unique"`
}

// Stops returns the number of connections (0 == direct).
func (o Offer) Stops() int {
	if len(o.Segments) == 0 {
		return 0
	}
	return len(o.Segments) - 1
}
