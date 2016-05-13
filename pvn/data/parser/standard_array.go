package parser

/* CHECKLIST
 * [ ] Uses interfaces as appropriate
 * [ ] Private package variables use underscore prefix
 * [ ] All parameters validated
 * [ ] All errors handled
 * [ ] Reviewed for concurrency safety
 * [ ] Code complete
 * [ ] Full test coverage
 */

import (
	"github.com/tidepool-org/platform/app"
	"github.com/tidepool-org/platform/pvn/data"
)

type StandardArray struct {
	context data.Context
	array   *[]interface{}
}

func NewStandardArray(context data.Context, array *[]interface{}) (*StandardArray, error) {
	if context == nil {
		return nil, app.Error("parser", "context is missing")
	}

	return &StandardArray{
		context: context,
		array:   array,
	}, nil
}

func (s *StandardArray) Array() *[]interface{} {
	return s.array
}

func (s *StandardArray) ParseBoolean(index int) *bool {
	if s.array == nil {
		return nil
	}

	if index < 0 || index >= len(*s.array) {
		return nil
	}

	rawValue := (*s.array)[index]

	booleanValue, ok := rawValue.(bool)
	if !ok {
		s.context.AppendError(index, ErrorTypeNotBoolean(rawValue))
		return nil
	}

	return &booleanValue
}

func (s *StandardArray) ParseInteger(index int) *int {
	if s.array == nil {
		return nil
	}

	if index < 0 || index >= len(*s.array) {
		return nil
	}

	rawValue := (*s.array)[index]

	integerValue, integerValueOk := rawValue.(int)
	if !integerValueOk {
		floatValue, floatValueOk := rawValue.(float64)
		if !floatValueOk {
			s.context.AppendError(index, ErrorTypeNotInteger(rawValue))
			return nil
		}
		integerValue = int(floatValue)
	}

	return &integerValue
}

func (s *StandardArray) ParseFloat(index int) *float64 {
	if s.array == nil {
		return nil
	}

	if index < 0 || index >= len(*s.array) {
		return nil
	}

	rawValue := (*s.array)[index]

	floatValue, floatValueOk := rawValue.(float64)
	if !floatValueOk {
		integerValue, integerValueOk := rawValue.(int)
		if !integerValueOk {
			s.context.AppendError(index, ErrorTypeNotFloat(rawValue))
			return nil
		}
		floatValue = float64(integerValue)
	}

	return &floatValue
}

func (s *StandardArray) ParseString(index int) *string {
	if s.array == nil {
		return nil
	}

	if index < 0 || index >= len(*s.array) {
		return nil
	}

	rawValue := (*s.array)[index]

	stringValue, ok := rawValue.(string)
	if !ok {
		s.context.AppendError(index, ErrorTypeNotString(rawValue))
		return nil
	}

	return &stringValue
}

func (s *StandardArray) ParseStringArray(index int) *[]string {
	if s.array == nil {
		return nil
	}

	if index < 0 || index >= len(*s.array) {
		return nil
	}

	rawValue := (*s.array)[index]

	stringArrayValue, stringArrayValueOk := rawValue.([]string)
	if !stringArrayValueOk {
		arrayValue, arrayValueOk := rawValue.([]interface{})
		if !arrayValueOk {
			s.context.AppendError(index, ErrorTypeNotArray(rawValue))
			return nil
		}

		stringArrayValue = []string{}
		parser, _ := NewStandardArray(s.context.NewChildContext(index), &arrayValue)
		for index := range arrayValue {
			var stringElement string
			if stringParsed := parser.ParseString(index); stringParsed != nil {
				stringElement = *stringParsed
			}
			stringArrayValue = append(stringArrayValue, stringElement)
		}
	}

	return &stringArrayValue
}

func (s *StandardArray) ParseObject(index int) *map[string]interface{} {
	if s.array == nil {
		return nil
	}

	if index < 0 || index >= len(*s.array) {
		return nil
	}

	rawValue := (*s.array)[index]

	objectValue, ok := rawValue.(map[string]interface{})
	if !ok {
		s.context.AppendError(index, ErrorTypeNotObject(rawValue))
		return nil
	}

	return &objectValue
}

func (s *StandardArray) ParseObjectArray(index int) *[]map[string]interface{} {
	if s.array == nil {
		return nil
	}

	if index < 0 || index >= len(*s.array) {
		return nil
	}

	rawValue := (*s.array)[index]

	objectArrayValue, objectArrayValueOk := rawValue.([]map[string]interface{})
	if !objectArrayValueOk {
		arrayValue, arrayValueOk := rawValue.([]interface{})
		if !arrayValueOk {
			s.context.AppendError(index, ErrorTypeNotArray(rawValue))
			return nil
		}

		parser, _ := NewStandardArray(s.context.NewChildContext(index), &arrayValue)
		for index := range arrayValue {
			var objectElement map[string]interface{}
			if objectParsed := parser.ParseObject(index); objectParsed != nil {
				objectElement = *objectParsed
			}
			objectArrayValue = append(objectArrayValue, objectElement)
		}
	}

	return &objectArrayValue
}

func (s *StandardArray) ParseInterface(index int) *interface{} {
	if s.array == nil {
		return nil
	}

	if index < 0 || index >= len(*s.array) {
		return nil
	}

	rawValue := (*s.array)[index]

	return &rawValue
}

func (s *StandardArray) ParseInterfaceArray(index int) *[]interface{} {
	if s.array == nil {
		return nil
	}

	if index < 0 || index >= len(*s.array) {
		return nil
	}

	rawValue := (*s.array)[index]

	arrayValue, ok := rawValue.([]interface{})
	if !ok {
		s.context.AppendError(index, ErrorTypeNotArray(rawValue))
		return nil
	}

	return &arrayValue
}

func (s *StandardArray) NewChildObjectParser(index int) data.ObjectParser {
	standardObject, _ := NewStandardObject(s.context.NewChildContext(index), s.ParseObject(index))
	return standardObject
}

func (s *StandardArray) NewChildArrayParser(index int) data.ArrayParser {
	standardArray, _ := NewStandardArray(s.context.NewChildContext(index), s.ParseInterfaceArray(index))
	return standardArray
}