// dto contains type definitions for "Data Transfer Objects", structures that
// represent their data ver-batim over nushell or caldav

package dto

import (
	"time"

	"github.com/emersion/go-webdav/caldav"
	"github.com/teambition/rrule-go"
)

type PropValue struct {
	Value  string
	Params map[string][]string
}

type PropValueList []PropValue

type CalendarList []caldav.Calendar

type TimeSegment struct {
	Now          time.Time
	Duration     time.Duration
	ActiveEvents []string
}

type Timeline []TimeSegment

// rrule.RRule but with gob encoding support
type RRule struct {
	*rrule.RRule
}

func (r RRule) GobEncode() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *RRule) GobDecode(s []byte) error {
	rule, err := rrule.StrToRRule(string(s))
	if err != nil {
		return err
	}
	*r = RRule{RRule: rule}
	return nil
}
