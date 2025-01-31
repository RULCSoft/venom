package venom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestProcessVariableAssigments(t *testing.T) {

	assign := AssignStep{}
	assign.Assignments = make(map[string]Assignment)
	assign.Assignments["assignVar"] = Assignment{
		From: "here.some.value",
	}
	assign.Assignments["assignVarWithRegex"] = Assignment{
		From:  "here.some.value",
		Regex: `this is (?s:(.*))`,
	}

	b, _ := yaml.Marshal(assign)
	t.Log("\n" + string(b))

	var stepIn TestStep
	assert.NoError(t, yaml.Unmarshal(b, &stepIn))

	tcVars := H{
		"here.some.value": "this is the \nvalue",
	}

	result, is, err := ProcessVariableAssigments("", tcVars, stepIn, TestLogger{t})
	assert.True(t, is)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	t.Log(result)
	assert.Equal(t, "map[assignVar:this is the \nvalue assignVarWithRegex:the \nvalue]", fmt.Sprint(result))

	var wrongStepIn TestStep
	b = []byte(`type: exec
script: echo 'foo'
`)
	assert.NoError(t, yaml.Unmarshal(b, &wrongStepIn))
	result, is, err = ProcessVariableAssigments("", tcVars, wrongStepIn, TestLogger{t})
	assert.False(t, is)
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.Empty(t, result)
}
