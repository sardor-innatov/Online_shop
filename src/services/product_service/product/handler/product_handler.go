package handler

import (
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/common/middleware"
	"online_shop/src/services/product_service/product/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductHanlder interface{}

type productHandler struct {
	service service.ProductService
}

func NewProductHandler(conn *pgxpool.Pool, group route.Group) {

	handler := productHandler{
		service: service.NewProductService(conn),
	}

	productGroup := group.Group("/product")
	{
		productGroup.Handle("POST", "", handler.Create, middleware.RequireAuth,middleware.CheckPermission(conn))
		productGroup.Handle("PUT", "/{id}", handler.Update, middleware.RequireAuth,middleware.CheckPermission(conn))
		productGroup.Handle("PUT", "/stock/{id}", handler.UpdateStock, middleware.RequireAuth,middleware.CheckPermission(conn))
		productGroup.Handle("DELETE", "/{id}", handler.Delete, middleware.RequireAuth,middleware.CheckPermission(conn))
		productGroup.Handle("GET", "/{id}", handler.GetById)
		productGroup.Handle("GET", "", handler.GetPage)
	}
}

// @Summary      create
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        product   body   dto.ProductCreateDto  true  "product"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /product [POST]
func (h *productHandler) Create(ctx route.Context) error {

	id, err := h.service.Create(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"id": id})
}

// @Summary      update
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Param        product   body   dto.ProductCreateDto  true  "product"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /product/{id} [PUT]
func (h *productHandler) Update(ctx route.Context) error {

	err := h.service.Update(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{})
}

// @Summary      update quantity of items
// @Description  write possitive intagers to add items, and negative intager to decrease quantity of items
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Param        quantity   query   int true  "quantity"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /product/stock/{id} [PUT]
func (h *productHandler) UpdateStock(ctx route.Context) error {

	err := h.service.UpdateStock(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{})
}

// @Summary      delete
// @Tags         product
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Success      201  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /product/{id} [DELETE]
func (h *productHandler) Delete(ctx route.Context) error {

	err := h.service.Delete(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{})
}

// @Summary      get product by id
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id   path   int64  true  "id"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /product/{id} [GET]
func (h *productHandler) GetById(ctx route.Context) error {

	dto, err := h.service.GetOne(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"product": dto})
}

// @Summary      get pages of product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        page   query   int64  true  "page"
// @Param        size   query   int64  true  "size"
// @Param        categoryId   query   int64  false  "categoryId"
// @Param        minPrice   query   float32  false  "minPrice"
// @Param        maxPrice   query   float32  false  "maxPrice"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /product [GET]
func (h *productHandler) GetPage(ctx route.Context) error {

	dtos, err := h.service.GetPage(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"products": dtos})
}
