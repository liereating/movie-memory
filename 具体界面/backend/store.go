package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"

	_ "modernc.org/sqlite"
)

// Movie 表示一部影片。
type Movie struct {
	Code     string `json:"code"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	Duration string `json:"duration"`
	Date     string `json:"date"`
	Cover    string `json:"cover"`
	Note     string `json:"note"`
}

var db *sql.DB

// openDB 打开（必要时创建）SQLite 数据库并建表。
func openDB(path string) error {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %w", err)
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS movies (
	code       TEXT PRIMARY KEY,
	title      TEXT NOT NULL,
	link       TEXT NOT NULL,
	duration   TEXT NOT NULL,
	date       TEXT NOT NULL,
	cover      TEXT NOT NULL,
	note       TEXT NOT NULL DEFAULT '',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_movies_title ON movies(title);
`)
	if err != nil {
		return fmt.Errorf("初始化表失败: %w", err)
	}
	return nil
}

// upsertMovies 批量插入或更新影片记录。
func upsertMovies(movies []Movie) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
INSERT INTO movies (code, title, link, duration, date, cover)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(code) DO UPDATE SET
	title = excluded.title,
	link = excluded.link,
	duration = excluded.duration,
	date = excluded.date,
	cover = excluded.cover,
	updated_at = CURRENT_TIMESTAMP
`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, m := range movies {
		if _, err := stmt.Exec(m.Code, m.Title, m.Link, m.Duration, m.Date, m.Cover); err != nil {
			return 0, fmt.Errorf("导入 %s 失败: %w", m.Code, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return len(movies), nil
}

// parseDate 将 "09 Oct 2026" 等英文日期解析为可比较的时间，解析失败返回零值。
func parseDate(s string) time.Time {
	if t, err := time.Parse("2 Jan 2006", s); err == nil {
		return t
	}
	return time.Time{}
}

// listMovies 查询影片列表，支持按 title/code 模糊搜索，并按日期倒序（最新在前）。
func listMovies(q string) ([]Movie, error) {
	query := `SELECT code, title, link, duration, date, cover, note FROM movies`
	var args []any
	if q != "" {
		like := "%" + q + "%"
		query += ` WHERE title LIKE ? OR code LIKE ?`
		args = append(args, like, like)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.Code, &m.Title, &m.Link, &m.Duration, &m.Date, &m.Cover, &m.Note); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// date 字段是 "DD Mon YYYY" 文本，SQL 字典序排序不正确，改为按真实日期倒序。
	sort.SliceStable(movies, func(i, j int) bool {
		return parseDate(movies[i].Date).After(parseDate(movies[j].Date))
	})
	return movies, nil
}

// getMovie 按 code 查询单部影片详情。
func getMovie(code string) (*Movie, error) {
	var m Movie
	err := db.QueryRow(
		`SELECT code, title, link, duration, date, cover, note FROM movies WHERE code = ?`,
		code,
	).Scan(&m.Code, &m.Title, &m.Link, &m.Duration, &m.Date, &m.Cover, &m.Note)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// saveNote 保存影片笔记。
func saveNote(code, note string) error {
	res, err := db.Exec(
		`UPDATE movies SET note = ?, updated_at = CURRENT_TIMESTAMP WHERE code = ?`,
		note, code,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("未找到影片 %s", code)
	}
	return nil
}

// countMovies 返回影片总数（用于启动日志）。
func countMovies() int {
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM movies`).Scan(&n); err != nil {
		log.Printf("统计影片数失败: %v", err)
		return 0
	}
	return n
}
