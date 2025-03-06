package entity

import (
	"fmt"
)

type SetClass struct {
	Name string `json:"name"`
}

var ValidValues = map[string]bool{
	"1": true,
	"2": true,
	"3": true,
	"4": true,
	"5": true,
	"6": true,
}

func (s *SetClass) Validate() error {
	if !ValidValues[s.Name] {
		return fmt.Errorf("name must be one of %v", getValidValues())
	}
	return nil
}

func getValidValues() []string {
	values := make([]string, 0, len(ValidValues))
	for key := range ValidValues {
		values = append(values, key)
	}
	return values
}
