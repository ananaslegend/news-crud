package model

import (
	"errors"
	"time"
)

var (
	ErrDateFromAfterDateTo = errors.New("DateFrom should be before DateTo")
)

type Filter struct {
	DateFrom time.Time
	DateTo   time.Time
}

func (f Filter) Validation() error {
	if f.DateFrom.After(f.DateTo) {
		return ErrDateFromAfterDateTo
	}
	return nil
}
