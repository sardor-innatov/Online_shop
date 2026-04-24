package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"online_shop/src/common/http/route"
	"online_shop/src/services/user_service/user/model"
	"online_shop/src/services/user_service/user/repository"
	dto "online_shop/src/services/user_service/user/service/dto_test"
	migrate "online_shop/src/services/user_service/user_migrate"

	"github.com/go-faker/faker/v4"
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

	service := NewUserService(conn)
	repo := repository.NewUserRepository(conn)

	dto := dto.UserCreateDto{}
	{
		err := faker.FakeData(&dto)
		{
			if err != nil {
				t.Errorf("failed generate fake data %v", err.Error())
			}
		}
	}
	body, _ := json.Marshal(dto)

	recorder := httptest.NewRecorder()

	ctx := route.Context{
		Request:  httptest.NewRequest(http.MethodPost, "/api/v1/user/1", bytes.NewReader(body)),
		Response: recorder,
	}

	id, err := repo.Insert(&model.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  "string123",
	}, ctx)
	{
		if err != nil {
			log.Fatal(err)
		}
	}

	strId := strconv.FormatInt(id, 10)
	ctx.Request.SetPathValue("id", strId)

	err = service.Update(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)

	err = repo.Delete(id, ctx)
	{
		if err != nil {
			log.Fatal(err)
		}
	}

}

func Test_Delete(t *testing.T) {

	conn := dbSetup()

	service := NewUserService(conn)
	repo := repository.NewUserRepository(conn)

	dto := dto.UserCreateDto{}
	{
		err := faker.FakeData(&dto)
		{
			if err != nil {
				t.Errorf("failed generate fake data %v", err.Error())
			}
		}
	}
	body, _ := json.Marshal(dto)

	recorder := httptest.NewRecorder()

	ctx := route.Context{
		Request:  httptest.NewRequest(http.MethodDelete, "/api/v1/user/1", bytes.NewReader(body)),
		Response: recorder,
	}

	id, err := repo.Insert(&model.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  "string123",
	}, ctx)
	{
		if err != nil {
			log.Fatal(err)
		}
	}

	strId := strconv.FormatInt(id, 10)
	ctx.Request.SetPathValue("id", strId)

	err = service.Delete(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)

}

func Test_GetOne(t *testing.T) {

	conn := dbSetup()

	service := NewUserService(conn)
	repo := repository.NewUserRepository(conn)

	dto := dto.UserCreateDto{}
	{
		err := faker.FakeData(&dto)
		{
			if err != nil {
				t.Errorf("failed generate fake data %v", err.Error())
			}
		}
	}

	recorder := httptest.NewRecorder()

	ctx := route.Context{
		Request:  httptest.NewRequest(http.MethodPost, "/api/v1/user/1", bytes.NewReader([]byte{})),
		Response: recorder,
	}

	id, err := repo.Insert(&model.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  "string123",
	}, ctx)
	{
		if err != nil {
			log.Fatal(err)
		}
	}

	strId := strconv.FormatInt(id, 10)
	ctx.Request.SetPathValue("id", strId)

	user, err := service.GetOne(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotNil(t, user)

	err = repo.Delete(id, ctx)
	{
		if err != nil {
			log.Fatal(err)
		}
	}

}

func Test_GetPage(t *testing.T) {

	conn := dbSetup()

	service := NewUserService(conn)
	repo := repository.NewUserRepository(conn)

	dto := dto.UserCreateDto{}
	{
		err := faker.FakeData(&dto)
		{
			if err != nil {
				t.Errorf("failed generate fake data %v", err.Error())
			}
		}
	}

	recorder := httptest.NewRecorder()

	ctx := route.Context{
		Request:  httptest.NewRequest(http.MethodGet, "/api/v1/user/page", bytes.NewReader([]byte{})),
		Response: recorder,
	}

	id, err := repo.Insert(&model.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  "string123",
	}, ctx)
	{
		if err != nil {
			log.Fatal(err)
		}
	}

	url := ctx.Request.URL.Query()
	url.Set("page", strconv.FormatInt(1,10))
	url.Set("size", strconv.FormatInt(5,10))

	ctx.Request.URL.RawQuery = url.Encode() 

	users, err := service.GetPage(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotNil(t, users)

	err = repo.Delete(id, ctx)
	{
		if err != nil {
			log.Fatal(err)
		}
	}

}