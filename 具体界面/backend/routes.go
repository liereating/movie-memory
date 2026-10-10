package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// 影片列表，?q= 可按标题或番号搜索
	api.GET("/movies", func(c *gin.Context) {
		movies, err := listMovies(strings.TrimSpace(c.Query("q")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": len(movies), "movies": movies})
	})

	// 影片详情（含笔记）
	api.GET("/movies/:code", func(c *gin.Context) {
		m, err := getMovie(c.Param("code"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if m == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "未找到影片 " + c.Param("code")})
			return
		}
		c.JSON(http.StatusOK, m)
	})

	// 保存笔记
	api.PUT("/movies/:code/note", func(c *gin.Context) {
		var body struct {
			Note string `json:"note"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求体应为 {\"note\": \"...\"}"})
			return
		}
		if err := saveNote(c.Param("code"), body.Note); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// 健康检查
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "count": countMovies()})
	})
}
