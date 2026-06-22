// Package connection checks whether a layover leaves enough time for the next
// flight. Self-transfer needs more (re-check bags, security again).
package connection

import "flightmeta/internal/offer"

// Minimum connection time (minutes), international.
const (
	OfficialMinMinutes     = 90
	SelfTransferMinMinutes = 150
)

// Risk levels written onto offer.Layover.Risk.
const (
	RiskSafe       = "safe"
	RiskRisky      = "risky"
	RiskInfeasible = "infeasible"
)

// RequiredMinutes returns the minimum connection time for a layover type.
func RequiredMinutes(selfTransfer bool) int {
	if selfTransfer {
		return SelfTransferMinMinutes
	}
	return OfficialMinMinutes
}

// Enrich sets Risk on each layover from its duration vs the required MCT.
// Self-transfer through a country that needs a transit visa is infeasible
// (you'd have to clear immigration to switch tickets).
func Enrich(offers []offer.Offer) {
	for i := range offers {
		for j := range offers[i].Layovers {
			lay := &offers[i].Layovers[j]
			required := RequiredMinutes(lay.SelfTransfer)
			mins := int(lay.Duration.Minutes())

			switch {
			case lay.SelfTransfer && lay.VisaStatus == "visa_required":
				lay.Risk = RiskInfeasible
			case mins >= required:
				lay.Risk = RiskSafe
			case mins*10 >= required*7: // >= 70% of required
				lay.Risk = RiskRisky
			default:
				lay.Risk = RiskInfeasible
			}
		}
	}
}

// WorstRisk returns the most severe layover risk in an offer ("" if none).
func WorstRisk(o offer.Offer) string {
	worst := ""
	rank := map[string]int{"": 0, RiskSafe: 1, RiskRisky: 2, RiskInfeasible: 3}
	for _, l := range o.Layovers {
		if rank[l.Risk] > rank[worst] {
			worst = l.Risk
		}
	}
	return worst
}
