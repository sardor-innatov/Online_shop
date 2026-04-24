package handler

import (
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/common/middleware"
	"online_shop/src/services/product_service/category/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler interface{}

type categoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(conn *pgxpool.Pool, group route.Group) {

	handler := categoryHandler{
		service: service.NewCategoryService(conn),
	}

	categoryGroup := group.Group("/category")
	{
		categoryGroup.Handle("POST", "", handler.Create, middleware.RequireAuth)
		categoryGroup.Handle("PUT", "/{id}", handler.Update, middleware.RequireAuth)
		categoryGroup.Handle("DELETE", "/{id}", handler.Delete, middleware.RequireAuth)
		categoryGroup.Handle("GET", "", handler.GetPage)
	}
}

// @Summary      create
// @Tags         category
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        category   body   dto.CategoryCreateDto  true  "category"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /category [POST]
func (h *categoryHandler) Create(ctx route.Context) error {

	id, err := h.service.Create(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"id": id})
}

// @Summary      update
// @Tags         category
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Param        category   body   dto.CategoryCreateDto  true  "category"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /category/{id} [PUT]
func (h *categoryHandler) Update(ctx route.Context) error {

	err := h.service.Update(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{})
}

// @Summary      delete
// @Tags         category
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Success      201  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /category/{id} [DELETE]
func (h *categoryHandler) Delete(ctx route.Context) error {

	err := h.service.Delete(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{})
}

// @Summary      get pages of category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        page   query   int64  true  "page"
// @Param        size   query   int64  true  "size"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /category [GET]
func (h *categoryHandler) GetPage(ctx route.Context) error {

	dtos, err := h.service.GetPage(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"categories": dtos})
}
