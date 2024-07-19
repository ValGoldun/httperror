package httperror

import (
	"encoding/json"
	"errors"
	business_errors "github.com/ValGoldun/business-errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

func WriteProblem(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case business_errors.Error:
		if e.IsCritical {
			serverProblem(ctx, errors.New("internal error"), e.Metadata)
			return
		}
		clientProblem(ctx, errors.New(e.Text), e.Metadata)
	case *json.UnmarshalTypeError:
		clientProblem(ctx, errors.New("invalid json type"), nil)
		return
	case *json.SyntaxError:
		clientProblem(ctx, errors.New("invalid json"), nil)
		return
	case validator.ValidationErrors:
		var fields []Field
		for _, field := range e {
			fields = append(fields, Field{Key: field.Field(), Error: field.Tag()})
		}
		clientProblemWithFields(ctx, errors.New("validation error"), fields)
		return
	default:
		if errors.Is(err, io.EOF) {
			clientProblem(ctx, errors.New("empty body"), nil)
			return
		}
		serverProblem(ctx, err, nil)
		return
	}
}
func serverProblem(ctx *gin.Context, err error, metadata business_errors.Metadata) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Problem{Error: "server problem", Metadata: metadata})
}

func clientProblem(ctx *gin.Context, err error, metadata business_errors.Metadata) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, Problem{Error: err.Error(), Metadata: metadata})
}

func clientProblemWithFields(ctx *gin.Context, err error, fields []Field) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, Problem{Error: err.Error(), Fields: fields})
}
