package models

// ZoneSourceRule defined the saling rules of the traffic
type ZoneSourceRule struct {
	SourceID        uint64
	FlatratePercent float64
	MinimalBids     GeoBidSlice
}
