package handler

import (
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/common/middleware"
	"online_shop/src/services/user_service/user/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler interface{}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(conn *pgxpool.Pool, group route.Group) {

	handler := userHandler{
		service: service.NewUserService(conn),
	}

	userGroup := group.Group("/user", middleware.RequireAuth,middleware.CheckPermission(conn))
	{
		userGroup.Handle("PUT", "/{id}", handler.Update)
		userGroup.Handle("DELETE", "/{id}", handler.Delete)
		userGroup.Handle("GET", "/{id}", handler.GetById)
		userGroup.Handle("GET", "", handler.GetPage)
	}
}

// @Summary      update
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Param        user   body   dto.UserUpdateDto  true  "user"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /user/{id} [PUT]
func (h *userHandler) Update(ctx route.Context) error {

	err := h.service.Update(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{})
}

// @Summary      delete
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Success      201  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /user/{id} [DELETE]
func (h *userHandler) Delete(ctx route.Context) error {

	err := h.service.Delete(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusNoContent, map[string]any{})
}

// @Summary      get user by id
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /user/{id} [GET]
func (h *userHandler) GetById(ctx route.Context) error {

	dto, err := h.service.GetOne(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"data": dto})
}

// @Summary      get pages of user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        page   query   int64  true  "page"
// @Param        size   query   int64  true  "size"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /user [GET]
func (h *userHandler) GetPage(ctx route.Context) error {

	dtos, err := h.service.GetPage(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"data": dtos})
}
