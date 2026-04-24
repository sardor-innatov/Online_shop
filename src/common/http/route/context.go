package route

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	keyValue map[string]any
}

func (c *Context) Bind(model any) error {

	err := json.NewDecoder(c.Request.Body).Decode(&model)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) PathParam(name string) string {

	var value string

	if value = c.Request.PathValue(name); value != "" { // if finds is such path param will return it
		return value
	}

	return "" // if no matches return empty string ("")
}

func (c *Context) QueryParam(name string) string {
	var value string

	if value = c.Request.URL.Query().Get(name); value != "" { // if finds matching query param will return it
		return value
	}

	return "" // if no matches return empty string ("")
}

func (c *Context) Set(key string, value any) {

	if c.keyValue == nil {
		c.keyValue = make(map[string]any)
	}
	c.keyValue[key] = value
}

func (c *Context) Get(key string) any {
	return c.keyValue[key]
}

func (c *Context) JSON(status int, content map[string]any) error {

	json, err := json.Marshal(content)
	{
		if err != nil {
			fmt.Println(err)

			return nil
		}
	}

	c.Response.WriteHeader(status)
	c.Response.Write(json)

	return nil
}

func (c *Context) GetIP() string {

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr 
	}
	return ip
}
