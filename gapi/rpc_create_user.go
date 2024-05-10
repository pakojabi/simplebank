package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/pakojabi/simplebank/db/sqlc"
	"github.com/pakojabi/simplebank/pb"
	"github.com/pakojabi/simplebank/util"
	"github.com/pakojabi/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)

			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

type validator struct {
	validatorFunc func (string) error
	fieldGetter func () string
	fieldName string
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	validators := []validator {
		{
			val.ValidateUsername,
			req.GetUsername,
			"username",
		},
		{
			val.ValidatePassword,
			req.GetPassword,
			"password",
		},
		{
			val.ValidateFullUsername,
			req.GetFullName,
			"full_name",
		},
		{
			val.ValidateEmail,
			req.GetEmail,
			"email",
		},
	}

	for _, validator := range validators {
		if err := validator.validatorFunc(validator.fieldGetter()); err != nil {
			violations = append(violations, fieldViolation(validator.fieldName, err))
		}
	}

	return violations
}
