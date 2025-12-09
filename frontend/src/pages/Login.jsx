"use client"

import { useState, useContext } from "react"
import { Link, useNavigate } from "react-router-dom"
import { Box, Typography, Alert, Paper } from "@mui/material"
import Button from "../components/ui/Button"
import Input from "../components/ui/Input"
import { authAPI } from "../services/api"
import { AuthContext } from "../context/AuthContext"

export default function Login() {
  const [email, setEmail] = useState("")
  const [otp, setOtp] = useState("")
  const [step, setStep] = useState("email")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const { login } = useContext(AuthContext)
  const navigate = useNavigate()

  const handleSendOTP = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError("")
    try {
      await authAPI.login(email)
      setStep("otp")
    } catch (err) {
      setError(err.response?.data?.error || "Failed to send OTP")
    }
    setLoading(false)
  }

  const handleVerifyOTP = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError("")
    try {
      const { data } = await authAPI.verifyOTP(email, otp)
      login(data.user, data.token)
      navigate("/")
    } catch (err) {
      setError(err.response?.data?.error || "Invalid OTP")
    }
    setLoading(false)
  }

  return (
    <Box sx={{ minHeight: "100vh", display: "flex", alignItems: "center", justifyContent: "center", p: 2 }}>
      <Paper sx={{ p: 4, maxWidth: 400, width: "100%" }}>
        <Typography variant="h4" gutterBottom align="center" fontWeight={700}>
          Welcome Back
        </Typography>
        <Typography variant="body2" color="text.secondary" align="center" mb={3}>
          Sign in to your crypto wallet
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {step === "email" ? (
          <form onSubmit={handleSendOTP}>
            <Input
              label="Email Address"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              sx={{ mb: 2 }}
            />
            <Button type="submit" variant="contained" fullWidth loading={loading}>
              Send OTP
            </Button>
          </form>
        ) : (
          <form onSubmit={handleVerifyOTP}>
            <Typography variant="body2" color="text.secondary" mb={2}>
              OTP sent to {email}
            </Typography>
            <Input label="Enter OTP" value={otp} onChange={(e) => setOtp(e.target.value)} required sx={{ mb: 2 }} />
            <Button type="submit" variant="contained" fullWidth loading={loading}>
              Verify OTP
            </Button>
            <Button variant="text" fullWidth onClick={() => setStep("email")} sx={{ mt: 1 }}>
              Change Email
            </Button>
          </form>
        )}

        <Typography variant="body2" align="center" mt={3}>
          Don't have an account?{" "}
          <Link to="/register" style={{ color: "#6366f1" }}>
            Register
          </Link>
        </Typography>
      </Paper>
    </Box>
  )
}
