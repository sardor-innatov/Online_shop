package service

import (
	"net/http"
	"online_shop/src/common/http/jwt"
	"online_shop/src/common/http/route"
	"online_shop/src/services/user_service/auth/dto"
	user_model "online_shop/src/services/user_service/user/model"
	user_repo "online_shop/src/services/user_service/user/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(ctx route.Context) (int64, error)
	Login(ctx route.Context) (string, error)
}

type authService struct {
	repo user_repo.UserRepository
}

func NewAuthService(conn *pgxpool.Pool) AuthService {
	return &authService{
		repo: user_repo.NewUserRepository(conn),
	}
}

func (s *authService) SignUp(ctx route.Context) (int64, error) {

	var dto dto.SignUpDto
	err := ctx.Bind(&dto)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return 0, err
		}
	}

	if len(dto.Password) > 72 {
		ctx.JSON(http.StatusBadRequest, map[string]any{"error": "password must be less then 72 characters"})
		return 0, bcrypt.ErrPasswordTooLong
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 10)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return 0, err
		}
	}

	dto.Password = string(hashedPassword)

	id, err := s.createUser(&dto, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "email must be unique"})
			return 0, err
		}
	}

	return id, nil
}

func (s *authService) Login(ctx route.Context) (string, error) {

	var dto dto.LogInDto
	err := ctx.Bind(&dto)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return "", err
		}
	}

	user, err := s.repo.SelectByEmail(dto.Email, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, map[string]any{"error": "invalid password or email"})
			return "", err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	{
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, map[string]any{"error": "invalid password or email"})
			return "", err
		}
	}

	// getting role

	role, err := s.repo.GetFirstRole(user.Id, ctx)
	{
		if err != nil {
			return "", err
		}
	}
	tokenModel := jwt.TokenCreateModel{
		Id:   user.Id,
		Role: role,
	}

	token, err := jwt.GenerateToken(tokenModel)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
	}

	return token, nil
}

func (s *authService) createUser(dto *dto.SignUpDto, ctx route.Context) (int64, error) {

	var user user_model.User
	{
		user.FirstName = dto.FirstName
		user.LastName = dto.LastName
		user.Email = dto.Email
		user.Password = dto.Password
	}

	id, err := s.repo.Insert(&user, ctx)
	{
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}
