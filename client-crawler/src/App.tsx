import "./App.css"
import { Navigate, Route, Routes } from "react-router-dom"

function App() {
  return (
    <Routes>
      <Route path="/" element={<></>} />
      <Route path="/url/:id" element={<></>} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
