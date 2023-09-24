package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/pkg/validator"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/server/http/api/response"
	"github.com/romankravchuk/eldorado/internal/services/auth/proto"
)

func HandleRegister(log *slog.Logger, client proto.AuthServiceClient) api.APIFunc {
	type req struct {
		Email     string `json:"email" validate:"required,email"`
		Username  string `json:"username" validate:"required,gte=3,lte=20"`
		Password  string `json:"password" validate:"required,gte=8,alphanum,lte=20"`
		FirstName string `json:"first_name" validate:"omitempty,alpha,gte=2"`
		LastName  string `json:"last_name" validate:"omitempty,alpha,gte=2"`
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
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

		ctx, cancel := context.WithTimeout(r.Context(), 150*time.Millisecond)
		defer cancel()

		resp, err := client.SignUp(ctx, &proto.SignUpRequest{
			Email:     input.Email,
			Username:  input.Username,
			Password:  input.Password,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		})
		if err != nil {
			msg := "internal server error"

			log.Error(msg, sl.Err(err), slog.Any("request_body", input))

			return response.APIError{
				Status:  http.StatusInternalServerError,
				Message: msg,
			}
		}
		if resp.Error != "" {
			msg := "invalid request"

			log.Error(msg, sl.Err(err), slog.Any("request_body", input), slog.Any("response", resp))

			return response.APIError{
				Status:  int(resp.Status),
				Message: msg,
			}
		}

		return response.JSON(w, http.StatusOK, response.M{"message": "ok"})
	}
}

func HandleGetTokenPairs(log *slog.Logger, client proto.AuthServiceClient) api.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return response.JSON(w, http.StatusOK, response.M{})
	}
}

func HandleRefreshToken(log *slog.Logger, client proto.AuthServiceClient) api.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return response.JSON(w, http.StatusOK, response.M{})
	}
}
