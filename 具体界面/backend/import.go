package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)

// movieFile 对应 movies1.json 的文件结构。
type movieFile struct {
	Actor  string  `json:"actor"`
	Count  int     `json:"count"`
	Movies []Movie `json:"movies"`
}

// importFromJSON 读取 JSON 文件并 upsert 到数据库，返回导入条数。
func importFromJSON(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	var mf movieFile
	if err := json.Unmarshal(data, &mf); err != nil {
		return 0, err
	}
	return upsertMovies(mf.Movies)
}

// runImport 执行手动导入：go run . -import <json路径>
func runImport(args []string) {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	path := fs.String("json", "", "要导入的 JSON 文件路径（必填）")
	dbPath := fs.String("db", "", "数据库文件路径，默认 movies.db")
	_ = fs.Parse(args)

	if *path == "" {
		log.Fatalf("用法: go run . -import -json <movies1.json路径> [-db <数据库路径>]")
	}
	p := *dbPath
	if p == "" {
		exeDir, _ := os.Executable()
		p = filepath.Join(filepath.Dir(exeDir), "movies.db")
	}
	if err := openDB(p); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	n, err := importFromJSON(*path)
	if err != nil {
		log.Fatalf("导入失败: %v", err)
	}
	log.Printf("导入完成，共写入/更新 %d 部影片", n)
}
