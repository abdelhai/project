package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const base = "/home/mustafa/code/hack/project/flash/"

func hello(c echo.Context) error {
	var p map[string]string
	err := c.Bind(&p)
	for k, v := range p {
		ioutil.WriteFile(base+"src/content/"+k, []byte(v), 0644)
	}
	fmt.Println("bind err", err, p)
	build := exec.Command("/usr/bin/node", base+"scripts/build.js")
	build.Dir = "flash"
	out, err := build.Output()
	fmt.Println(out, err)
	if err != nil {
		return nil // todo
	}
	return c.String(200, "hello")
	// js, err := ioutil.ReadFile(base + "public/build/bundle.js")
	// fmt.Println(js, err)
	// css, err := ioutil.ReadFile(base + "public/build/bundle.css")
	// fmt.Println(css, err)
	// return c.JSON(http.StatusOK, map[string]string{
	// 	"js":  string(js),
	// 	"css": string(css),
	// })
}

func main() {
	api := echo.New()
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())

	// Routes
	api.POST("/", hello)

	// Start server
	api.Logger.Fatal(api.Start(":7000"))

}
