package internal

import (
	"errors"
	"strconv"
	"time"
)

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

func (f *Frequency) String() string {
	return strconv.Itoa(f.Value) + string(f.Unit)
}

func (f *Frequency) Duration() time.Duration {
	return FrequencyUnits[f.Unit] * time.Duration(f.Value)
}

func (f *Frequency) Validate() error {
	if _, ok := FrequencyUnits[f.Unit]; !ok {
		return errors.New("invalid frequency unit, must be one of: m, h, d")
	} else if f.Unit == Minute && f.Value < 5 {
		return errors.New("frequency must be at least 5 minutes")
	}
	return nil
}
