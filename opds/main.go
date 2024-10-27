package main

import (
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo/v4"
)

type Server struct {
}

func (s *Server) Mount(e *echo.Echo) {
	base := e.Group("/opds")
	base.GET("", s.GetHello)

	base.GET("/download", func(c echo.Context) error {
		// Set Content-Type for EPUB
		c.Response().Header().Set(echo.HeaderContentType, "application/epub+zip")

		// Set Content-Disposition to attachment to prompt download
		c.Response().Header().Set("Content-Disposition", `attachment; filename="book.epub"`)

		// Serve the file
		return c.File("book.epub")
	})
}

func (s *Server) GetHello(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	return c.Render(http.StatusOK, "index", "World")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
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
