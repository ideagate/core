package endpoint

import (
	"testing"

	entityContext "github.com/bayu-aditya/ideagate/backend/core/model/entity/context"
	pbEndpoint "github.com/bayu-aditya/ideagate/backend/model/gen-go/core/endpoint"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEndpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Endpoint Suite")
}

var _ = Describe("Variable", func() {
	mockStepId := "mockStepId"
	mockAnotherStep := "mockAnotherStep"
	mockAnotherStep2 := "mockAnotherStep2"

	mockCtxData := &entityContext.ContextData{
		Req: entityContext.ContextRequestData{
			Header: map[string]any{
				"header_1": "value_header_1",
				"header_2": "value_header_2",
			},
			Query: map[string]any{
				"query_1": "value_query_1",
				"query_2": 12345,
			},
			Json: map[string]any{
				"json_1": "value_json_1",
				"json_2": map[string]any{
					"json_2_a": 123,
				},
			},
		},
		Step: map[string]entityContext.ContextStepData{
			mockStepId: {
				Var: map[string]any{
					"var_1": "value_var_1",
					"var_2": "value_var_2",
				},
				Data: entityContext.ContextStepDataBody{
					Body: map[string]any{
						"current_body_1": "value_body_1",
						"current_body_2": true,
					},
					Query: map[string]any{
						"current_query_1": map[string]any{
							"col_a": "val_a_1",
							"col_b": "val_b_1",
						},
					},
					StatusCode: 200,
				},
			},
			mockAnotherStep: {
				Var: map[string]any{
					"var_3": "value_var_3",
					"var_4": 123.45,
				},
				Data: entityContext.ContextStepDataBody{
					Body: map[string]any{
						"body_1": "value_body_1",
						"body_2": true,
					},
					Query: map[string]any{
						"query_1": []any{
							map[string]any{"col_a": "val_a_1", "col_b": "val_b_1"},
							map[string]any{"col_a": "val_a_2", "col_b": "val_b_2"},
						},
					},
					StatusCode: 204,
				},
				Out: map[string]any{
					"out_1": "value_out_1",
				},
			},
			mockAnotherStep2: {},
		},
	}

	runTest := func(variable *Variable, wantResult any, wantErr bool) {
		got, err := variable.GetValue(mockStepId, mockCtxData)

		if wantResult == nil {
			Expect(got).To(BeNil())
		} else {
			Expect(got).To(Equal(wantResult))
		}

		if wantErr {
			Expect(err).To(HaveOccurred())
		} else {
			Expect(err).To(BeNil())
		}
	}

	Describe("From Request", func() {
		Context("Header", func() {
			It("{{.Req.Header.<Key>}} - exist", func() {
				variable := Variable{
					Value: "{{.Req.Header.header_1}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
				}
				got, err := variable.GetValue(mockStepId, mockCtxData)
				Expect(got).To(Equal("value_header_1"))
				Expect(err).To(BeNil())
			})
			It("{{.Req.Header.<Key>}} - not exist", func() {
				variable := Variable{
					Value:   "{{.Req.Header.unknown}}",
					Type:    pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					Default: "default_value",
				}
				got, err := variable.GetValue(mockStepId, mockCtxData)
				Expect(got).To(BeNil())
				Expect(err).To(BeNil())
			})
			It("{{.Req.Header.<Key>}} - using default", func() {
				variable := Variable{
					Value:    "{{.Req.Header.unknown}}",
					Type:     pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					Required: true,
					Default:  "default_value",
				}
				got, err := variable.GetValue(mockStepId, mockCtxData)
				Expect(got).To(Equal("default_value"))
				Expect(err).To(BeNil())
			})
		})
		Context("Query", func() {
			It("{{.Req.Query.<Key>}} - exist", func() {
				variable := Variable{
					Value: "{{.Req.Query.query_2}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_INT,
				}
				got, err := variable.GetValue(mockStepId, mockCtxData)
				Expect(got).To(Equal(int64(12345)))
				Expect(err).To(BeNil())
			})
			It("{{.Req.Query.<Key>}} - not exist", func() {
				variable := Variable{
					Value:   "{{.Req.Query.unknown}}",
					Type:    pbEndpoint.VariableType_VARIABLE_TYPE_INT,
					Default: "123",
				}
				got, err := variable.GetValue(mockStepId, mockCtxData)
				Expect(got).To(BeNil())
				Expect(err).To(BeNil())
			})
			It("{{.Req.Query.<Key>}} - using default", func() {
				variable := Variable{
					Value:    "{{.Req.Query.unknown}}",
					Type:     pbEndpoint.VariableType_VARIABLE_TYPE_INT,
					Required: true,
					Default:  "123",
				}
				got, err := variable.GetValue(mockStepId, mockCtxData)
				Expect(got).To(Equal(int64(123)))
				Expect(err).To(BeNil())
			})
		})
		Context("Json", func() {
			It("{{.Req.Json.<Key>}} - exist", func() {
				runTest(&Variable{
					Value: "{{.Req.Json.json_1}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
				}, "value_json_1", false)
			})
			It("{{.Req.Json.<Key>}} - exist nested", func() {
				runTest(&Variable{
					Value: "{{.Req.Json.json_2.json_2_a}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_INT,
				}, int64(123), false)
			})
			It("{{.Req.Json.<Key>}} - not exist", func() {
				runTest(&Variable{
					Value:   "{{.Req.Json.unknown.unknown}}",
					Type:    pbEndpoint.VariableType_VARIABLE_TYPE_INT,
					Default: "456",
				}, nil, false)
			})
			It("{{.Req.Json.<Key>}} - using default", func() {
				runTest(&Variable{
					Value:    "{{.Req.Json.unknown.unknown}}",
					Type:     pbEndpoint.VariableType_VARIABLE_TYPE_INT,
					Required: true,
					Default:  "456",
				}, int64(456), false)
			})
		})
	})
	Describe("From Another Step", func() {
		Context("Variable", func() {
			It("{{.Step.<StepId>.Var.<Key>}} - invalid step id", func() {
				runTest(&Variable{
					Value: "{{.Step.unknown.Var.var_1}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
				}, nil, false)
			})
			It("{{.Step.<StepId>.Var.<Key>}} - exist", func() {
				runTest(&Variable{
					Value: "{{.Step.mockAnotherStep.Var.var_4}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_FLOAT,
				}, 123.45, false)
			})
			It("{{.Step.<StepId>.Var.<Key>}} - not exist", func() {
				runTest(&Variable{
					Value:   "{{.Step.mockAnotherStep.Var.unknown}}",
					Type:    pbEndpoint.VariableType_VARIABLE_TYPE_FLOAT,
					Default: "789.123",
				}, nil, false)
			})
			It("{{.Step.<StepId>.Var.<Key>}} - using default", func() {
				runTest(&Variable{
					Value:    "{{.Step.mockAnotherStep.Var.unknown}}",
					Type:     pbEndpoint.VariableType_VARIABLE_TYPE_FLOAT,
					Required: true,
					Default:  "789.123",
				}, 789.123, false)
			})
		})
		Context("Data", func() {
			Context("Body", func() {
				It("{{.Step.<StepId>.Data.Body.<Key>}} - invalid step id", func() {
					runTest(&Variable{
						Value: "{{.Step.unknown.Data.Body.var_1}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, nil, false)
				})
				It("{{.Step.<StepId>.Data.Body.<Key>}} - exist", func() {
					runTest(&Variable{
						Value: "{{.Step.mockAnotherStep.Data.Body.body_2}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_BOOL,
					}, true, false)
				})
				It("{{.Step.<StepId>.Data.Body.<Key>}} - not exist", func() {
					runTest(&Variable{
						Value: "{{.Step.mockAnotherStep.Data.Body.unknown}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_BOOL,
					}, nil, false)
				})
				It("{{.Step.<StepId>.Data.Body.<Key>}} - using default", func() {
					runTest(&Variable{
						Value:    "{{.Step.mockAnotherStep.Data.Body.unknown}}",
						Type:     pbEndpoint.VariableType_VARIABLE_TYPE_BOOL,
						Required: true,
						Default:  "true",
					}, true, false)
				})
			})
			Context("Query", func() {
				It("{{.Step.<StepId>.Data.Query.<QueryId>}} - invalid step id", func() {
					runTest(&Variable{
						Value: "{{.Step.unknown.Data.Query.query_1}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, nil, false)
				})
				It("{{.Step.<StepId>.Data.Query.<QueryId>}} - invalid query id", func() {
					runTest(&Variable{
						Value: "{{(index .Step.mockAnotherStep.Data.Query.unknown 0).col_a}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, nil, false)
				})
				It("{{.Step.<StepId>.Data.Query.<QueryId>}} - not exist, index > length", func() {
					runTest(&Variable{
						Value: "{{(index .Step.mockAnotherStep.Data.Query.query_1 10).col_a}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, nil, false)
				})
				It("{{.Step.<StepId>.Data.Query.<QueryId>}} - invalid property", func() {
					runTest(&Variable{
						Value: "{{(index .Step.mockAnotherStep.Data.Query.query_1 0).unknown}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, nil, false)
				})
				It("{{.Step.<StepId>.Data.Query.<QueryId>}} - exist index 0", func() {
					runTest(&Variable{
						Value: "{{(index .Step.mockAnotherStep.Data.Query.query_1 0).col_a}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, "val_a_1", false)
				})
				It("{{.Step.<StepId>.Data.Query.<QueryId>}} - using default", func() {
					runTest(&Variable{
						Value:    "{{(index .Step.mockAnotherStep.Data.Query.unknown 0).col_a}}",
						Type:     pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
						Required: true,
						Default:  "default_value",
					}, "default_value", false)
				})
			})
			Context("StatusCode", func() {
				It("{{.Step.<StepId>.Data.StatusCode}} - invalid step id", func() {
					runTest(&Variable{
						Value: "{{.Step.unknown.Data.StatusCode}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_INT,
					}, nil, false)
				})
				It("{{.Step.<StepId>.Data.StatusCode}} - exist", func() {
					runTest(&Variable{
						Value: "{{.Step.mockAnotherStep.Data.StatusCode}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_INT,
					}, int64(204), false)
				})
				It("{{.Step.<StepId>.Data.StatusCode}} - empty", func() {
					runTest(&Variable{
						Value: "{{.Step.mockAnotherStep2.Data.StatusCode}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_INT,
					}, int64(0), false)
				})
				It("{{.Step.<StepId>.Data.StatusCode}} - using default", func() {
					runTest(&Variable{
						Value:    "{{.Step.mockAnotherStep2.Data.StatusCode}}",
						Type:     pbEndpoint.VariableType_VARIABLE_TYPE_INT,
						Required: true,
						Default:  "500",
					}, int64(500), false)
				})
			})
		})
		Context("Output", func() {
			It("{{.Step.<StepId>.Out.<Key>}} - invalid step id", func() {
				runTest(&Variable{
					Value: "{{.Step.unknown.Out.out_1}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
				}, nil, false)
			})
			It("{{.Step.<StepId>.Out.<Key>}} - exist", func() {
				runTest(&Variable{
					Value: "{{.Step.mockAnotherStep.Out.out_1}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
				}, "value_out_1", false)
			})
			It("{{.Step.<StepId>.Out.<Key>}} - not exist", func() {
				runTest(&Variable{
					Value: "{{.Step.mockAnotherStep.Out.unknown}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
				}, nil, false)
			})
			It("{{.Step.<StepId>.Out.<Key>}} - using default", func() {
				runTest(&Variable{
					Value:    "{{.Step.mockAnotherStep.Out.unknown}}",
					Type:     pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					Required: true,
					Default:  "default_value",
				}, "default_value", false)
			})
		})
	})
	Describe("From Current Step", func() {
		Context("Variable", func() {
			It("{{.Var.<Key>}} - exist", func() {
				runTest(&Variable{
					Value: "{{.Var.var_1}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
				}, "value_var_1", false)
			})
			It("{{.Var.<Key>}} - not exist", func() {
				runTest(&Variable{
					Value: "{{.Var.unknown}}",
					Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
				}, nil, false)
			})
			It("{{.Var.<Key>}} - using default", func() {
				runTest(&Variable{
					Value:    "{{.Var.unknown}}",
					Type:     pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					Required: true,
					Default:  "default_value",
				}, "default_value", false)
			})
		})
		Context("Data", func() {
			Context("Body", func() {
				It("{{.Data.Body.<Key>}} - exist", func() {
					runTest(&Variable{
						Value: "{{.Data.Body.current_body_1}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, "value_body_1", false)
				})
				It("{{.Data.Body.<Key>}} - not exist", func() {
					runTest(&Variable{
						Value: "{{.Data.Body.unknown}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, nil, false)
				})
				It("{{.Data.Body.<Key>}} - using default", func() {
					runTest(&Variable{
						Value:    "{{.Data.Body.unknown}}",
						Type:     pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
						Required: true,
						Default:  "default_value",
					}, "default_value", false)
				})
			})
			Context("Query", func() {
				It("{{.Data.Query.<QueryId>}} - exist", func() {
					runTest(&Variable{
						Value: "{{.Data.Query.current_query_1.col_a}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, "val_a_1", false)
				})
				It("{{.Data.Query.<QueryId>}} - not exist", func() {
					runTest(&Variable{
						Value: "{{.Data.Query.unknown.unknown}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
					}, nil, false)
				})
				It("{{.Data.Query.<QueryId>}} - using default", func() {
					runTest(&Variable{
						Value:    "{{.Data.Query.unknown.unknown}}",
						Type:     pbEndpoint.VariableType_VARIABLE_TYPE_STRING,
						Required: true,
						Default:  "default_value",
					}, "default_value", false)
				})
			})
			Context("StatusCode", func() {
				It("{{.Data.StatusCode}} - exist", func() {
					runTest(&Variable{
						Value: "{{.Data.StatusCode}}",
						Type:  pbEndpoint.VariableType_VARIABLE_TYPE_INT,
					}, int64(200), false)
				})
			})
		})
	})
})

func TestVariable_isEmptyValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil",
			args: args{
				value: nil,
			},
			want: true,
		},
		{
			name: "int64 - empty",
			args: args{
				value: int64(0),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Variable{}
			if got := v.isEmptyValue(tt.args.value); got != tt.want {
				t.Errorf("isEmptyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
