package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"online_shop/src/common/helper"
	"online_shop/src/common/http/route"
	"online_shop/src/services/user_service/role/dto"
	"online_shop/src/services/user_service/role/repository"
	migrate "online_shop/src/services/user_service/user_migrate"

	"testing"
	"github.com/go-openapi/testify/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func dbSetup() *pgxpool.Pool {

	conn, err := pgxpool.New(context.Background(), "postgres://postgres:M067Ss19@localhost:5432/TestingDB?sslmode=disable")
	{
		if err != nil {
			println(err.Error())
		}
	}

	err = migrate.Migrate(conn)
	{
		if err != nil {
			log.Fatal(err)
		}
	}

	return conn
}

func Test_Update(t *testing.T) {

	conn := dbSetup()

	service := NewRoleService(conn)
	repo := repository.NewRoleRepository(conn)

	dto := dto.RoleCreateDto{
		RoleName: "banned_user",
		Permissions: helper.JsonObject{},
	}
	
	body, _ := json.Marshal(dto)

	recorder := httptest.NewRecorder()

	ctx := route.Context{
		Request:  httptest.NewRequest(http.MethodPost, "/api/v1/role", bytes.NewReader(body)),
		Response: recorder,
	}

	id, err := service.CreateRole(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)

	err = repo.Delete(id, ctx)
	{
		if err != nil {
			log.Fatal(err)
		}
	}

}

