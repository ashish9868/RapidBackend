package handlers

import (
	"net/http"
	"strconv"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/gin-gonic/gin"
)

func paginatedResponse(ctx *gin.Context, results any, total int) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "50"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 50
	}
	from := 0
	to := total
	if total > 0 {
		from = (page-1)*perPage + 1
	}
	totalPages := 1
	if perPage > 0 {
		totalPages = (total + perPage - 1) / perPage
		if totalPages < 1 {
			totalPages = 1
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"results":     results,
		"total":       total,
		"page":        page,
		"from":        from,
		"to":          to,
		"total_pages": totalPages,
	})
}

func jsonError(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{"global": message})
}

func getAuthUser(ctx *gin.Context, app *core.App) (any, bool) {
	val, exists := ctx.Get(gin.AuthUserKey)
	if !exists || val == nil {
		jsonError(ctx, http.StatusUnauthorized, "Unauthorized")
		return nil, false
	}
	return val, true
}
