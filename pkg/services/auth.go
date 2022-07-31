package services

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sabilimaulana/sebelbucks-auth-service/pkg/db"
	"github.com/sabilimaulana/sebelbucks-auth-service/pkg/models"
	"github.com/sabilimaulana/sebelbucks-auth-service/pkg/pb"
	"github.com/sabilimaulana/sebelbucks-auth-service/pkg/utils"
)

type Server struct {
	H   db.Handler
	Jwt utils.JwtWrapper
	pb.UnimplementedAuthServiceServer
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}, nil
	}

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)
	user.UUID = uuid.NewString()

	s.H.DB.Create(&user)

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) RegisterAdmin(ctx context.Context, req *pb.RegisterAdminRequest) (*pb.RegisterAdminResponse, error) {
	var user models.User
	var admin models.Admin

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		return &pb.RegisterAdminResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}, nil
	}

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)
	user.UUID = uuid.NewString()

	admin.UserUUID = user.UUID

	s.H.DB.Create(&user)
	s.H.DB.Create(&admin)

	return &pb.RegisterAdminResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) LoginAdmin(ctx context.Context, req *pb.LoginAdminRequest) (*pb.LoginAdminResponse, error) {
	var user models.User
	var admin models.Admin

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
		return &pb.LoginAdminResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		return &pb.LoginAdminResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	// Check is the user admin
	if result := s.H.DB.Where(&models.Admin{UserUUID: user.UUID}).First(&admin); result.Error != nil {
		return &pb.LoginAdminResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &pb.LoginAdminResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	var user models.User
	var admin models.Admin
	isAdmin := true

	if result := s.H.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	if result := s.H.DB.Where(&models.Admin{UserUUID: user.UUID}).First(&admin); result.Error != nil {
		isAdmin = false
	}

	return &pb.ValidateResponse{
		Status:   http.StatusOK,
		UserUuid: user.UUID,
		IsAdmin:  isAdmin,
	}, nil
}
