// Package offer holds the flight offer model shared across the app. No
// dependencies on purpose — everything else imports it.
package offer

import "time"

// ConnectionType is how the legs are ticketed: one ticket (protected) or
// separate tickets (self-transfer — cheaper, but a missed connection is on you).
type ConnectionType string

const (
	Official     ConnectionType = "official"
	SelfTransfer ConnectionType = "self_transfer"
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
	SelfTransfer bool          `json:"selfTransfer"`
	// filled in after search (connection + visa enrichment)
	Risk        string `json:"risk,omitempty"`        // safe | risky | infeasible
	VisaStatus  string `json:"visaStatus,omitempty"`  // no_visa | twov | visa_required | unknown
	TransitNote string `json:"transitNote,omitempty"`
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
	Unique     bool           `json:"unique"` // combo the big aggregators don't show
}

// Stops returns the number of connections (0 == direct).
func (o Offer) Stops() int {
	if len(o.Segments) == 0 {
		return 0
	}
	return len(o.Segments) - 1
}
