package model

import "github.com/oklog/ulid/v2"

type Prowler struct {
	UserID  ulid.ULID   `json:"user_id"`
	ID      ulid.ULID   `json:"id"`
	Type    ProwlerType `json:"type"`
	Query   string      `json:"query"`
	Targets Targets     `json:"targets"`
}

func NewProwl(a ...any) *Prowler {
	p := &Prowler{
		UserID: a[0].(ulid.ULID),
		ID:     a[1].(ulid.ULID),
		Type:   a[2].(ProwlerType),
	}
	if len(a) > 3 {
		p.Query = a[3].(string)
	}
	if len(a) > 4 {
		p.Targets = a[4].(Targets)
	}
	return p
}

func (p *Prowler) Hunt() {
	if p.Type == SearchProwlType {

	}
}

func (p *Prowler) HasTargets() bool {
	return p.Targets != nil && !p.Targets.None()
}

func (p *Prowler) IsTarget(t string) bool {
	return p.HasTargets() && p.Targets.Exist(t)
}
