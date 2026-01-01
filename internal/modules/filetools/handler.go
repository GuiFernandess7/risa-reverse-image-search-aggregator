package filetools

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	interfaces "github.com/GuiFernandess7/risa/internal/repository/interfaces"
	auth "github.com/GuiFernandess7/risa/internal/services/auth"
	engine "github.com/GuiFernandess7/risa/internal/services/engine"
	utils "github.com/GuiFernandess7/risa/pkg/utils"
	"github.com/labstack/echo/v4"
)

const SEARCH_COST = 1

func (imgH ImageHandler) UploadImage(c echo.Context) error {
	log.Println("[STARTING] - Verifying available credits...")

	user, err := auth.GetAuthUser(c)
	fmt.Printf("[STARTING] - Authentication: %v", user.ID)
	err = auth.VerifyUserCredits(imgH.DB, user.ID, SEARCH_COST)
	if err != nil {
		switch err {
		case auth.ErrInsufficientCredits:
			return echo.NewHTTPError(http.StatusPaymentRequired, "insufficient credits")
		case auth.ErrCreditBalanceNotFound:
			return echo.NewHTTPError(http.StatusForbidden, "credit account not initialized")
		default:
			return echo.ErrInternalServerError
		}
	}

	log.Println("[STARTING] - Calling route /image/upload...")
	srcFile, err := utils.GetFileObject(c, "file")
	defer srcFile.Close()

	log.Println("[RUNNING] - Reading file...")
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, srcFile)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading file")
	}

	engineName := c.FormValue("engine")
	searchService, asyncService, err := engine.GetEngine(engineName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if asyncService != nil {
		log.Println("[RUNNING] - Starting async search service...")
		jobID, err := asyncService.Start(interfaces.SearchInput{
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

	log.Println("[RUNNING] - Starting search service...")
	result, err := searchService.Search(interfaces.SearchInput{
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
	log.Println("[STARTING] - Calling route /image/check/status...")
	allowedParams := []string{"engine", "job_id"}
	if err := utils.ValidateRequestParams(c, allowedParams); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid parameters",
		})
	}

	engineName := c.QueryParam("engine")
	jobID := c.QueryParam("job_id")
	_, asyncService, err := engine.GetEngine(engineName)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "invalid engine",
		})
	}

	if asyncService == nil {
		log.Printf("[ERROR] - this engine does not support async status")
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

	return c.JSON(http.StatusOK, map[string]any{
		"engine": engineName,
		"result": result,
	})
}
