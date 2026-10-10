# 电影详情库

网页应用：后端 Go + Gin + SQLite，前端 Vite + React。

## 目录结构

```
具体界面/
├── backend/              # Go 后端
│   ├── main.go           # 服务入口（含 -import 手动导入模式）
│   ├── routes.go         # REST API 路由
│   ├── store.go          # SQLite 存取逻辑
│   ├── import.go         # JSON 导入
│   └── movies.db         # SQLite 数据库（首次导入后生成）
└── frontend/             # React 前端（Vite）
```

## 启动步骤

### 1. 导入数据（首次或数据更新后）

```bash
cd backend
go build -o movies-server.exe .
.\movies-server.exe -import -json "..\..\清洗脚本\movies1.json"
```

> 数据源 `movies1.json` 位于 `清洗脚本` 目录，由清洗脚本从 `movies.json` 按前缀 SNOS/SSIS/SSINI 筛出。

### 2. 启动后端（端口 8080）

```bash
cd backend
.\movies-server.exe
```

### 3. 启动前端（端口 5173）

```bash
cd frontend
npm install
npm run dev
```

浏览器打开 http://localhost:5173

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/movies | 影片列表，`?q=关键词` 按片名/番号搜索 |
| GET | /api/movies/:code | 影片详情（含笔记） |
| PUT | /api/movies/:code/note | 保存笔记，body `{"note":"..."}` |

## 使用

- 列表页：卡片网格展示封面/番号/时长/日期，顶部搜索框可搜索
- 详情页：点击卡片进入，展示基本信息 + 我的笔记（编辑后点"保存笔记"），笔记存 SQLite，重启不丢失
