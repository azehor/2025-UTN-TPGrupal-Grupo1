import type { FC } from 'react'
import './App.css'
import Home from "./pages/Home"
import { Link, Route, Routes } from 'react-router-dom'



const App: FC = () => {
  return (
    <div>
      <nav>
        <Link to="/">Inicio</Link>
      </nav>

      <Routes>
        <Route path="/" element={<Home />} />
      </Routes>
    </div>
  )
}


export default App
