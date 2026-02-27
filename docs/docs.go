package docs

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed swagger.yaml
var swaggerSpec embed.FS

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
  <meta charset="UTF-8">
  <title>RISA API - Documentação</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html { box-sizing: border-box; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "/api/docs/swagger.yaml",
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
      layout: "BaseLayout",
      deepLinking: true,
    });
  </script>
</body>
</html>`

func RegisterDocsRoutes(g *echo.Group) {
	g.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, swaggerUIHTML)
	})

	g.GET("/docs/swagger.yaml", func(c echo.Context) error {
		data, err := swaggerSpec.ReadFile("swagger.yaml")
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.Blob(http.StatusOK, "application/x-yaml", data)
	})
}
