import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { listMovies } from '../db'

// 格式化日期：09 Oct 2026 -> 2026-10-09
function formatDate(s) {
  const m = s.match(/^(\d{1,2})\s+(\w{3})\s+(\d{4})$/i)
  if (!m) return s
  const months = { Jan: 1, Feb: 2, Mar: 3, Apr: 4, May: 5, Jun: 6, Jul: 7, Aug: 8, Sep: 9, Oct: 10, Nov: 11, Dec: 12 }
  const mm = String(months[m[2]] ?? '').padStart(2, '0')
  const dd = String(m[1]).padStart(2, '0')
  return mm ? `${m[3]}-${mm}-${dd}` : s
}

export default function ListPage() {
  const [movies, setMovies] = useState([])
  const [loading, setLoading] = useState(true)
  const [q, setQ] = useState('')
  const [error, setError] = useState('')

  async function fetchList(keyword) {
    setLoading(true)
    try {
      const list = await listMovies(keyword)
      setMovies(list)
      setError('')
    } catch (e) {
      setError('加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchList('')
  }, [])

  function onSubmit(e) {
    e.preventDefault()
    fetchList(q.trim())
  }

  return (
    <div className="page">
      <header className="topbar">
        <h1>电影详情库</h1>
        <span className="count">{movies.length} 部</span>
      </header>

      <form className="search" onSubmit={onSubmit}>
        <input
          type="text"
          placeholder="按片名或番号搜索…"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />
        <button type="submit">搜索</button>
        {q && (
          <button type="button" className="ghost" onClick={() => { setQ(''); fetchList('') }}>
            清空
          </button>
        )}
      </form>

      {error && <p className="error">{error}</p>}
      {loading && <p className="hint">加载中…</p>}
      {!loading && movies.length === 0 && <p className="hint">没有符合条件的影片</p>}

      <div className="grid">
        {movies.map((m) => (
          <Link key={m.code} to={`/movie/${encodeURIComponent(m.code)}`} className="card">
            <img src={m.cover} alt={m.code} loading="lazy" />
            <div className="card-body">
              <div className="code">{m.code}</div>
              <div className="meta">{formatDate(m.date)} · {m.duration}</div>
              <div className="title">{m.title}</div>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
