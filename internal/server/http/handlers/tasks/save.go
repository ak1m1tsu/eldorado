package tasks

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/romankravchuk/eldorado/internal/data"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/pkg/validator"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
)

type TaskCreater interface {
	Create(ctx context.Context, userID string, task data.Task) (data.Task, error)
}

func HandleCreateTask(log *slog.Logger, creater TaskCreater) api.APIFunc {
	const op = "server.http.handlers.CreateTask"

	type req struct {
		Title       string `json:"title" validate:"required,min=3,max=100"`
		Description string `json:"description" validate:"required,min=3,max=500"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", r.Header.Get(api.RequestIDHeader)),
		)

		input := new(req)
		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusBadRequest,
				Message: msg,
			}
		}

		if err := validator.ValidateStruct(*input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			}
		}

		userID, ok := r.Context().Value(api.UserIDKey).(string)
		if !ok {
			msg := "forbidden"

			log.Error(msg, slog.Any("input", input))

			return response.APIError{
				Status:  http.StatusForbidden,
				Message: msg,
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 150*time.Millisecond)
		defer cancel()

		t, err := creater.Create(
			ctx,
			userID,
			data.Task{Title: input.Title, Description: input.Description},
		)
		if err != nil {
			msg := "internal server error"

			log.Error(msg, sl.Err(err))

			return response.APIError{
				Status:  http.StatusInternalServerError,
				Message: msg,
			}
		}

		return response.JSON(w, http.StatusCreated, response.M{
			"task": t,
		})
	}
}
