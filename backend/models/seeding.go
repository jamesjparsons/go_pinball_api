package models

type SeedingMethod string

const (
	SeedingMethodAverage  SeedingMethod = "AVERAGE"
	SeedingMethodRank     SeedingMethod = "RANK"
	SeedingMethodRandom   SeedingMethod = "RANDOM"
	SeedingMethodIFPARank SeedingMethod = "IFPA_RANK"
)
