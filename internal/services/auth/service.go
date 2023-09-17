package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/romankravchuk/eldorado/internal/data"
	"github.com/romankravchuk/eldorado/internal/pkg/jwt"
	"github.com/romankravchuk/eldorado/internal/pkg/logger"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/pkg/validator"
	"github.com/romankravchuk/eldorado/internal/services"
	"github.com/romankravchuk/eldorado/internal/services/auth/proto"
	"github.com/romankravchuk/eldorado/internal/storages"
	"github.com/romankravchuk/eldorado/internal/storages/users"
	"github.com/romankravchuk/eldorado/internal/storages/users/pg"
	"golang.org/x/crypto/bcrypt"
)

type Option func(*Service) error

func WithUsersStorage(users users.Storage) Option {
	return func(s *Service) error {
		if users == nil {
			return services.ErrNilUsersStorage
		}

		s.users = users
		return nil
	}
}

func WithUsersPostgresStorage(url string) Option {
	return func(s *Service) error {
		pool, err := storages.NewDBPool("postgres", url)
		if err != nil {
			return err
		}

		users, err := pg.New(pool)
		if err != nil {
			return err
		}

		return WithUsersStorage(users)(s)
	}
}

func WithAccessCreds(pvKey, pbKey string, ttl time.Duration) Option {
	return func(s *Service) error {
		pem, pub, err := decodeRSAKeys(pbKey, pbKey)
		if err != nil {
			return err
		}

		s.access = data.RSACredentials{
			PrivateKey: pem,
			PublicKey:  pub,
			TTL:        ttl,
		}
		return nil
	}
}

func WithRefreshCreds(pvKey, pbKey string, ttl time.Duration) Option {
	return func(s *Service) error {
		pem, pub, err := decodeRSAKeys(pbKey, pbKey)
		if err != nil {
			return err
		}

		s.refresh = data.RSACredentials{
			PrivateKey: pem,
			PublicKey:  pub,
			TTL:        ttl,
		}
		return nil
	}
}

func WithLogger(log *slog.Logger) Option {
	return func(s *Service) error {
		if log == nil {
			log = logger.New("local", os.Stderr)
		}

		s.log = log
		return nil
	}
}

type Service struct {
	users users.Storage

	log *slog.Logger

	access  data.RSACredentials
	refresh data.RSACredentials

	proto.UnsafeAuthServiceServer
}

func New(opts ...Option) (*Service, error) {
	s := &Service{}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Service) SignUp(ctx context.Context, in *proto.SignUpRequest) (*proto.Response, error) {
	const op = "services.auth.SignUp"

	log := s.log.With("op", op)

	if err := validator.ValidateStruct(in); err != nil {
		msg := "failed to validate request"

		log.Error(msg, sl.Err(err))

		return &proto.Response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		msg := "failed to generate hash from password"

		log.Error(msg, sl.Err(err))

		return &proto.Response{
			Status: http.StatusInternalServerError,
			Error:  msg,
		}, nil
	}

	err = s.users.Save(ctx, &data.User{
		Email:             in.Email,
		Username:          in.Username,
		Name:              in.Username,
		EncryptedPassword: string(hash),
	})
	if err != nil {
		if errors.Is(err, users.ErrAlreadyExists) {
			msg := "the user with given email already exists"

			log.Error(msg, sl.Err(err))

			return &proto.Response{
				Status: http.StatusBadRequest,
				Error:  msg,
			}, nil
		}

		msg := "failed to create user"

		log.Error(msg, sl.Err(err))

		return &proto.Response{
			Status: http.StatusInternalServerError,
			Error:  msg,
		}, nil
	}

	return &proto.Response{
		Status: http.StatusOK,
	}, nil
}

