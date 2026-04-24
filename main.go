package main

import (
	_ "online_shop/src/common/config"
	_ "online_shop/docs"
	cmd "online_shop/src"
)

// @title           Online-Shop API
// @version         1.0
// @BasePath  /api/v1
// @host      localhost:8080
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Введите токен в формате: Bearer <your_token>
func main() {

	cmd.Exec()

}
