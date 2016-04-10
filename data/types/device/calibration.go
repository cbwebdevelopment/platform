package device

import (
	"github.com/tidepool-org/platform/data/types"

	"github.com/tidepool-org/platform/validate"
)

type Calibration struct {
	Value *float64 `json:"value" bson:"value" valid:"bloodvalue"`
	Units *string  `json:"units" bson:"units" valid:"mmolmgunits"`
	Base  `bson:",inline"`
}

func (b Base) makeCalibration(datum types.Datum, errs validate.ErrorProcessing) *Calibration {
	calibration := &Calibration{
		Value: datum.ToFloat64(types.BloodValueField.Name, errs),
		Units: datum.ToString(types.MmolOrMgUnitsField.Name, errs),
		Base:  b,
	}
	types.GetPlatformValidator().SetErrorReasons(failureReasons).Struct(calibration, errs)
	return calibration
}
