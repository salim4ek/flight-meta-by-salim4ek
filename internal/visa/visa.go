// Package visa works out transit-visa status per passport at each layover.
// Data is hand-seeded and indicative (see data/visa_rules.json), not official —
// hence the disclaimer. A real feed (Timatic/sherpa) can drop in behind Resolver.
package visa

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"flightmeta/internal/offer"
)

//go:embed data/airports.json data/visa_rules.json
var dataFS embed.FS

// Status is the transit-visa outcome for a passport at a transit country.
type Status string

const (
	NoVisa       Status = "no_visa"
	TWOV         Status = "twov" // transit without visa, conditional
	VisaRequired Status = "visa_required"
	Unknown      Status = "unknown"
)

// Rule is one transit rule. SelfTransferStatus, if set, overrides Status for
// self-transfer layovers (you exit landside, so airside-only rules don't help).
type Rule struct {
	Status             Status `json:"status"`
	SelfTransferStatus Status `json:"selfTransferStatus,omitempty"`
	Hours              int    `json:"hours,omitempty"`
	Note               string `json:"note"`
	Conditions         string `json:"conditions,omitempty"`
}

// Resolver maps airports to countries and looks up transit rules.
type Resolver struct {
	airports   map[string]string             // IATA -> ISO2 country
	rules      map[string]map[string]Rule    // passport -> country -> Rule
	disclaimer string
}

type airportEntry struct {
	Country string `json:"country"`
}

type rulesFile struct {
	Disclaimer string                     `json:"disclaimer"`
	Asof       string                     `json:"asof"`
	Rules      map[string]map[string]Rule `json:"rules"`
}

// Load reads the embedded seed data and returns a ready Resolver.
func Load() (*Resolver, error) {
	ab, err := dataFS.ReadFile("data/airports.json")
	if err != nil {
		return nil, fmt.Errorf("read airports: %w", err)
	}
	var rawAirports map[string]json.RawMessage
	if err := json.Unmarshal(ab, &rawAirports); err != nil {
		return nil, fmt.Errorf("parse airports: %w", err)
	}
	airports := make(map[string]string, len(rawAirports))
	for code, raw := range rawAirports {
		var a airportEntry
		if err := json.Unmarshal(raw, &a); err != nil || a.Country == "" {
			continue // tolerate "_comment" strings and any malformed rows
		}
		airports[strings.ToUpper(code)] = strings.ToUpper(a.Country)
	}

	rb, err := dataFS.ReadFile("data/visa_rules.json")
	if err != nil {
		return nil, fmt.Errorf("read rules: %w", err)
	}
	var rf rulesFile
	if err := json.Unmarshal(rb, &rf); err != nil {
		return nil, fmt.Errorf("parse rules: %w", err)
	}
	rules := make(map[string]map[string]Rule, len(rf.Rules))
	for passport, byCountry := range rf.Rules {
		m := make(map[string]Rule, len(byCountry))
		for country, rule := range byCountry {
			m[strings.ToUpper(country)] = rule
		}
		rules[strings.ToUpper(passport)] = m
	}

	return &Resolver{airports: airports, rules: rules, disclaimer: rf.Disclaimer}, nil
}

// Disclaimer returns the indicative-data warning to show to users.
func (r *Resolver) Disclaimer() string { return r.disclaimer }

// Country returns the ISO2 country for an airport (ok=false if unknown).
func (r *Resolver) Country(iata string) (string, bool) {
	c, ok := r.airports[strings.ToUpper(iata)]
	return c, ok
}

// lookup returns the rule for a passport at a country, falling back to the
// passport's DEFAULT entry, then to Unknown.
func (r *Resolver) lookup(passport, country string) Rule {
	if pm, ok := r.rules[strings.ToUpper(passport)]; ok {
		if rule, ok := pm[strings.ToUpper(country)]; ok {
			return rule
		}
		if def, ok := pm["DEFAULT"]; ok {
			return def
		}
	}
	return Rule{Status: Unknown, Note: "Правила транзита не определены — уточните у авиакомпании/консульства"}
}

// Enrich fills VisaStatus and TransitNote on every layover of every offer for
// the given passport.
func (r *Resolver) Enrich(offers []offer.Offer, passport string) {
	for i := range offers {
		for j := range offers[i].Layovers {
			lay := &offers[i].Layovers[j]
			country, ok := r.Country(lay.Airport)
			if !ok {
				lay.VisaStatus = string(Unknown)
				lay.TransitNote = "Страна хаба неизвестна — уточните транзит"
				continue
			}
			rule := r.lookup(passport, country)
			status := rule.Status
			if lay.SelfTransfer && rule.SelfTransferStatus != "" {
				status = rule.SelfTransferStatus
			}
			lay.VisaStatus = string(status)
			lay.TransitNote = note(country, status, rule)
		}
	}
}

func note(country string, status Status, rule Rule) string {
	var b strings.Builder
	b.WriteString(country)
	b.WriteString(": ")
	switch status {
	case NoVisa:
		b.WriteString(rule.Note)
	case TWOV:
		if rule.Hours > 0 {
			b.WriteString(fmt.Sprintf("безвиз-транзит до %dч", rule.Hours))
		} else {
			b.WriteString("безвиз-транзит (TWOV)")
		}
	case VisaRequired:
		b.WriteString("нужна транзитная виза")
	default:
		b.WriteString(rule.Note)
	}
	if rule.Conditions != "" {
		b.WriteString(" · ")
		b.WriteString(rule.Conditions)
	}
	return b.String()
}
