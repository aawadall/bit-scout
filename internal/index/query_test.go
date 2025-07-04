package index

import (
	"testing"

	"github.com/aawadall/bit-scout/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestParseQuery_SimpleEquals(t *testing.T) {
	q, err := ParseQuery("filename=main.go")
	assert.NoError(t, err)
	assert.Len(t, q.Conditions, 1)
	assert.Equal(t, "filename", q.Conditions[0].Dimension)
	assert.Equal(t, OpEquals, q.Conditions[0].Operator)
	assert.Equal(t, "main.go", q.Conditions[0].Value)
}

func TestParseQuery_AndConditions(t *testing.T) {
	q, err := ParseQuery("filename=main.go and fileExtension=go")
	assert.NoError(t, err)
	assert.Len(t, q.Conditions, 2)
	assert.Equal(t, "filename", q.Conditions[0].Dimension)
	assert.Equal(t, "fileExtension", q.Conditions[1].Dimension)
}

func TestQueryCondition_Evaluate_Equals(t *testing.T) {
	doc := models.Document{Meta: map[string]string{"filename": "main.go"}}
	cond := QueryCondition{Dimension: "filename", Operator: OpEquals, Value: "main.go"}
	match, err := cond.Evaluate(doc)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestQueryCondition_Evaluate_NotEquals(t *testing.T) {
	doc := models.Document{Meta: map[string]string{"filename": "main.go"}}
	cond := QueryCondition{Dimension: "filename", Operator: OpNotEquals, Value: "other.go"}
	match, err := cond.Evaluate(doc)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestQueryCondition_Evaluate_Contains(t *testing.T) {
	doc := models.Document{Meta: map[string]string{"filename": "main.go"}}
	cond := QueryCondition{Dimension: "filename", Operator: OpContains, Value: "main"}
	match, err := cond.Evaluate(doc)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestQueryCondition_Evaluate_Numeric(t *testing.T) {
	doc := models.Document{Meta: map[string]string{"fileSize": "100"}}
	cond := QueryCondition{Dimension: "fileSize", Operator: OpGreater, Value: "10"}
	match, err := cond.Evaluate(doc)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestQuery_Evaluate_AllConditionsMatch(t *testing.T) {
	doc := models.Document{Meta: map[string]string{"filename": "main.go", "fileExtension": "go"}}
	q, _ := ParseQuery("filename=main.go and fileExtension=go")
	match, err := q.Evaluate(doc)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestQuery_Evaluate_ConditionFails(t *testing.T) {
	doc := models.Document{Meta: map[string]string{"filename": "main.go", "fileExtension": "py"}}
	q, _ := ParseQuery("filename=main.go and fileExtension=go")
	match, err := q.Evaluate(doc)
	assert.NoError(t, err)
	assert.False(t, match)
}

func TestParseQuery_InvalidFormat(t *testing.T) {
	_, err := ParseQuery("invalidquery")
	assert.Error(t, err)
}
