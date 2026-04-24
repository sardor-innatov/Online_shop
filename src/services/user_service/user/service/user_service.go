package service

import (
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/services/user_service/user/dto"
	"online_shop/src/services/user_service/user/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	Update(ctx route.Context) error
	Delete(ctx route.Context) error
	GetOne(ctx route.Context) (*dto.UserGetDto, error)
	GetPage(ctx route.Context) ([]dto.UserGetDto, error)
}

type userService struct {
	conn *pgxpool.Pool
	repo repository.UserRepository
}

func NewUserService(conn *pgxpool.Pool) UserService {
	return &userService{
		conn: conn,
		repo: repository.NewUserRepository(conn),
	}
}

func (s *userService) Update(ctx route.Context) error {

	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid user id"})
			return err
		}
	}

	var dto dto.UserUpdateDto
	err = ctx.Bind(&dto)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return err
		}
	}

	user, err := s.repo.SelectOne(id, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "user not found"})
			return err
		}
	}

	{
		user.FirstName = dto.FirstName
		user.LastName = dto.LastName
		user.Email = dto.Email
	}

	err = s.repo.Udpate(id, user, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return err
		}
	}

	return nil
}

func (s *userService) Delete(ctx route.Context) error {

	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid user id"})
			return err
		}
	}

	_, err = s.repo.SelectOne(id, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "user not found"})
			return err
		}
	}

	err = s.repo.Delete(id,ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return err
		}
	}

	return nil
}

func (s *userService) GetOne(ctx route.Context) (*dto.UserGetDto, error) {

	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid user id"})
			return nil, err
		}
	}

	user, err := s.repo.SelectOne(id,ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "user not found"})
			return nil, err
		}
	}

	dto := dto.UserGetDto{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return &dto, nil
}

func (s *userService) GetPage(ctx route.Context) ([]dto.UserGetDto, error) {

	pageStr := ctx.QueryParam("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid user id"})
			return nil, err
		}
		if page <= 0 {
			page = 1
		}
	}

	sizeStr := ctx.QueryParam("size")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid user id"})
			return nil, err
		}
		if size <= 0 {
			size = 1
		}
		if size > 25 {
			size = 25
		}
	}

	users, err := s.repo.SelectWithPage((size*page - size), size,ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "user not found"})
			return nil, err
		}
	}

	var dtos []dto.UserGetDto

	for _, user := range users {

		dtos = append(dtos, dto.UserGetDto{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}

	return dtos, nil
}
