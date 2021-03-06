package status_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	testData "github.com/tidepool-org/platform/data/test"
	"github.com/tidepool-org/platform/data/types/device"
	"github.com/tidepool-org/platform/service"
)

func NewRawObject() map[string]interface{} {
	rawObject := testData.RawBaseObject()
	rawObject["type"] = "deviceEvent"
	rawObject["subType"] = "status"
	rawObject["duration"] = 0
	rawObject["status"] = "suspended"
	rawObject["reason"] = map[string]interface{}{"suspended": "manual"}
	return rawObject
}

func NewMeta() interface{} {
	return &device.Meta{
		Type:    "deviceEvent",
		SubType: "status",
	}
}

var _ = Describe("Status", func() {
	Context("duration", func() {
		DescribeTable("invalid when", testData.ExpectFieldNotValid,
			Entry("is less than 0", NewRawObject(), "duration", -1,
				[]*service.Error{testData.ComposeError(service.ErrorValueNotGreaterThanOrEqualTo(-1, 0), "/duration", NewMeta())},
			),
		)

		DescribeTable("valid when", testData.ExpectFieldIsValid,
			Entry("is 0", NewRawObject(), "duration", 0),
			Entry("is max of 999999999999999999", NewRawObject(), "duration", 999999999999999999),
		)
	})

	Context("status", func() {
		DescribeTable("invalid when", testData.ExpectFieldNotValid,
			Entry("is empty", NewRawObject(), "status", "",
				[]*service.Error{testData.ComposeError(service.ErrorValueStringNotOneOf("", []string{"resumed", "suspended"}), "/status", NewMeta())},
			),
			Entry("is not one of the predefined types", NewRawObject(), "status", "bad",
				[]*service.Error{testData.ComposeError(service.ErrorValueStringNotOneOf("bad", []string{"resumed", "suspended"}), "/status", NewMeta())},
			),
		)

		DescribeTable("valid when", testData.ExpectFieldIsValid,
			Entry("is suspended type", NewRawObject(), "status", "resumed"),
			Entry("is suspended type", NewRawObject(), "status", "suspended"),
		)
	})

	Context("reason", func() {
		DescribeTable("valid when", testData.ExpectFieldIsValid,
			Entry("is manual", NewRawObject(), "reason", map[string]interface{}{"suspended": "manual"}),
			Entry("is automatic", NewRawObject(), "reason", map[string]interface{}{"suspended": "automatic"}),
		)
	})
})
