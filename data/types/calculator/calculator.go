package calculator

import (
	"reflect"

	validator "gopkg.in/bluesuncorp/validator.v8"

	"github.com/tidepool-org/platform/data/types"
	"github.com/tidepool-org/platform/validate"
)

func init() {
	types.GetPlatformValidator().RegisterValidation(insulinSensitivityField.Tag, InsulinSensitivityValidator)
	types.GetPlatformValidator().RegisterValidation(insulinOnBoardField.Tag, InsulinOnBoardValidator)
}

type Event struct {
	*Recommended        `json:"recommended,omitempty" bson:"recommended,omitempty"`
	*BloodGlucoseTarget `json:"bgTarget,omitempty" bson:"bgTarget,omitempty"`
	*Bolus              `json:"bolus,omitempty" bson:"bolus,omitempty"`

	CarbohydrateInput  *int     `json:"carbInput,omitempty" bson:"carbInput,omitempty" valid:"omitempty,required"`
	InsulinOnBoard     *float64 `json:"insulinOnBoard,omitempty" bson:"insulinOnBoard,omitempty" valid:"omitempty,insulinvalue"`
	InsulinSensitivity *int     `json:"insulinSensitivity,omitempty" bson:"insulinSensitivity,omitempty" valid:"omitempty,insulinsensitivity"`
	BloodGlucoseInput  *float64 `json:"bgInput,omitempty" bson:"bgInput,omitempty" valid:"-"`
	Units              *string  `json:"units" bson:"units" valid:"-"`
	types.Base         `bson:",inline"`
}

type Recommended struct {
	Carbohydrate *float64 `json:"carb" bson:"carb" valid:"required"`
	Correction   *float64 `json:"correction" bson:"correction"`
	Net          *float64 `json:"net" bson:"net" valid:"required"`
}

type Bolus struct {
	Type     *string `json:"type" bson:"type" valid:"-"`
	SubType  *string `json:"subType" bson:"subType" valid:"bolussubtype"`
	Time     *string `json:"time" bson:"time" valid:"timestr"`
	DeviceID *string `json:"deviceId" bson:"deviceId" valid:"required"`
}

type BloodGlucoseTarget struct {
	High *float64 `json:"high" bson:"high" valid:"-"`
	Low  *float64 `json:"low" bson:"low" valid:"-"`
}

//Name is currently `wizard` for backwards compatatbilty but will be migrated to `calculator`
const Name = "wizard"

var (
	carbohydrateInputField = types.DatumFieldInformation{
		DatumField: &types.DatumField{Name: "carbInput"},
		Tag:        "required",
		Message:    "This is a required field",
	}

	bloodGlucoseInputField = types.DatumFieldInformation{
		DatumField: &types.DatumField{Name: "bgInput"},
		Tag:        types.BloodGlucoseValueField.Tag,
		Message:    types.BloodGlucoseValueField.Message,
	}

	insulinSensitivityField = types.DatumFieldInformation{
		DatumField: &types.DatumField{Name: "insulinSensitivity"},
		Tag:        "insulinsensitivity",
		Message:    "This is a required field",
	}

	insulinOnBoardField = types.DatumFieldInformation{
		DatumField: &types.DatumField{Name: "insulinOnBoard"},
		Tag:        "insulinvalue",
		Message:    "This is a required field",
	}

	carbField = types.DatumFieldInformation{
		DatumField: &types.DatumField{Name: "carb"},
		Tag:        "required",
		Message:    "This is a required field",
	}

	netField = types.DatumFieldInformation{
		DatumField: &types.DatumField{Name: "net"},
		Tag:        "required",
		Message:    "This is a required field",
	}

	correctionField = types.DatumFieldInformation{
		DatumField: &types.DatumField{Name: "correction"},
		Tag:        "required",
		Message:    "This is a required field",
	}

	failureReasons = validate.FailureReasons{
		//"BloodGlucoseInput": validate.ValidationInfo{FieldName: bloodGlucoseInputField.Name, Message: bloodGlucoseInputField.Message},
		"CarbohydrateInput": validate.ValidationInfo{FieldName: carbohydrateInputField.Name, Message: carbohydrateInputField.Message},
		//"Units":             validate.ValidationInfo{FieldName: types.MmolOrMgUnitsField.Name, Message: types.MmolOrMgUnitsField.Message},
		"InsulinOnBoard": validate.ValidationInfo{FieldName: types.BolusSubTypeField.Name, Message: types.BolusSubTypeField.Message},

		"Bolus.SubType":  validate.ValidationInfo{FieldName: "bolus/" + types.BolusSubTypeField.Name, Message: types.BolusSubTypeField.Message},
		"Bolus.DeviceID": validate.ValidationInfo{FieldName: "bolus/" + types.BaseDeviceIDField.Name, Message: types.BaseDeviceIDField.Message},
		"Bolus.Time":     validate.ValidationInfo{FieldName: "bolus/" + types.TimeStringField.Name, Message: types.TimeStringField.Message},

		//"BloodGlucoseTarget.High": validate.ValidationInfo{FieldName: "bgTarget/high", Message: types.BloodGlucoseValueField.Message},
		//"BloodGlucoseTarget.Low":  validate.ValidationInfo{FieldName: "bgTarget/low", Message: types.BloodGlucoseValueField.Message},

		"Recommended.Net":          validate.ValidationInfo{FieldName: "recommended/" + netField.Name, Message: netField.Message},
		"Recommended.Correction":   validate.ValidationInfo{FieldName: "recommended/" + correctionField.Name, Message: correctionField.Message},
		"Recommended.Carbohydrate": validate.ValidationInfo{FieldName: "recommended/" + carbField.Name, Message: carbField.Message},
	}
)

