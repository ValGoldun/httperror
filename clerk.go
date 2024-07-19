package httperror

import (
	"encoding/json"
	"errors"
	business_errors "github.com/ValGoldun/business-errors"
	"github.com/ValGoldun/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

type ProblemWriter struct {
	logger logger.Logger
}

func (pw ProblemWriter) Problem(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case business_errors.Error:
		if e.IsCritical {
			pw.serverProblem(ctx, errors.New("internal error"), e.Metadata)
			return
		}
		pw.clientProblem(ctx, errors.New(e.Text), e.Metadata)
	case *json.UnmarshalTypeError:
		pw.clientProblem(ctx, errors.New("invalid json type"), nil)
		return
	case *json.SyntaxError:
		pw.clientProblem(ctx, errors.New("invalid json"), nil)
		return
	case validator.ValidationErrors:
		var fields = make(Fields, len(e))
		for index, field := range e {
			fields[index] = Field{Key: field.Field(), Error: field.Tag()}
		}
		pw.clientProblemWithFields(ctx, errors.New("validation error"), fields)
		return
	default:
		if errors.Is(err, io.EOF) {
			pw.clientProblem(ctx, errors.New("empty body"), nil)
			return
		}
		pw.serverProblem(ctx, err, nil)
		return
	}
}
func (pw ProblemWriter) serverProblem(ctx *gin.Context, err error, metadata business_errors.Metadata) {
	pw.logger.Error(err.Error(), metadata.LoggerFields()...)

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Problem{Error: "server problem", Metadata: metadata})
}

func (pw ProblemWriter) clientProblem(ctx *gin.Context, err error, metadata business_errors.Metadata) {
	pw.logger.Warn(err.Error(), metadata.LoggerFields()...)

	ctx.AbortWithStatusJSON(http.StatusBadRequest, Problem{Error: err.Error(), Metadata: metadata})
}

func (pw ProblemWriter) clientProblemWithFields(ctx *gin.Context, err error, fields Fields) {
	pw.logger.Warn(err.Error(), business_errors.Metadata{
		"validation_error": fields.String(),
	}.LoggerFields()...)

	ctx.AbortWithStatusJSON(http.StatusBadRequest, Problem{Error: err.Error(), Fields: fields})
}
