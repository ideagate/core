package endpoint

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"

	entityContext "github.com/bayu-aditya/ideagate/backend/core/model/entity/context"
	pbEndpoint "github.com/bayu-aditya/ideagate/backend/model/gen-go/core/endpoint"
	"github.com/spf13/cast"
)

type Variable pbEndpoint.Variable

func (v *Variable) GetValue(stepId string, ctxData *entityContext.ContextData) (interface{}, error) {
	var result any = v.Value

	// get value from context
	result = v.getValueFromTemplate(stepId, ctxData, v.Value)

	// parse value by type
	result, err := v.parseValueByType(result, v.Type)
	if err != nil {
		return nil, err
	}

	// check is value empty
	if v.isEmptyValue(result) && v.Required {
		// set into default and parse the default value
		result, err = v.parseValueByType(v.Default, v.Type)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (v *Variable) GetValueString(stepId string, ctxData *entityContext.ContextData) (string, error) {
	value, err := v.GetValue(stepId, ctxData)
	if err != nil {
		return "", err
	}

	return cast.ToStringE(value)
}

func (v *Variable) getValueFromTemplate(stepId string, ctxData *entityContext.ContextData, templateValue string) interface{} {
	tmpl, err := template.New("").Parse(templateValue)

	if err != nil {
		return nil
	}

	type dataTemplateType struct {
		Req  entityContext.ContextRequestData
		Step map[string]entityContext.ContextStepData
		Var  map[string]any
		Data entityContext.ContextStepDataBody
	}

	data := dataTemplateType{
		Req:  ctxData.Req,
		Step: ctxData.Step,
		Var:  ctxData.Step[stepId].Var,
		Data: ctxData.Step[stepId].Data,
	}

	var resultBuffer bytes.Buffer
	if err = tmpl.Execute(&resultBuffer, data); err != nil {
		return nil
	}

	result := resultBuffer.String()
	result = strings.ReplaceAll(result, "<no value>", "")

	if result == "" {
		return nil
	}

	return result
}

func (v *Variable) parseValueByType(value interface{}, varType pbEndpoint.VariableType) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	switch varType {
	case pbEndpoint.VariableType_VARIABLE_TYPE_STRING:
		return cast.ToStringE(value)

	case pbEndpoint.VariableType_VARIABLE_TYPE_INT:
		return cast.ToInt64E(value)

	case pbEndpoint.VariableType_VARIABLE_TYPE_FLOAT:
		return cast.ToFloat64E(value)

	case pbEndpoint.VariableType_VARIABLE_TYPE_BOOL:
		return cast.ToBoolE(value)

	case pbEndpoint.VariableType_VARIABLE_TYPE_OBJECT:
		return value, nil
	}

	return value, nil
}

func (v *Variable) isEmptyValue(value interface{}) bool {
	if value == nil {
		return true
	}

	return reflect.ValueOf(value).IsZero()
}
