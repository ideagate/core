package context

import (
	"sync"
)

// ContextData data
type ContextData struct {
	sync.RWMutex
	Req  ContextRequestData         `json:",omitempty"` // data from http request
	Step map[string]ContextStepData `json:",omitempty"` // map[StepId]StepData
}

type ContextRequestData struct {
	Header map[string]any `json:",omitempty"`
	Query  map[string]any `json:",omitempty"` // map[queryVar]Value
	Json   map[string]any `json:",omitempty"` // map[jsonVar]Value
}

type ContextStepData struct {
	Var  map[string]any      `json:",omitempty"` // map[Var]Value. Data from step variables
	Data ContextStepDataBody `json:",omitempty"` // body response. For database in JSON form
	Out  map[string]any      `json:",omitempty"` // map[OutputVar]Value
}

type ContextStepDataBody struct {
	Body       any            `json:",omitempty"`
	Query      map[string]any `json:",omitempty"`
	StatusCode int            `json:"status_code"`
}

func (ctxData *ContextData) SetRequestQuery(query map[string]any) {
	ctxData.Lock()
	ctxData.Req.Query = query
	ctxData.Unlock()
}

func (ctxData *ContextData) SetRequestJson(json map[string]any) {
	ctxData.Lock()
	ctxData.Req.Json = json
	ctxData.Unlock()
}

func (ctxData *ContextData) GetStep(stepId string) ContextStepData {
	if ctxData.Step == nil {
		ctxData.Step = make(map[string]ContextStepData)
	}

	stepData, _ := ctxData.Step[stepId]
	return stepData
}

func (ctxData *ContextData) SetStepStatusCode(stepId string, statusCode int) {
	ctxData.Lock()
	stepData := ctxData.GetStep(stepId)
	stepData.Data.StatusCode = statusCode
	ctxData.Step[stepId] = stepData
	ctxData.Unlock()
}

func (ctxData *ContextData) SetStepDataBody(stepId string, body any) {
	ctxData.Lock()
	stepData := ctxData.GetStep(stepId)
	stepData.Data.Body = body
	ctxData.Step[stepId] = stepData
	ctxData.Unlock()
}

func (ctxData *ContextData) SetStepVariable(stepId string, data map[string]any) {
	ctxData.Lock()
	stepData := ctxData.GetStep(stepId)
	stepData.Var = data
	ctxData.Step[stepId] = stepData
	ctxData.Unlock()
}

func (ctxData *ContextData) SetStepOutput(stepId string, data map[string]any) {
	ctxData.Lock()
	stepData := ctxData.GetStep(stepId)
	stepData.Out = data
	ctxData.Step[stepId] = stepData
	ctxData.Unlock()
}
