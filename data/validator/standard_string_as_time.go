package validator

/* CHECKLIST
 * [x] Uses interfaces as appropriate
 * [x] Private package variables use underscore prefix
 * [x] All parameters validated
 * [x] All errors handled
 * [x] Reviewed for concurrency safety
 * [x] Code complete
 * [x] Full test coverage
 */

import (
	"time"

	"github.com/tidepool-org/platform/data"
)

type StandardStringAsTime struct {
	context     data.Context
	reference   interface{}
	stringValue *string
	timeValue   *time.Time
	timeLayout  string
}

func NewStandardStringAsTime(context data.Context, reference interface{}, stringValue *string, timeLayout string) *StandardStringAsTime {
	if context == nil {
		return nil
	}
	if timeLayout == "" {
		return nil
	}

	standardStringAsTime := &StandardStringAsTime{
		context:     context,
		reference:   reference,
		stringValue: stringValue,
		timeLayout:  timeLayout,
	}
	standardStringAsTime.parse()
	return standardStringAsTime
}

func (s *StandardStringAsTime) Exists() data.Time {
	if s.stringValue == nil {
		s.context.AppendError(s.reference, ErrorValueDoesNotExist())
	}
	return s
}

func (s *StandardStringAsTime) After(limit time.Time) data.Time {
	if s.timeValue != nil {
		if !s.timeValue.After(limit) {
			s.context.AppendError(s.reference, ErrorTimeNotAfter(*s.timeValue, limit, s.timeLayout))
		}
	}
	return s
}

func (s *StandardStringAsTime) AfterNow() data.Time {
	if s.timeValue != nil {
		if !s.timeValue.After(time.Now()) {
			s.context.AppendError(s.reference, ErrorTimeNotAfterNow(*s.timeValue, s.timeLayout))
		}
	}
	return s
}

func (s *StandardStringAsTime) Before(limit time.Time) data.Time {
	if s.timeValue != nil {
		if !s.timeValue.Before(limit) {
			s.context.AppendError(s.reference, ErrorTimeNotBefore(*s.timeValue, limit, s.timeLayout))
		}
	}
	return s
}

func (s *StandardStringAsTime) BeforeNow() data.Time {
	if s.timeValue != nil {
		if !s.timeValue.Before(time.Now()) {
			s.context.AppendError(s.reference, ErrorTimeNotBeforeNow(*s.timeValue, s.timeLayout))
		}
	}
	return s
}

func (s *StandardStringAsTime) parse() {
	if s.stringValue != nil {
		if timeValue, err := time.Parse(s.timeLayout, *s.stringValue); err != nil {
			s.context.AppendError(s.reference, ErrorTimeNotValid(*s.stringValue, s.timeLayout))
		} else {
			s.timeValue = &timeValue
		}
	}
}