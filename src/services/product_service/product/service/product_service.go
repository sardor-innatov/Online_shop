package service

import (
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/services/product_service/product/dto"
	"online_shop/src/services/product_service/product/model"
	"online_shop/src/services/product_service/product/repository"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService interface {
	Create(ctx route.Context) (int64, error)
	Update(ctx route.Context) error
	UpdateStock(ctx route.Context) error
	Delete(ctx route.Context) error
	GetOne(ctx route.Context) (*dto.ProductGetDto, error)
	GetPage(ctx route.Context) ([]dto.ProductGetDto, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(conn *pgxpool.Pool) ProductService {
	return &productService{
		repo: repository.NewProductRepository(conn),
	}
}

func (s *productService) Create(ctx route.Context) (int64, error) {

	var dto dto.ProductCreateDto
	err := ctx.Bind(&dto)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return 0, err
		}
	}

	product := model.Product{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Stock:       dto.Stock,
		CategoryId:  dto.CategoryId,
	}

	id, err := s.repo.Insert(&product, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "failed to create product :" + err.Error()})
			return 0, err
		}
	}

	return id, nil
}

func (s *productService) Update(ctx route.Context) error {

	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid product id"})
			return err
		}
	}

	var dto dto.ProductCreateDto
	err = ctx.Bind(&dto)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return err
		}
	}

	product, err := s.repo.SelectOne(id,ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "product not found"})
			return err
		}
	}

	{
		product.Name = dto.Name
		product.Description = dto.Description
		product.Price = dto.Price
		product.Stock = dto.Stock
		product.CategoryId = dto.CategoryId
		product.UpdatedAt = time.Now()
	}

	err = s.repo.Update(id, product, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return err
		}
	}

	return nil
}

func (s *productService) UpdateStock(ctx route.Context) error {

	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid product id"})
			return err
		}
	}

	quantityStr := ctx.QueryParam("quantity")
	quantity, err := strconv.ParseInt(quantityStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid quantity"})
			return err
		}
	}

	_, err = s.repo.SelectOne(id, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "product not found"})
			return err
		}
	}

	err = s.repo.UpdateStock(id, int(quantity), ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "failed to update stock"})
			return err
		}
	}

	return nil
}

func (s *productService) Delete(ctx route.Context) error {

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
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "product not found"})
			return err
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

func (s *productService) GetOne(ctx route.Context) (*dto.ProductGetDto, error) {

	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid user id"})
			return nil, err
		}
	}

	dto, err := s.repo.SelectOneWithCategoryName(id, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "product not found"})
			return nil, err
		}
	}

	return dto, nil
}

func (s *productService) GetPage(ctx route.Context) ([]dto.ProductGetDto, error) {

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

	var categoryId int64
	categoryIdStr := ctx.QueryParam("categoryId")
	if categoryIdStr != "" {
		categoryId, err = strconv.ParseInt(categoryIdStr, 10, 64)
		{
			if err != nil {
				ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid query param"})
				return nil, err
			}
			if categoryId <= 0 {
				categoryId = 0
			}
		}
	}

	var minPrice float64
	minPriceStr := ctx.QueryParam("minPrice")
	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 10)
		{
			if err != nil {
				ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid query param"})
				return nil, err
			}
			if minPrice <= 0 {
				minPrice = 0
			}
		}
	}

	var maxPrice float64
	maxPriceStr := ctx.QueryParam("maxPrice")
	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 10)
		{
			if err != nil {
				ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid query param"})
				return nil, err
			}
			if maxPrice <= 0 {
				maxPrice = 0
			}
		}
	}

	products, err := s.repo.SelectWithPage((size*page - size), size, categoryId, float32(minPrice), float32(maxPrice), ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "product not found"})
			return nil, err
		}
	}

	return products, nil
}
