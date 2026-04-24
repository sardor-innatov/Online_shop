package service

import (
	"log"
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/services/user_service/role/dto"
	"online_shop/src/services/user_service/role/model"
	"online_shop/src/services/user_service/role/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleService interface {
	CreateRole(ctx route.Context) (int64, error)
	// Update(ctx route.Context) error
	// Delete(ctx route.Context) error
	// GetById(ctx route.Context) (*model.Role, error)
	// GetAllRoles(ctx route.Context) ([]model.Role, error)
	// GetAllPermissions(ctx route.Context) ([]model.Permission, error)
	// GiveRole(ctx route.Context, dto *dto.UserRoleCreate) error
}

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(conn *pgxpool.Pool) RoleService {
	return &roleService{
		repo: repository.NewRoleRepository(conn),
	}
}

func (s *roleService) CreateRole(ctx route.Context) (int64, error) {

	var dto *dto.RoleCreateDto
	err := ctx.Bind(&dto)
	{
		if err != nil {
			log.Println("failed to read from json", err)
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
			return 0, err
		}
	}
	log.Println("bindong json")

	role := model.Role{
		RoleName:    dto.RoleName,
		Permissions: dto.Permissions,
	}

	id, err := s.repo.Insert(&role, ctx)
	{
		if err != nil {
			log.Println("failed to create role", err.Error())
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return 0, err
		}
	}

	return id, nil
}

// func (s *roleService) Update(ctx route.Context) error {

// 	var dto dto.RoleCreateDto

// 	idStr := ctx.PathParam("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	{
// 		if err != nil {
// 			log.Println(err.Error())
// 			return ctx.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid quiz ID"})
// 		}
// 	}

// 	err = ctx.Bind(&dto)
// 	{
// 		if err != nil {
// 			log.Println("failed to read from json", err)
// 			return ctx.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
// 		}
// 	}

// 	var role model.Role
// 	err = s.db.Table("roles").
// 		Where("id = ?", id).
// 		First(&role).Error
// 	{
// 		if err != nil {
// 			log.Println("role not found", err.Error())
// 			return ctx.JSON(http.StatusNotFound, map[string]any{"error": err.Error()})
// 		}
// 	}

// 	role = model.Role{
// 		Id:          id,
// 		RoleName:    dto.RoleName,
// 		Permissions: dto.Permissions,
// 	}

// 	result := s.db.Save(&role)
// 	err = result.Error
// 	{
// 		if err != nil {
// 			log.Println("failed to update role", err.Error())
// 			return ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
// 		}
// 	}

// 	return nil
// }

// func (s *roleService) Delete(ctx route.Context) error {

// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	{
// 		if err != nil {
// 			log.Println(err.Error())
// 			return ctx.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid quiz ID"})
// 		}
// 	}

// 	var role model.Role
// 	err = s.db.Table("roles").
// 		Where("id = ?", id).
// 		First(&role).Error
// 	{
// 		if err != nil {
// 			log.Println("role not found", err.Error())
// 			return ctx.JSON(http.StatusNotFound, map[string]any{"error": err.Error()})
// 		}
// 	}

// 	result := s.db.Delete(&role)
// 	err = result.Error
// 	{
// 		if err != nil {
// 			log.Println("failed to delete", err.Error())
// 			return ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
// 		}
// 	}

// 	return nil
// }

// func (s *roleService) GetById(ctx route.Context) (*model.Role, error) {

// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	{
// 		if err != nil {
// 			log.Println(err.Error())
// 			return nil, ctx.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid role id"})
// 		}
// 	}

// 	var role model.Role
// 	err = s.db.Table("roles").
// 		Where("id = ?", id).
// 		First(&role).Error
// 	{
// 		if err != nil {
// 			log.Println("role not found", err.Error())
// 			return nil, ctx.JSON(http.StatusNotFound, map[string]any{"error": err.Error()})
// 		}
// 	}

// 	return &role, nil
// }

// func (s *roleService) GetAllRoles(ctx route.Context) ([]model.Role, error) {

// 	var roles []model.Role

// 	result := s.db.Find(&roles)
// 	err := result.Error
// 	{
// 		if err != nil {
// 			log.Println("failed to get roles", err.Error())
// 			return nil, ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
// 		}
// 	}

// 	return roles, nil
// }

// func (s *roleService) GetAllPermissions(ctx route.Context) ([]model.Permission, error) {

// 	var permissions []model.Permission

// 	result := s.db.Find(&permissions)
// 	err := result.Error
// 	{
// 		if err != nil {
// 			log.Println("failed to get roles", err.Error())
// 			return nil, ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
// 		}
// 	}

// 	return permissions, nil
// }

// func (s *roleService) GiveRole(ctx route.Context, dto *dto.UserRoleCreate) error {
// 	err := ctx.Bind(&dto)
// 	{
// 		if err != nil {
// 			log.Println("failed to read from json", err)
// 			return ctx.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
// 		}
// 	}
// 	log.Println("bindong json")
// 	// checking role existence
// 	var role model.Role
// 	err = s.db.Table("roles").
// 		Where("id = ?", dto.RoleId).
// 		First(&role).Error
// 	{
// 		if err != nil {
// 			log.Println("role not found", err.Error())
// 			return ctx.JSON(http.StatusNotFound, map[string]any{"error": err.Error()})
// 		}
// 	}
// 	log.Println("checked role existance")
// 	// checking user existence
// 	var user user_model.User
// 	err = s.db.Table("users").
// 		Where("id = ?", dto.UserId).
// 		First(&user).Error
// 	{
// 		if err != nil {
// 			log.Println("user not found", err.Error())
// 			return ctx.JSON(http.StatusNotFound, map[string]any{"error": err.Error()})
// 		}
// 	}
// 	log.Println("checked user existence")

// 	result := s.db.Table("user_roles").Create(&dto)
// 	err = result.Error
// 	{
// 		if err != nil {
// 			log.Println("failed to give role", err.Error())
// 			return ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
// 		}
// 	}
// 	log.Println("created role")

// 	return nil
// }
