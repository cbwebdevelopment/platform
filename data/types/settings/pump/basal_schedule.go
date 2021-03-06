package pump

import "github.com/tidepool-org/platform/data"

type BasalSchedule struct {
	Rate  *float64 `json:"rate,omitempty" bson:"rate,omitempty"`
	Start *int     `json:"start,omitempty" bson:"start,omitempty"`
}

func NewBasalSchedule() *BasalSchedule {
	return &BasalSchedule{}
}

func (b *BasalSchedule) Parse(parser data.ObjectParser) {
	b.Rate = parser.ParseFloat("rate")
	b.Start = parser.ParseInteger("start")
}

func (b *BasalSchedule) Validate(validator data.Validator) {
	validator.ValidateFloat("rate", b.Rate).Exists().InRange(0.0, 100.0)
	validator.ValidateInteger("start", b.Start).Exists().InRange(0, 86400000)
}

func (b *BasalSchedule) Normalize(normalizer data.Normalizer) {
}

func parseBasalSchedule(parser data.ObjectParser) *BasalSchedule {
	var basalSchedule *BasalSchedule
	if parser.Object() != nil {
		basalSchedule = NewBasalSchedule()
		basalSchedule.Parse(parser)
		parser.ProcessNotParsed()
	}
	return basalSchedule
}

func parseBasalScheduleArray(parser data.ArrayParser) *[]*BasalSchedule {
	var basalScheduleArray *[]*BasalSchedule
	if parser.Array() != nil {
		basalScheduleArray = &[]*BasalSchedule{}
		for index := range *parser.Array() {
			*basalScheduleArray = append(*basalScheduleArray, parseBasalSchedule(parser.NewChildObjectParser(index)))
		}
		parser.ProcessNotParsed()
	}
	return basalScheduleArray
}

func ParseBasalSchedulesMap(parser data.ObjectParser) *map[string]*[]*BasalSchedule {
	var basalScheduleMap *map[string]*[]*BasalSchedule
	if parser.Object() != nil {
		basalScheduleMap = &map[string]*[]*BasalSchedule{}
		for key := range *parser.Object() {
			(*basalScheduleMap)[key] = parseBasalScheduleArray(parser.NewChildArrayParser(key))
		}
		parser.ProcessNotParsed()
	}
	return basalScheduleMap
}
