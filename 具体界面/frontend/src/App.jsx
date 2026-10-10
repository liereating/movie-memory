import { Route, Routes } from 'react-router-dom'
import ListPage from './pages/ListPage'
import DetailPage from './pages/DetailPage'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<ListPage />} />
      <Route path="/movie/:code" element={<DetailPage />} />
    </Routes>
  )
}
