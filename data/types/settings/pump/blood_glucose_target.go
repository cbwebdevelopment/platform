package pump

import (
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/data/blood/glucose"
)

type BloodGlucoseTarget struct {
	glucose.Target `bson:",inline"`

	Start *int `json:"start,omitempty" bson:"start,omitempty"`
}

func ParseBloodGlucoseTarget(parser data.ObjectParser) *BloodGlucoseTarget {
	var bloodGlucoseTarget *BloodGlucoseTarget
	if parser.Object() != nil {
		bloodGlucoseTarget = NewBloodGlucoseTarget()
		bloodGlucoseTarget.Parse(parser)
		parser.ProcessNotParsed()
	}
	return bloodGlucoseTarget
}

func ParseBloodGlucoseTargetArray(parser data.ArrayParser) *[]*BloodGlucoseTarget {
	var bloodGlucoseTargetArray *[]*BloodGlucoseTarget
	if parser.Array() != nil {
		bloodGlucoseTargetArray = &[]*BloodGlucoseTarget{}
		for index := range *parser.Array() {
			*bloodGlucoseTargetArray = append(*bloodGlucoseTargetArray, ParseBloodGlucoseTarget(parser.NewChildObjectParser(index)))
		}
		parser.ProcessNotParsed()
	}
	return bloodGlucoseTargetArray
}

func NewBloodGlucoseTarget() *BloodGlucoseTarget {
	return &BloodGlucoseTarget{}
}

func (t *BloodGlucoseTarget) Parse(parser data.ObjectParser) {
	t.Target.Parse(parser)

	t.Start = parser.ParseInteger("start")
}

func (t *BloodGlucoseTarget) Validate(validator data.Validator, units *string) {
	t.Target.Validate(validator, units)

	validator.ValidateInteger("start", t.Start).Exists().InRange(0, 86400000)
}
