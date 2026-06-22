// Package combiner builds routes from per-leg quotes: a direct offer plus
// self-transfer combos through a list of hubs. Implements sources.Adapter.
package combiner

import (
	"context"
	"time"

	"flightmeta/internal/offer"
	"flightmeta/internal/sources"
)

// Hub is a connecting airport to try, with the layover (minutes) scheduled there.
type Hub struct {
	Code       string
	LayoverMin int
}

// DefaultHubs — common connecting points. Layover values give a mix of
// safe/risky connections.
var DefaultHubs = []Hub{
	{"IST", 180},
	{"DXB", 200},
	{"DOH", 170},
	{"AUH", 165},
	{"PEK", 240},
	{"TAS", 160},
	{"DEL", 190},
	{"BEG", 100}, // tight on purpose
}

// directPremium marks up the direct fare a bit so combos can undercut it.
const directPremium = 130 // percent

type Combiner struct {
	legs sources.LegSource
	hubs []Hub
}

func New(legs sources.LegSource, hubs []Hub) *Combiner {
	if len(hubs) == 0 {
		hubs = DefaultHubs
	}
	return &Combiner{legs: legs, hubs: hubs}
}

func (c *Combiner) Name() string { return "combiner" }

func (c *Combiner) Search(ctx context.Context, q sources.Query) ([]offer.Offer, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	dep := q.DepartDate
	if dep.IsZero() {
		dep = time.Now().Add(14 * 24 * time.Hour)
	}
	day := time.Date(dep.Year(), dep.Month(), dep.Day(), 8, 0, 0, 0, time.UTC) // first leg departs 08:00 UTC

	var offers []offer.Offer

	// Direct (official single ticket), with a convenience premium.
	if lq, ok, err := c.legs.CheapestLeg(ctx, q.Origin, q.Destination, dep, q); err == nil && ok {
		arr := day.Add(time.Duration(lq.DurationMin) * time.Minute)
		offers = append(offers, offer.Offer{
			ID:         "direct-" + q.Origin + q.Destination,
			Source:     c.Name(),
			Connection: offer.Official,
			Currency:   lq.Currency,
			PriceMinor: lq.PriceMinor * directPremium / 100,
			DeepLink:   lq.DeepLink,
			Segments:   []offer.Segment{seg(q.Origin, q.Destination, day, arr, lq)},
			Carriers:   []string{lq.Carrier},
		})
	}

	// Self-transfer via each hub.
	for _, h := range c.hubs {
		if h.Code == q.Origin || h.Code == q.Destination {
			continue
		}
		l1, ok1, err1 := c.legs.CheapestLeg(ctx, q.Origin, h.Code, dep, q)
		l2, ok2, err2 := c.legs.CheapestLeg(ctx, h.Code, q.Destination, dep, q)
		if err1 != nil || err2 != nil || !ok1 || !ok2 {
			continue
		}
		dep1 := day
		arr1 := dep1.Add(time.Duration(l1.DurationMin) * time.Minute)
		dep2 := arr1.Add(time.Duration(h.LayoverMin) * time.Minute)
		arr2 := dep2.Add(time.Duration(l2.DurationMin) * time.Minute)

		offers = append(offers, offer.Offer{
			ID:         "self-" + q.Origin + h.Code + q.Destination,
			Source:     c.Name(),
			Connection: offer.SelfTransfer,
			Currency:   l1.Currency,
			PriceMinor: l1.PriceMinor + l2.PriceMinor,
			DeepLink:   l1.DeepLink, // first-leg link; per-leg purchase shown in UI later
			Unique:     true,
			Segments: []offer.Segment{
				seg(q.Origin, h.Code, dep1, arr1, l1),
				seg(h.Code, q.Destination, dep2, arr2, l2),
			},
			Layovers: []offer.Layover{{
				Airport:      h.Code,
				Duration:     time.Duration(h.LayoverMin) * time.Minute,
				SelfTransfer: true,
			}},
			Carriers: dedupe(l1.Carrier, l2.Carrier),
		})
	}

	return offers, nil
}

func seg(from, to string, dep, arr time.Time, lq sources.LegQuote) offer.Segment {
	return offer.Segment{
		From: from, To: to,
		DepartUTC: dep, ArriveUTC: arr,
		MarketingCarrier: lq.Carrier,
		FlightNumber:     lq.FlightNumber,
	}
}

func dedupe(a, b string) []string {
	if a == b {
		return []string{a}
	}
	return []string{a, b}
}
