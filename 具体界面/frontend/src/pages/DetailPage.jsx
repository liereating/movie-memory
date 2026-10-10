import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getMovie, saveNote as persistNote } from '../db'

export default function DetailPage() {
  const { code } = useParams()
  const [movie, setMovie] = useState(null)
  const [note, setNote] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    setLoading(true)
    getMovie(code)
      .then((m) => {
        setMovie(m)
        setNote(m?.note ?? '')
        setError('')
      })
      .catch(() => setError('加载失败'))
      .finally(() => setLoading(false))
  }, [code])

  async function saveNote() {
    setSaving(true)
    setSaved(false)
    try {
      await persistNote(code, note)
      setSaved(true)
    } catch (e) {
      setError('保存失败，请重试')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="page"><p className="hint">加载中…</p></div>
  if (!movie) return <div className="page"><p className="error">{error || '未找到该影片'}</p><Link to="/">返回列表</Link></div>

  return (
    <div className="page">
      <Link to="/" className="back">← 返回列表</Link>
      <div className="detail">
        <img className="detail-cover" src={movie.cover} alt={movie.code} />
        <div className="detail-info">
          <h2 className="code">{movie.code}</h2>
          <p className="title">{movie.title}</p>
          <dl className="fields">
            <dt>时长</dt><dd>{movie.duration}</dd>
            <dt>发行日期</dt><dd>{movie.date}</dd>
            <dt>番号</dt><dd>{movie.code}</dd>
            <dt>原始链接</dt><dd><a href={movie.link} target="_blank" rel="noreferrer">打开页面 ↗</a></dd>
          </dl>

          <div className="note-box">
            <label htmlFor="note">我的笔记</label>
            <textarea
              id="note"
              rows="8"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              placeholder="写下对这部片子的评价、备忘…"
            />
            <div className="note-actions">
              <button onClick={saveNote} disabled={saving}>
                {saving ? '保存中…' : '保存笔记'}
              </button>
              {saved && <span className="saved">已保存</span>}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
