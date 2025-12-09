"use client"

import { Routes, Route, Navigate } from "react-router-dom"
import { useState, useEffect } from "react"
import Layout from "./components/Layout"
import Login from "./pages/Login"
import Register from "./pages/Register"
import Dashboard from "./pages/Dashboard"
import SendMoney from "./pages/SendMoney"
import Transactions from "./pages/Transactions"
import Wallet from "./pages/Wallet"
import Mining from "./pages/Mining"
import BlockExplorer from "./pages/BlockExplorer"
import Zakat from "./pages/Zakat"
import Logs from "./pages/Logs"
import { AuthContext } from "./context/AuthContext"

function App() {
  const [user, setUser] = useState(null)
  const [token, setToken] = useState(localStorage.getItem("token"))

  useEffect(() => {
    const storedUser = localStorage.getItem("user")
    if (storedUser) {
      setUser(JSON.parse(storedUser))
    }
  }, [])

  const login = (userData, authToken) => {
    setUser(userData)
    setToken(authToken)
    localStorage.setItem("token", authToken)
    localStorage.setItem("user", JSON.stringify(userData))
  }

  const logout = () => {
    setUser(null)
    setToken(null)
    localStorage.removeItem("token")
    localStorage.removeItem("user")
  }

  return (
    <AuthContext.Provider value={{ user, token, login, logout }}>
      <Routes>
        <Route path="/login" element={!token ? <Login /> : <Navigate to="/" />} />
        <Route path="/register" element={!token ? <Register /> : <Navigate to="/" />} />
        <Route path="/" element={token ? <Layout /> : <Navigate to="/login" />}>
          <Route index element={<Dashboard />} />
          <Route path="send" element={<SendMoney />} />
          <Route path="transactions" element={<Transactions />} />
          <Route path="wallet" element={<Wallet />} />
          <Route path="mining" element={<Mining />} />
          <Route path="blocks" element={<BlockExplorer />} />
          <Route path="zakat" element={<Zakat />} />
          <Route path="logs" element={<Logs />} />
        </Route>
      </Routes>
    </AuthContext.Provider>
  )
}

export default App
