import { Capacitor } from '@capacitor/core'
import seed from './seed/movies.json'

// 影片数据层：完全本地，无后端。
// - Android App（原生平台）：手机本地 SQLite（@capacitor-community/sqlite）
// - 桌面浏览器预览：localStorage（种子数据首次写入）
// 首次启动时把内置的 39 条种子数据导入，笔记等修改保存在本机。

const DB_NAME = 'movies'
const LS_KEY = 'moviedb_v1'

const MONTHS = { Jan: 1, Feb: 2, Mar: 3, Apr: 4, May: 5, Jun: 6, Jul: 7, Aug: 8, Sep: 9, Oct: 10, Nov: 11, Dec: 12 }

// "09 Oct 2026" -> 可比较的时间戳，解析失败返回 0
export function parseDate(s) {
  const m = /^(\d{1,2})\s+(\w{3})\s+(\d{4})$/.exec(s || '')
  if (!m) return 0
  const mo = MONTHS[m[2]]
  if (!mo) return 0
  return Date.UTC(+m[3], mo - 1, +m[1])
}

function sortMovies(rows) {
  return [...rows].sort((a, b) => parseDate(b.date) - parseDate(a.date))
}

// ---------- 原生平台（Android）：本地 SQLite ----------
async function initNative() {
  const { CapacitorSQLite, SQLiteConnection } = await import('@capacitor-community/sqlite')
  const sqlite = new SQLiteConnection(CapacitorSQLite)
  const db = await sqlite.createConnection(DB_NAME, false, 'no-encryption', 1, false)
  await db.open()
  await db.execute(`
    CREATE TABLE IF NOT EXISTS movies (
      code     TEXT PRIMARY KEY,
      title    TEXT NOT NULL,
      link     TEXT NOT NULL,
      duration TEXT NOT NULL,
      date     TEXT NOT NULL,
      cover    TEXT NOT NULL,
      note     TEXT NOT NULL DEFAULT ''
    )
  `)
  const res = await db.query('SELECT COUNT(*) AS n FROM movies')
  if ((res.values?.[0]?.n ?? 0) === 0) {
    for (const m of seed) {
      await db.run(
        'INSERT OR REPLACE INTO movies (code, title, link, duration, date, cover, note) VALUES (?, ?, ?, ?, ?, ?, ?)',
        [m.code, m.title, m.link, m.duration, m.date, m.cover, m.note ?? ''],
      )
    }
  }
  return {
    async list(q) {
      const r = q
        ? await db.query('SELECT * FROM movies WHERE title LIKE ? OR code LIKE ?', [`%${q}%`, `%${q}%`])
        : await db.query('SELECT * FROM movies')
      return sortMovies(r.values ?? [])
    },
    async get(code) {
      const r = await db.query('SELECT * FROM movies WHERE code = ?', [code])
      return r.values?.[0] ?? null
    },
    async saveNote(code, note) {
      await db.run('UPDATE movies SET note = ? WHERE code = ?', [note, code])
    },
  }
}

// ---------- 浏览器（桌面预览）：localStorage ----------
function initWeb() {
  function load() {
    let rows = null
    try {
      rows = JSON.parse(localStorage.getItem(LS_KEY))
    } catch {
      /* 忽略损坏数据，重新播种 */
    }
    if (!Array.isArray(rows) || rows.length === 0) {
      rows = seed.map((m) => ({ ...m, note: m.note ?? '' }))
      persist(rows)
    }
    return rows
  }
  function persist(rows) {
    localStorage.setItem(LS_KEY, JSON.stringify(rows))
  }
  return {
    async list(q) {
      const k = (q || '').toLowerCase()
      const rows = load().filter(
        (m) => m.title.toLowerCase().includes(k) || m.code.toLowerCase().includes(k),
      )
      return sortMovies(rows)
    },
    async get(code) {
      return load().find((m) => m.code === code) ?? null
    },
    async saveNote(code, note) {
      const rows = load()
      const m = rows.find((x) => x.code === code)
      if (m) {
        m.note = note
        persist(rows)
      }
    },
  }
}

let storePromise = null
function getStore() {
  if (!storePromise) {
    storePromise = Capacitor.isNativePlatform()
      ? initNative()
      : Promise.resolve(initWeb())
  }
  return storePromise
}

// 供页面调用的三个接口
export async function listMovies(q = '') {
  return (await getStore()).list(q)
}
export async function getMovie(code) {
  return (await getStore()).get(code)
}
export async function saveNote(code, note) {
  return (await getStore()).saveNote(code, note)
}
