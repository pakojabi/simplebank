package gapi

import (
	"context"
	"database/sql"

	db "github.com/pakojabi/simplebank/db/sqlc"
	"github.com/pakojabi/simplebank/pb"
	"github.com/pakojabi/simplebank/util"
	"github.com/pakojabi/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if violations := validateLoginUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User does not exist: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "Error retrieving user: %s", err)
	}

	if err = util.CheckPassword(req.GetPassword(), user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Authentication failed: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.Make(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create access token: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.Make(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create refresh token: %s", err)
	}

	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session %s", err)
	}
	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	validators := []validator{
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
	}

	for _, validator := range validators {
		if err := validator.validatorFunc(validator.fieldGetter()); err != nil {
			violations = append(violations, fieldViolation(validator.fieldName, err))
		}
	}

	return violations
}
