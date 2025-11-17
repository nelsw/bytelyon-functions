package model

import "time"

type FrequencyUnit string

const (
	Minute FrequencyUnit = "m"
	Hour   FrequencyUnit = "h"
	Day    FrequencyUnit = "d"
)

var FrequencyUnits = map[FrequencyUnit]time.Duration{
	Minute: time.Minute,
	Hour:   time.Hour,
	Day:    time.Hour * 24,
}

type Frequency struct {
	Unit  FrequencyUnit `json:"unit"`
	Value int           `json:"value"`
}

func (f Frequency) Duration() time.Duration {
	return FrequencyUnits[f.Unit] * time.Duration(f.Value)
}
