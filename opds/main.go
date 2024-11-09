package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"opds/models"
	"os"
	"text/template"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Server struct {
}

func (s *Server) Mount(e *echo.Echo) {
	base := e.Group("/opds")
	base.GET("/catalogs/:id", s.GetCatalog)
}

func (s *Server) GetCatalog(c echo.Context) error {
	ctx := context.Background()
	apiKey := os.Getenv("GOOGLE_API_KEY")
	driveService, err := drive.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Print(err)
		return c.JSON(400, err)
	}
	catalogID := c.Param("id")
	query := fmt.Sprintf("'%s' in parents", catalogID)

	r, err := driveService.Files.List().Q(query).SupportsAllDrives(true).IncludeItemsFromAllDrives(true).Do()
	if err != nil {
		fmt.Print(err)
		return c.JSON(400, err)
	}

	feeds := models.ConvertFilesToFeeds(r.Files)

	return c.Render(http.StatusOK, "index", feeds)
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	server := Server{}
	server.Mount(e)

	isLambda := os.Getenv("LAMBDA")

	if isLambda == "TRUE" {
		lambdaAdapter := &LambdaAdapter{Echo: e}
		lambda.Start(lambdaAdapter.Handler)
	} else {
		e.Logger.Fatal(e.Start(":1323"))
	}
}
