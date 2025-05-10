package models

type GroupOrdering string

const (
	GroupOrderingRandom GroupOrdering = "RANDOM"
	GroupOrderingSeeded GroupOrdering = "SEEDED"
)
