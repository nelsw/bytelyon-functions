package model

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

type DateTime time.Time

func (v *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.Trim(s, `"`) // Remove quotes from the JSON string
	if s == "" || s == "null" {
		return nil // Handle empty or null strings
	}

	t, err := time.Parse(time.RFC1123, s) // Parse using your custom format
	if err != nil {
		return err
	}

	*v = DateTime(t)
	return nil
}

func (v *DateTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, `"`) // Remove quotes from the JSON string
	if s == "" || s == "null" {
		return nil // Handle empty or null strings
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*v = DateTime(t)
	return nil
}

func (v *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*v).Format(time.RFC3339) + `"`), nil
}

func (v *DateTime) ULID() ulid.ULID {
	return NewUlidFromTime(time.Time(*v))
}
