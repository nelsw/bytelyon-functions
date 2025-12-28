package model

import (
	"errors"
	"time"
)

type Frequency struct {
	Unit  FrequencyUnit `json:"unit"`
	Value int           `json:"value"`
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
