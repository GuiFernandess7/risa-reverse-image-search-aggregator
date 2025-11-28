package filetools

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/GuiFernandess7/risa/internal/services/search"
	"github.com/labstack/echo/v4"
)

func (imgH ImageHandler) UploadImage(c echo.Context) error {
	log.Println("[STARTING] - Calling route /image/upload...")

	// Convert to utils -----
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "Error retrieving file from form")
	}

	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error opening file")
	}
	defer src.Close()

	log.Println("[STARTING] - Reading file...")
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, src)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading file")
	}

	// Convert to utils -----
	
	engineName := c.FormValue("engine")
	searchService, asyncService, err := search.GetEngine(engineName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if asyncService != nil {
		log.Println("[STARTING] - Starting async search service...")
		jobID, err := asyncService.Start(search.SearchInput{
			ImageBytes: buf.Bytes(),
		})
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"engine": asyncService.Name(),
			"job_id": jobID,
		})
	}

	log.Println("[STARTING] - Starting search service...")
	result, err := searchService.Search(search.SearchInput{
		ImageBytes: buf.Bytes(),
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"engine": searchService.Name(),
		"result": result,
	})
}

func (imgH ImageHandler) CheckStatusAsync(c echo.Context) error {
	engineName := c.QueryParam("engine")
	jobID := c.QueryParam("job_id")

	if engineName == "" || jobID == "" {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "engine and job_id are required",
		})
	}

	_, asyncService, err := search.GetEngine(engineName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid engine",
		})
	}

	if asyncService == nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "this engine does not support async status",
		})
	}

	result, err := asyncService.Check(jobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"engine": engineName,
		"result": result,
	})
}
