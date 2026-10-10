package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// openBrowser 延迟一小段时间后调用系统默认浏览器打开指定地址。
func openBrowser(url string) {
	go func() {
		time.Sleep(600 * time.Millisecond)
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		case "darwin":
			cmd = exec.Command("open", url)
		default:
			cmd = exec.Command("xdg-open", url)
		}
		if err := cmd.Start(); err != nil {
			log.Printf("自动打开浏览器失败: %v", err)
		}
	}()
}

func main() {
	// 手动导入模式：go run . -import -json <movies1.json路径>
	for _, a := range os.Args[1:] {
		if a == "-import" {
			runImport(os.Args[2:])
			return
		}
	}

	// 数据库文件放在 backend 目录下。
	dbPath := os.Getenv("MOVIE_DB_PATH")
	if dbPath == "" {
		exeDir, _ := os.Executable()
		dbPath = filepath.Join(filepath.Dir(exeDir), "movies.db")
	}
	if err := openDB(dbPath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	log.Printf("数据库已就绪，当前共有 %d 部影片", countMovies())

	r := gin.Default()
	registerRoutes(r)

	// 如果前端已构建，静态托管 dist 目录，方便单端口访问。
	dist := filepath.Join("..", "frontend", "dist")
	if _, err := os.Stat(dist); err == nil {
		r.Static("/assets", filepath.Join(dist, "assets"))
		r.StaticFile("/", filepath.Join(dist, "index.html"))
		r.NoRoute(func(c *gin.Context) {
			// 非 /api 请求回退到前端入口（Vue 前端路由）。
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.JSON(404, gin.H{"error": "接口不存在"})
				return
			}
			c.File(filepath.Join(dist, "index.html"))
		})
		log.Println("已托管前端构建产物 frontend/dist")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("服务启动，监听 http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
