package httperror

import business_errors "github.com/ValGoldun/business-errors"

type Problem struct {
	Error    string                   `json:"error"`
	Fields   []Field                  `json:"fields,omitempty"`
	Metadata business_errors.Metadata `json:"metadata,omitempty"`
}

type Field struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}
