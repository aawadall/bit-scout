package index

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aawadall/bit-scout/internal/models"
	"github.com/rs/zerolog/log"
)

// QueryOperator represents comparison operators
type QueryOperator string

const (
	OpEquals    QueryOperator = "="
	OpNotEquals QueryOperator = "!="
	OpLess      QueryOperator = "<"
	OpLessEq    QueryOperator = "<="
	OpGreater   QueryOperator = ">"
	OpGreaterEq QueryOperator = ">="
	OpContains  QueryOperator = "contains"
)

// QueryCondition represents a single condition in a query
type QueryCondition struct {
	Dimension string
	Operator  QueryOperator
	Value     string
}

// Query represents a parsed query with conditions
type Query struct {
	Conditions []QueryCondition
	RawQuery   string
}

// ParseQuery parses a query string into a Query struct
func ParseQuery(queryStr string) (*Query, error) {
	query := &Query{
		RawQuery:   queryStr,
		Conditions: []QueryCondition{},
	}

	// Split by AND/OR operators (for now, we'll treat everything as AND)
	// This is a simple implementation - can be extended for OR logic
	parts := strings.Split(queryStr, " and ")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		condition, err := parseCondition(part)
		if err != nil {
			return nil, fmt.Errorf("failed to parse condition '%s': %w", part, err)
		}

		query.Conditions = append(query.Conditions, condition)
	}

	log.Debug().Msgf("Parsed query '%s' into %d conditions", queryStr, len(query.Conditions))
	return query, nil
}

// parseCondition parses a single condition like "fileExtension=go" or "fileSize<10"
func parseCondition(conditionStr string) (QueryCondition, error) {
	// Regex to match: dimension operator value
	// Supports: =, !=, <, <=, >, >=, contains
	re := regexp.MustCompile(`^(\w+)\s*(=|!=|<=|>=|<|>|contains)\s*(.+)$`)
	matches := re.FindStringSubmatch(conditionStr)

	if len(matches) != 4 {
		return QueryCondition{}, fmt.Errorf("invalid condition format: %s", conditionStr)
	}

	dimension := matches[1]
	operator := QueryOperator(matches[2])
	value := strings.TrimSpace(matches[3])

	// Remove quotes if present
	if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
		(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
		value = value[1 : len(value)-1]
	}

	return QueryCondition{
		Dimension: dimension,
		Operator:  operator,
		Value:     value,
	}, nil
}

// Evaluate evaluates a query against a document
func (q *Query) Evaluate(doc models.Document) (bool, error) {
	for _, condition := range q.Conditions {
		matches, err := condition.Evaluate(doc)
		if err != nil {
			return false, fmt.Errorf("condition evaluation failed: %w", err)
		}

		if !matches {
			return false, nil // AND logic - if any condition fails, document doesn't match
		}
	}

	return true, nil
}

// Evaluate evaluates a single condition against a document
func (c *QueryCondition) Evaluate(doc models.Document) (bool, error) {
	// Get the value from document metadata
	docValue, exists := doc.Meta[c.Dimension]
	if !exists {
		// If dimension doesn't exist, try to get from document properties
		switch c.Dimension {
		case "filename":
			docValue = doc.Meta["filename"]
		case "path":
			docValue = doc.Source
		case "text":
			docValue = doc.Text
		default:
			return false, nil // Dimension not found, condition fails
		}
	}

	if docValue == "" {
		return false, nil
	}

	switch c.Operator {
	case OpEquals:
		return strings.EqualFold(docValue, c.Value), nil

	case OpNotEquals:
		return !strings.EqualFold(docValue, c.Value), nil

	case OpContains:
		return strings.Contains(strings.ToLower(docValue), strings.ToLower(c.Value)), nil

	case OpLess, OpLessEq, OpGreater, OpGreaterEq:
		// Try to convert to numeric comparison
		return c.evaluateNumeric(docValue)

	default:
		return false, fmt.Errorf("unsupported operator: %s", c.Operator)
	}
}

// evaluateNumeric handles numeric comparisons
func (c *QueryCondition) evaluateNumeric(docValue string) (bool, error) {
	// Try to parse as float64 for numeric comparison
	docNum, err := strconv.ParseFloat(docValue, 64)
	if err != nil {
		// If not numeric, fall back to string comparison
		switch c.Operator {
		case OpLess:
			return docValue < c.Value, nil
		case OpLessEq:
			return docValue <= c.Value, nil
		case OpGreater:
			return docValue > c.Value, nil
		case OpGreaterEq:
			return docValue >= c.Value, nil
		default:
			return false, fmt.Errorf("unsupported numeric operator: %s", c.Operator)
		}
	}

	queryNum, err := strconv.ParseFloat(c.Value, 64)
	if err != nil {
		return false, fmt.Errorf("query value '%s' is not numeric", c.Value)
	}

	switch c.Operator {
	case OpLess:
		return docNum < queryNum, nil
	case OpLessEq:
		return docNum <= queryNum, nil
	case OpGreater:
		return docNum > queryNum, nil
	case OpGreaterEq:
		return docNum >= queryNum, nil
	default:
		return false, fmt.Errorf("unsupported numeric operator: %s", c.Operator)
	}
}
