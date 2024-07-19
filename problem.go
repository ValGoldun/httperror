package httperror

import (
	"fmt"
	business_errors "github.com/ValGoldun/business-errors"
	"strings"
)

type Problem struct {
	Error    string                   `json:"error"`
	Fields   []Field                  `json:"fields,omitempty"`
	Metadata business_errors.Metadata `json:"metadata,omitempty"`
}

type Field struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

type Fields []Field

func (fields Fields) String() string {
	var formatted = make([]string, len(fields))

	for index, field := range fields {
		formatted[index] = fmt.Sprintf("%s: %s", field.Key, field.Error)
	}

	return strings.Join(formatted, ", ")
}
