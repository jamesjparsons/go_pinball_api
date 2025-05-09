package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Map of player count to points per position
type PointDistributionMap map[string][]float64

// Value implements the driver.Valuer interface for database serialization
func (p PointDistributionMap) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface for database deserialization
func (p *PointDistributionMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal PointDistributionMap value: %v", value)
	}
	return json.Unmarshal(bytes, p)
}
