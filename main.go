package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TODO
const base = "/home/mustafa/code/hack/project/flash/"
const bucket = "haxxorapp"

var sess = session.Must(session.NewSession(&aws.Config{
	Region: aws.String("eu-central-1"),
}))

// Create an uploader with the session and default options
var uploader = s3manager.NewUploader(sess)

func uploadf(name string) error {
	f, err := os.Open(name)
	if err != nil {
		fmt.Printf("failed to open file %q, %v", name, err)
		return fmt.Errorf("failed to open file %q, %v", name, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(strings.Replace(name, base, "", 1)),
		Body:   f,
	})
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %v\n", result.Location)
	return nil
}

func getfiles(root string) *[]string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil
	}
	return &files
	// for _, file := range files {
	// 	fmt.Println(file)
	// }
}

func deploy(c echo.Context) error {
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

	files := getfiles(base + "build")

	for _, file := range *files {
		uploadf(file)
	}
	fmt.Println(files)
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
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// Routes
	api.POST("/", deploy)

	// Start server
	api.Logger.Fatal(api.Start(":7000"))

}
