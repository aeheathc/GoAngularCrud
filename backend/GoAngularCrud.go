package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	model "github.com/aeheathc/GoAngularCrud/model"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
)

func main() {
	echoServer := echo.New()
	echoServer.GET(   "/users",     getUsers)
	echoServer.GET(   "/users/:id", getUser)
	echoServer.POST(  "/users",     postUser)
	echoServer.PUT(   "/users/:id", putUser)
	echoServer.DELETE("/users/:id", deleteUser)
	//echoServer.RouteNotFound("*", noRoute)
	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	  }))
	port, found := os.LookupEnv("HTTP_PORT")
	if !found {
		port = "8888"
	}
	echoServer.Logger.Fatal(echoServer.Start(":" + port))
}

func noRoute(httpContext echo.Context) error {
	fmt.Println("Route not found")
	return httpContext.String(http.StatusNotFound, "Invalid route")
}

func getUsers(httpContext echo.Context) error {
	fmt.Println("GET /users")
	users, usererr := model.GetUsers()
	if usererr != nil {
		return httpContext.String(http.StatusInternalServerError, "error querying DB: "+usererr.Error())
	}
	return httpContext.JSON(http.StatusOK, users)
}

func getUser(httpContext echo.Context) error {
	fmt.Println("GET /users/:id")
	id, iderr := strconv.Atoi(httpContext.Param("id"))
	if iderr != nil {
		return httpContext.String(http.StatusBadRequest, "id must be an integer")
	}
	user, sqlerr := model.GetUser(id)
	if sqlerr != nil {
		return httpContext.String(http.StatusInternalServerError, "error querying DB: "+sqlerr.Error())
	}
	if user == nil {
		return httpContext.NoContent(http.StatusNotFound)
	}
	return httpContext.JSON(http.StatusOK, *user)
}

func postUser(httpContext echo.Context) error {
	fmt.Println("POST /users")
	user := new(model.UserNoId)
	if err := httpContext.Bind(user); err != nil {
		return err
	}
	fulluser, err := user.Post()
	if err != nil {
		return httpContext.String(http.StatusInternalServerError, "error querying DB: "+err.Error())
	}
	return httpContext.JSON(http.StatusCreated, *fulluser)
}

func putUser(httpContext echo.Context) error {
	fmt.Println("PUT /users/:id")
	id, iderr := strconv.Atoi(httpContext.Param("id"))
	if iderr != nil {
		return httpContext.String(http.StatusBadRequest, "id must be an integer")
	}
	user := new(model.UserNoId)
	if err := httpContext.Bind(user); err != nil {
		return err
	}
	created, sqlerr := user.Put(id)
	if sqlerr != nil {
		return httpContext.String(http.StatusInternalServerError, "error querying DB: "+sqlerr.Error())
	}
	if created {
		return httpContext.NoContent(http.StatusCreated)
	}

	return httpContext.NoContent(http.StatusOK)
}

func deleteUser(httpContext echo.Context) error {
	fmt.Println("DELETE /users/:id")
	id, iderr := strconv.Atoi(httpContext.Param("id"))
	if iderr != nil {
		return httpContext.String(http.StatusBadRequest, "id must be an integer")
	}
	found, err := model.DeleteUser(id)
	if err != nil {
		return httpContext.String(http.StatusInternalServerError, "error querying DB: "+err.Error())
	}
	if !found {
		return httpContext.NoContent(http.StatusNotFound)
	}
	
	return httpContext.NoContent(http.StatusOK)
}