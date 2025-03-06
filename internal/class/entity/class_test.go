package entity_test

import (
	"testing"

	"github.com/ghulammuzz/misterblast/internal/class/entity"
	"github.com/stretchr/testify/assert"
)

func TestValidate_ValidClass(t *testing.T) {
	tests := []struct {
		name        string
		input       entity.SetClass
		expectError bool
	}{
		{"Valid name 1", entity.SetClass{Name: "1"}, false},
		{"Valid name 2", entity.SetClass{Name: "2"}, false},
		{"Invalid name", entity.SetClass{Name: "invalid"}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetValidValues(t *testing.T) {
	expected := []string{"1", "2", "3", "4", "5", "6"}
	actual := entity.ValidValues

	assert.Len(t, actual, len(expected))
	for _, val := range expected {
		assert.True(t, actual[val])
	}
}