func buildRecommended(recommendedDatum types.Datum, errs validate.ErrorProcessing) *Recommended {
	return &Recommended{
		Carbohydrate: recommendedDatum.ToFloat64(carbField.Name, errs),
		Correction:   recommendedDatum.ToFloat64(correctionField.Name, errs),
		Net:          recommendedDatum.ToFloat64(netField.Name, errs),
	}
}

func buildBolus(bolusDatum types.Datum, errs validate.ErrorProcessing) *Bolus {
	bolusType := "bolus"
	return &Bolus{
		Type:     &bolusType,
		SubType:  bolusDatum.ToString(types.BolusSubTypeField.Name, errs),
		Time:     bolusDatum.ToString(types.TimeStringField.Name, errs),
		DeviceID: bolusDatum.ToString("deviceId", errs),
	}
}

func buildBloodGlucoseTarget(units *string, bgTargetDatum types.Datum, errs validate.ErrorProcessing) *BloodGlucoseTarget {
	bgTarget := &BloodGlucoseTarget{
		High: bgTargetDatum.ToFloat64("high", errs),
		Low:  bgTargetDatum.ToFloat64("low", errs),
	}
	bgTarget.High, _ = types.NewBloodGlucoseValidation(bgTarget.High, units).SetValueErrorPath("bgTarget/high").ValidateAndConvertBloodGlucoseValue(errs)
	bgTarget.Low, _ = types.NewBloodGlucoseValidation(bgTarget.Low, units).SetValueErrorPath("bgTarget/low").ValidateAndConvertBloodGlucoseValue(errs)

	return bgTarget
}

func Build(datum types.Datum, errs validate.ErrorProcessing) *Event {

	originalBloodGlucoseUnits := datum.ToString(types.MmolOrMgUnitsField.Name, errs)

	var bolus *Bolus
	bolusDatum, ok := datum["bolus"].(map[string]interface{})
	if ok {
		bolus = buildBolus(bolusDatum, errs)
	}

	var bloodGlucoseTarget *BloodGlucoseTarget
	bloodGlucoseTargetDatum, ok := datum["bgTarget"].(map[string]interface{})
	if ok {
		bloodGlucoseTarget = buildBloodGlucoseTarget(originalBloodGlucoseUnits, bloodGlucoseTargetDatum, errs)
	}

	var recommended *Recommended
	recommendedDatum, ok := datum["recommended"].(map[string]interface{})
	if ok {
		recommended = buildRecommended(recommendedDatum, errs)
	}

	event := &Event{
		Recommended:        recommended,
		Bolus:              bolus,
		BloodGlucoseTarget: bloodGlucoseTarget,
		CarbohydrateInput:  datum.ToInt(carbohydrateInputField.Name, errs),
		InsulinOnBoard:     datum.ToFloat64(insulinOnBoardField.Name, errs),
		InsulinSensitivity: datum.ToInt(insulinSensitivityField.Name, errs),
		BloodGlucoseInput:  datum.ToFloat64(bloodGlucoseInputField.Name, errs),
		Units:              originalBloodGlucoseUnits,
		Base:               types.BuildBase(datum, errs),
	}

	event.BloodGlucoseInput, event.Units = types.NewBloodGlucoseValidation(event.BloodGlucoseInput, event.Units).SetValueAllowedToBeEmpty(true).ValidateAndConvertBloodGlucoseValue(errs)

	types.GetPlatformValidator().SetFailureReasons(failureReasons).Struct(event, errs)

	return event
}

func InsulinSensitivityValidator(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	_, ok := field.Interface().(float64)
	if !ok {
		return false
	}
	//TODO: correct validation here
	return true
}

func InsulinOnBoardValidator(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	iob, ok := field.Interface().(float64)
	if !ok {
		return false
	}
	return iob >= 0.0
}
