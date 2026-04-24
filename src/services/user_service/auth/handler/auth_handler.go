package handler

import (
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/services/user_service/auth/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler interface {
}

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(conn *pgxpool.Pool, group route.Group) {

	handler := &authHandler{
		service: service.NewAuthService(conn),
	}

	authGroup := group.Group("/auth")
	{
		authGroup.Handle("POST", "", handler.SignUp)
		authGroup.Handle("POST", "/login", handler.LogIn)
	}
}

// @Summary      sign up
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth   body   dto.SignUpDto  true  "auth"
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       500  {object}  nil
// @Router       /auth [POST]
func (h *authHandler) SignUp(ctx route.Context) error {

	id, err := h.service.SignUp(ctx)
	{
		if err != nil {

			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"id": id})
}

// @Summary      login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth   body   dto.LogInDto  true  "auth"
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       404  {object}  nil
// @Failed       500  {object}  nil
// @Router       /auth/login [POST]
func (h *authHandler) LogIn(ctx route.Context) error {

	token, err := h.service.Login(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"token": token})
}

// // @Summary      get my profile
// // @Tags         auth
// // @Accept       json
// // @Produce      json
// // @Security     ApiKeyAuth
// // @Success      200  {object}  nil
// // @Failed       400  {object}  nil
// // @Failed       404  {object}  nil
// // @Failed       500  {object}  nil
// // @Router       /auth/me [GET]
// func (h *authHandler) GetMe(ctx echo.Context) error {
// 	log.Println("get me")

// 	dto, err := h.service.Me(ctx)
// 	{
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return ctx.JSON(http.StatusOK, echo.Map{"user": dto})
// }

// // @Summary      refresh token
// // @Tags         auth
// // @Accept       json
// // @Produce      json
// // @Security     ApiKeyAuth
// // @Param        auth   body   auth_dto.RefreshTokenDto  true  "auth"
// // @Success      200  {object}  nil
// // @Failed       400  {object}  nil
// // @Failed       404  {object}  nil
// // @Failed       500  {object}  nil
// // @Router       /auth/refresh [POST]
// func (h *authHandler) RefreshToken(ctx echo.Context) error {
// 	log.Println("refresh token")
// 	var dto auth_dto.RefreshTokenDto

// 	err := h.service.RefreshToken(ctx, &dto)
// 	{
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
