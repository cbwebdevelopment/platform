package continuous_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/tidepool-org/platform/data/context"
	"github.com/tidepool-org/platform/data/normalizer"
	"github.com/tidepool-org/platform/data/types/base"
	"github.com/tidepool-org/platform/data/types/base/continuous"
	"github.com/tidepool-org/platform/data/types/base/testing"
	"github.com/tidepool-org/platform/data/types/common/bloodglucose"
	"github.com/tidepool-org/platform/data/validator"
	"github.com/tidepool-org/platform/service"
)

var _ = Describe("Continuous", func() {
	var rawObject = testing.RawBaseObject()
	var meta = &base.Meta{
		Type: "cbg",
	}

	rawObject["type"] = "cbg"
	rawObject["units"] = "mmol/L"
	rawObject["value"] = 5

	Context("units", func() {
		DescribeTable("units when", testing.ExpectFieldNotValid,
			Entry("is empty", rawObject, "units", "",
				[]*service.Error{testing.ComposeError(validator.ErrorStringNotOneOf("", []string{bloodglucose.Mmoll, bloodglucose.MmolL, bloodglucose.Mgdl, bloodglucose.MgdL}), "/units", meta)},
			),
			Entry("is not one of the predefined values", rawObject, "units", "wrong",
				[]*service.Error{testing.ComposeError(validator.ErrorStringNotOneOf("wrong", []string{bloodglucose.Mmoll, bloodglucose.MmolL, bloodglucose.Mgdl, bloodglucose.MgdL}), "/units", meta)},
			),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("is mmol/l", rawObject, "units", "mmol/l"),
			Entry("is mmol/L", rawObject, "units", "mmol/L"),
			Entry("is mg/dl", rawObject, "units", "mg/dl"),
			Entry("is mg/dL", rawObject, "units", "mg/dL"),
		)
	})

	Context("value", func() {
		DescribeTable("value when", testing.ExpectFieldNotValid,
			Entry("is less than 0", rawObject, "value", -0.1,
				[]*service.Error{testing.ComposeError(validator.ErrorFloatNotInRange(-0.1, bloodglucose.MgdLFromValue, bloodglucose.MgdLToValue), "/value", meta)},
			),
			Entry("is greater than 1000", rawObject, "value", 1000.1,
				[]*service.Error{testing.ComposeError(validator.ErrorFloatNotInRange(1000.1, bloodglucose.MgdLFromValue, bloodglucose.MgdLToValue), "/value", meta)},
			),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("is above 0", rawObject, "value", 0.1),
			Entry("is below max", rawObject, "value", 990.85745),
			Entry("is an integer", rawObject, "value", 380),
		)
	})

	Context("normalized when mmol/L", func() {
		DescribeTable("normalization", func(val, expected float64) {
			continuousBg, err := continuous.New()
			Expect(err).To(BeNil())
			continuousBg.Value = &val
			continuousBg.Units = &bloodglucose.Mmoll

			testContext := context.NewStandard()
			standardNormalizer, err := normalizer.NewStandard(testContext)
			Expect(err).To(BeNil())
			continuousBg.Normalize(standardNormalizer)
			Expect(continuousBg.Units).To(Equal(&bloodglucose.MmolL))
			Expect(continuousBg.Value).To(Equal(&expected))
		},
			Entry("is expected lower bg value", 3.7, 3.7),
			Entry("is below max", 54.99, 54.99),
			Entry("is expected upper bg value", 23.0, 23.0),
		)
	})

	Context("normalized when mg/dL", func() {
		DescribeTable("normalization", func(val, expected float64) {
			continuousBg, err := continuous.New()
			Expect(err).To(BeNil())
			continuousBg.Value = &val
			continuousBg.Units = &bloodglucose.Mgdl

			testContext := context.NewStandard()
			standardNormalizer, err := normalizer.NewStandard(testContext)
			Expect(err).To(BeNil())
			continuousBg.Normalize(standardNormalizer)
			Expect(continuousBg.Units).To(Equal(&bloodglucose.MmolL))
			Expect(continuousBg.Value).To(Equal(&expected))
		},
			Entry("is expected lower bg value", 60.0, 3.33044879462732),
			Entry("is below max", 990.85745, 55.0),
			Entry("is expected upper bg value", 400.0, 22.202991964182132),
		)
	})
})