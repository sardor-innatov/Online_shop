package service

import (
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/services/product_service/category/dto"
	"online_shop/src/services/product_service/category/model"
	"online_shop/src/services/product_service/category/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryService interface {
	Create(ctx route.Context) (int64, error)
	Update(ctx route.Context) error
	Delete(ctx route.Context) error
	GetPage(ctx route.Context) ([]dto.CategoryGetDto, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(conn *pgxpool.Pool) CategoryService {
	return &categoryService{
		repo: repository.NewCategoryRepository(conn),
	}
}

func (s *categoryService) Create(ctx route.Context) (int64, error) {

	var dto dto.CategoryCreateDto

	err := ctx.Bind(dto)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return 0, err
		}
	}

	category := model.Category{
		Name: dto.Name,
	}

	id, err := s.repo.Insert(&category, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return 0, err
		}
	}

	return id, nil
}

func (s *categoryService) Update(ctx route.Context) error {

	var dto dto.CategoryCreateDto
	err := ctx.Bind(dto)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return err
		}
	}

	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid product id"})
			return err
		}
	}

	category, err := s.repo.SelectOne(id, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "category not found"})
		}
	}

	{
		category.Name = dto.Name
	}

	err = s.repo.Update(id, category, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return err
		}
	}

	return nil
}

func (s *categoryService) Delete(ctx route.Context) error {

	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid product id"})
			return err
		}
	}

	_, err = s.repo.SelectOne(id, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "category not found"})
		}
	}

	err = s.repo.Delete(id, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return err
		}
	}

	return nil
}

func (s *categoryService) GetPage(ctx route.Context) ([]dto.CategoryGetDto, error) {

	pageStr := ctx.QueryParam("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid query param"})
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
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid query param"})
			return nil, err
		}

		if size <= 0 {
			size = 1
		}
		if size > 25 {
			size = 25
		}
	}

	categories, err := s.repo.SelectPage(size, page*size-size, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return nil, err
		}
	}

	var dtos []dto.CategoryGetDto
	for _, category := range categories {
		dto := dto.CategoryGetDto{
			Id:   category.Id,
			Name: category.Name,
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}
