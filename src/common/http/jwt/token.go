package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"online_shop/src/common/config"
	"strings"
	"time"
)

type TokenModel struct {
	Id        int64     `json:"id"`		   	 // user_id
	Role      string    `json:"role"`	   	 // user_role
	ExpiresAt time.Time `json:"expiresAt"`   // token expiration time
	CreatedAt time.Time `json:"createdAt"`   // token creation time
}

type TokenCreateModel struct {
	Id        int64     `json:"id"`		   	 // user_id
	Role      string    `json:"role"`	   	 // user_role
}

func GenerateToken(data TokenCreateModel) (string, error) {

	projectEnv := config.ProjectEnv()

	tokenModel := TokenModel{
		Id: data.Id,
		Role: data.Role,
		ExpiresAt: time.Now().Add((time.Hour*24)*time.Duration(time.Duration(projectEnv.JwtExpire))),
		CreatedAt: time.Now(),
	}

	json, err := json.Marshal(tokenModel)
	{
		if err != nil {
			return "", err
		}
	}

	encodedJson := base64.RawURLEncoding.EncodeToString(json)

	h := hmac.New(sha256.New, []byte(projectEnv.JwtSecret))
	h.Write([]byte(encodedJson))
	signature := h.Sum(nil)

	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	return fmt.Sprintf("%s.%s", encodedJson, encodedSignature), nil
}

func ValidateToken(token string) ([]byte, bool) {

	tokenParts := strings.Split(token, ".")
	{
		if len(tokenParts) != 2 {
			return nil, false
		}
	}

	encodedJson := tokenParts[0]
	encodedSignature := tokenParts[1]

	signature, err := base64.RawURLEncoding.DecodeString(encodedSignature)
	{
		if err != nil{
			panic(err)
		}
	}

	projectEnv := config.ProjectEnv()

	h := hmac.New(sha256.New, []byte(projectEnv.JwtSecret))
	h.Write([]byte(encodedJson))
	newSignature := h.Sum(nil)

	ok := hmac.Equal(signature, newSignature)
	{
		if !ok{
			return nil, false
		}
	}

	json, _ := base64.RawURLEncoding.DecodeString(encodedJson)

	return json, true
}

func GetClaims(jsonData []byte) (*TokenModel, error) {

	var claims TokenModel
	err := json.Unmarshal(jsonData, &claims)
	{
		if err != nil{
			return nil, err
		}
	}
	
	return &claims, nil
}

