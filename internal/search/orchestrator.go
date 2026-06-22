// Package search fans out a query to all configured sources concurrently,
// enforces a per-source timeout, then merges, filters and ranks the results.
package search

import (
	"context"
	"log/slog"
	"time"

	"flightmeta/internal/offer"
	"flightmeta/internal/rank"
	"flightmeta/internal/sources"
	"flightmeta/internal/visa"
)

// Orchestrator coordinates the concurrent fan-out across data sources.
type Orchestrator struct {
	sources       []sources.Adapter
	sourceTimeout time.Duration
	visa          *visa.Resolver // optional; enriches layovers with transit-visa info
	log           *slog.Logger
}

// New builds an orchestrator over the given source adapters. resolver may be
// nil (no transit-visa enrichment).
func New(log *slog.Logger, sourceTimeout time.Duration, resolver *visa.Resolver, srcs ...sources.Adapter) *Orchestrator {
	return &Orchestrator{sources: srcs, sourceTimeout: sourceTimeout, visa: resolver, log: log}
}

// Result is the merged, ranked search response plus per-source diagnostics.
type Result struct {
	Offers         []offer.Offer `json:"offers"`
	Sources        []SourceStat  `json:"sources"`
	VisaDisclaimer string        `json:"visaDisclaimer,omitempty"`
}

// SourceStat reports how each source performed (count, error, latency).
type SourceStat struct {
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Err    string `json:"error,omitempty"`
	TookMs int64  `json:"tookMs"`
}

// Search queries every source concurrently. A slow or failing source is
// isolated by its own timeout and never blocks the others; its error is
// surfaced in SourceStat rather than failing the whole request.
func (o *Orchestrator) Search(ctx context.Context, q sources.Query) Result {
	type partial struct {
		stat   SourceStat
		offers []offer.Offer
	}
	ch := make(chan partial, len(o.sources))

	for _, src := range o.sources {
		go func(src sources.Adapter) {
			sctx, cancel := context.WithTimeout(ctx, o.sourceTimeout)
			defer cancel()

			start := time.Now()
			offers, err := src.Search(sctx, q)
			stat := SourceStat{
				Name:   src.Name(),
				Count:  len(offers),
				TookMs: time.Since(start).Milliseconds(),
			}
			if err != nil {
				stat.Err = err.Error()
				o.log.Warn("source failed", "source", src.Name(), "err", err)
			}
			ch <- partial{stat: stat, offers: offers}
		}(src)
	}

	var all []offer.Offer
	stats := make([]SourceStat, 0, len(o.sources))
	for range o.sources {
		p := <-ch
		stats = append(stats, p.stat)
		all = append(all, p.offers...)
	}

	// Transit-visa enrichment must run before ranking, since the
	// OnlyVisaFreeTransit filter reads each layover's VisaStatus.
	res := Result{Sources: stats}
	if o.visa != nil {
		o.visa.Enrich(all, q.Passport)
		res.VisaDisclaimer = o.visa.Disclaimer()
	}
	res.Offers = rank.Apply(all, q.Filters)
	return res
}
