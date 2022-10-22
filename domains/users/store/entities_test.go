package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringifyGender(t *testing.T) {
	table := []struct {
		name   string
		gender Gender
		exp    string
	}{
		{
			name:   "Male Test",
			gender: GenderMale,
			exp:    "Male",
		},
		{
			name:   "Female Test",
			gender: GenderFemale,
			exp:    "Female",
		},
		{
			name:   "Others",
			gender: GenderOthers,
			exp:    "Others",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.gender.String()
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestGetGenderFromString(t *testing.T) {
	table := []struct {
		name   string
		gender string
		exp    Gender
	}{
		{
			name:   "Male Test",
			gender: "Male",
			exp:    GenderMale,
		},
		{
			name:   "Female Test",
			gender: "Female",
			exp:    GenderFemale,
		},
		{
			name:   "Others Test",
			gender: "Others",
			exp:    GenderOthers,
		},
		{
			name:   "Accomodate Test",
			gender: "Non-Binary",
			exp:    GenderOthers,
		},
		{
			name:   "Lower Case Male Test",
			gender: "male",
			exp:    GenderMale,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := GenderFromString(tt.gender)
			require.Equal(t, tt.exp, got)
		})
	}
}
