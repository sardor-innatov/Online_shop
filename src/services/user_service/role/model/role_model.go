package model

import "online_shop/src/common/helper"

type Role struct {
	Id          int64             
	RoleName    string            
	Permissions helper.JsonObject 
}

type Permission struct {
	Id           int64  
	Path         string 
	Method       string 
	EndpointName string 
}

type UserRole struct {
	Id     int64 
	UserId int64 
	RoleId int64 
}
