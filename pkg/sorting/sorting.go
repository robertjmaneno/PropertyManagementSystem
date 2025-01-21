package sorting

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Direction represents sort direction
type Direction string

const (
	ASC  Direction = "asc"
	DESC Direction = "desc"
)

// Field represents a sortable field
type Field struct {
	Name          string    `json:"name"`
	Direction     Direction `json:"direction"`
	CaseSensitive bool      `json:"case_sensitive,omitempty"`
}

// Options represents sorting options
type Options struct {
	Fields []Field `json:"sort" form:"sort"`
}

// NewOptions creates a new Options instance with default sorting
func NewOptions(defaultField string, defaultDirection Direction) *Options {
	return &Options{
		Fields: []Field{
			{
				Name:          defaultField,
				Direction:     defaultDirection,
				CaseSensitive: false,
			},
		},
	}
}

// AddField adds a sort field
func (o *Options) AddField(name string, direction Direction, caseSensitive bool) {
	o.Fields = append(o.Fields, Field{
		Name:          name,
		Direction:     direction,
		CaseSensitive: caseSensitive,
	})
}

// ParseFromString parses sort options from string
// Format: "field1:asc:cs,field2:desc" (cs = case sensitive)
func (o *Options) ParseFromString(s string) error {
	if s == "" {
		return nil
	}

	fields := strings.Split(s, ",")
	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) < 2 || len(parts) > 3 {
			return fmt.Errorf("invalid sort format: %s", field)
		}

		name := strings.TrimSpace(parts[0])
		direction := strings.ToLower(strings.TrimSpace(parts[1]))

		if direction != string(ASC) && direction != string(DESC) {
			return fmt.Errorf("invalid sort direction: %s", direction)
		}

		// Check for case sensitivity flag
		caseSensitive := false
		if len(parts) == 3 && strings.ToLower(strings.TrimSpace(parts[2])) == "cs" {
			caseSensitive = true
		}

		o.AddField(name, Direction(direction), caseSensitive)
	}

	return nil
}

// ValidateFields validates sort fields against allowed fields
func (o *Options) ValidateFields(allowedFields map[string]bool) error {
	for _, field := range o.Fields {
		if !allowedFields[field.Name] {
			return fmt.Errorf("invalid sort field: %s", field.Name)
		}
	}
	return nil
}

// Apply applies sorting to a GORM query
func (o *Options) Apply(db *gorm.DB) *gorm.DB {
	if len(o.Fields) == 0 {
		return db
	}

	var orders []string
	for _, field := range o.Fields {
		orderExpr := field.Name
		if !field.CaseSensitive {
			// For case-insensitive sorting, use LOWER() function
			orderExpr = fmt.Sprintf("LOWER(%s)", field.Name)
		}
		orders = append(orders, fmt.Sprintf("%s %s", orderExpr, strings.ToUpper(string(field.Direction))))
	}

	return db.Order(strings.Join(orders, ", "))
}

// BuildSearchQuery builds a case-sensitive or case-insensitive search query
func BuildSearchQuery(query string, caseSensitive bool, fields ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if query == "" || len(fields) == 0 {
			return db
		}

		q := db
		for i, field := range fields {
			var condition string
			var searchValue string

			if caseSensitive {
				condition = field + " LIKE ?"
				searchValue = "%" + query + "%"
			} else {
				condition = "LOWER(" + field + ") LIKE ?"
				searchValue = "%" + strings.ToLower(query) + "%"
			}

			if i == 0 {
				q = q.Where(condition, searchValue)
			} else {
				q = q.Or(condition, searchValue)
			}
		}
		return q
	}
} 