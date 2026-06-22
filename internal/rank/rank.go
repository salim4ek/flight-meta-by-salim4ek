// Package rank applies user filters and orders offers. Ranking is intentionally
// simple for now (cheapest first); richer scoring (transit safety, route
// uniqueness boosts) comes in later phases.
package rank

import (
	"sort"
	"strings"

	"flightmeta/internal/offer"
	"flightmeta/internal/sources"
)

// Apply filters the offers and returns a new, price-sorted slice.
func Apply(offers []offer.Offer, f sources.Filters) []offer.Offer {
	out := make([]offer.Offer, 0, len(offers))
	for _, o := range offers {
		if !passStops(o, f.Stops) {
			continue
		}
		if o.Connection == offer.SelfTransfer && !f.AllowSelfTransfer {
			continue
		}
		if f.OnlyVisaFreeTransit && needsTransitVisa(o) {
			continue
		}
		if !passAirlines(o, f.IncludeAirlines, f.ExcludeAirlines) {
			continue
		}
		out = append(out, o)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].PriceMinor < out[j].PriceMinor
	})
	return out
}

// needsTransitVisa reports whether any layover definitely requires a transit
// visa for the traveller (per the visa module's enrichment).
func needsTransitVisa(o offer.Offer) bool {
	for _, l := range o.Layovers {
		if l.VisaStatus == "visa_required" {
			return true
		}
	}
	return false
}

func passStops(o offer.Offer, mode sources.StopsMode) bool {
	s := o.Stops()
	switch mode {
	case sources.StopsDirect:
		return s == 0
	case sources.StopsOne:
		return s == 1
	case sources.StopsOnePlus:
		return s >= 1
	default:
		return true
	}
}

func passAirlines(o offer.Offer, include, exclude []string) bool {
	has := func(set []string, code string) bool {
		for _, c := range set {
			if strings.EqualFold(c, code) {
				return true
			}
		}
		return false
	}
	// Exclude wins: drop the offer if any of its carriers is excluded.
	for _, c := range o.Carriers {
		if has(exclude, c) {
			return false
		}
	}
	// Include (whitelist): keep only if at least one carrier matches.
	if len(include) > 0 {
		for _, c := range o.Carriers {
			if has(include, c) {
				return true
			}
		}
		return false
	}
	return true
}
