package utils

import (
	"time"
)

type ITimeGenerator interface {
	Now() time.Time
}

type TimeGenerator struct{}

func (t *TimeGenerator) Now() time.Time {
	return time.Now()
}

func NewTimeGenerator() ITimeGenerator {
	return &TimeGenerator{}
}

func GetCurrentTime() time.Time {
	return NewTimeGenerator().Now()
}