func (s *Service) Token(ctx context.Context, in *proto.TokenRequest) (*proto.TokenResponse, error) {
	const op = "services.auth.Token"

	log := s.log.With("op", op)

	if err := validator.ValidateStruct(in); err != nil {
		msg := "failed to validate request"

		log.Error(msg, sl.Err(err))

		return &proto.TokenResponse{
			Meta: &proto.Response{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			},
		}, nil
	}

	u, err := s.users.FindByEmail(ctx, in.GetEmail())
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			msg := "the user with given email not found"

			log.Error(msg, sl.Err(err))

			return &proto.TokenResponse{
				Meta: &proto.Response{
					Status: http.StatusNotFound,
					Error:  msg,
				},
			}, nil
		}

		msg := "failed to find user by email"

		log.Error(msg, sl.Err(err))

		return &proto.TokenResponse{
			Meta: &proto.Response{
				Status: http.StatusInternalServerError,
				Error:  msg,
			},
		}, nil
	}

	if bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(in.GetPassword())) != nil {
		msg := "invalid password"

		log.Error(msg, sl.Err(err))

		return &proto.TokenResponse{
			Meta: &proto.Response{
				Status: http.StatusBadRequest,
				Error:  msg,
			},
		}, nil
	}

	access, err := jwt.CreateToken(
		&data.TokenPayload{
			ID:     uuid.NewString(),
			UserID: u.ID,
			Email:  u.Email,
		},
		s.access.TTL,
		s.access.PrivateKey,
	)
	if err != nil {
		msg := "failed to create access token"

		log.Error(msg, sl.Err(err))

		return &proto.TokenResponse{
			Meta: &proto.Response{
				Status: http.StatusInternalServerError,
				Error:  msg,
			},
		}, nil
	}

	refresh, err := jwt.CreateToken(
		&data.TokenPayload{
			ID:     uuid.NewString(),
			UserID: u.ID,
			Email:  u.Email,
		},
		s.refresh.TTL,
		s.refresh.PrivateKey,
	)
	if err != nil {
		msg := "failed to create refresh token"

		log.Error(msg, sl.Err(err))

		return &proto.TokenResponse{
			Meta: &proto.Response{
				Status: http.StatusInternalServerError,
				Error:  msg,
			},
		}, nil
	}

	return &proto.TokenResponse{
		Meta: &proto.Response{
			Status: http.StatusOK,
		},
		AccessToken:  access.Token,
		RefreshToken: refresh.Token,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, in *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	const op = "services.auth.Token"

	log := s.log.With("op", op)

	if err := validator.ValidateStruct(in); err != nil {
		msg := "failed to validate request"

		log.Error(msg, sl.Err(err))

		return &proto.RefreshResponse{
			Meta: &proto.Response{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			},
		}, nil
	}

	payload, err := jwt.ValidateToken(in.GetRefresh(), s.refresh.PublicKey)
	if err != nil {
		msg := "refresh token is invalid"

		log.Error(msg, sl.Err(err))

		return &proto.RefreshResponse{
			Meta: &proto.Response{
				Status: http.StatusForbidden,
				Error:  err.Error(),
			},
		}, nil
	}

	access, err := jwt.CreateToken(
		payload,
		s.access.TTL,
		s.access.PrivateKey,
	)
	if err != nil {
		msg := "failed to create access token"

		log.Error(msg, sl.Err(err))

		return &proto.RefreshResponse{
			Meta: &proto.Response{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			},
		}, nil
	}

	return &proto.RefreshResponse{
		Meta: &proto.Response{
			Status: http.StatusOK,
		},
		AccessToken: access.Token,
	}, nil
}

func (s *Service) ResetPassword(ctx context.Context, in *proto.ResetPasswordRequest) (*proto.Response, error) {
	const op = "services.auth.ResetPassword"

	log := s.log.With("op", op)

	if err := validator.ValidateStruct(in); err != nil {
		msg := "failed to validate request"

		log.Error(msg, sl.Err(err))

		return &proto.Response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &proto.Response{Status: http.StatusOK}, nil
}

func (s *Service) ConfirmSingUp(ctx context.Context, in *proto.ConfirmSignUpRequest) (*proto.Response, error) {
	const op = "services.auth.ConfirmSingUp"

	log := s.log.With("op", op)

	if err := validator.ValidateStruct(in); err != nil {
		msg := "failed to validate request"

		log.Error(msg, sl.Err(err))

		return &proto.Response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	return &proto.Response{Status: http.StatusOK}, nil
}

func decodeRSAKeys(private, public string) (pem, pub []byte, err error) {
	pem, err = base64.StdEncoding.DecodeString(private)
	if err != nil {
		return
	}

	pub, err = base64.StdEncoding.DecodeString(public)
	return
}
