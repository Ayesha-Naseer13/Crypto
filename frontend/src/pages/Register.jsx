"use client"

import { useState } from "react"
import { Link, useNavigate } from "react-router-dom"
import { Box, Typography, Alert, Paper } from "@mui/material"
import Button from "../components/ui/Button"
import Input from "../components/ui/Input"
import { authAPI } from "../services/api"

export default function Register() {
  const [formData, setFormData] = useState({ email: "", fullName: "", cnic: "" })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)
  const navigate = useNavigate()

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError("")
    try {
      await authAPI.register(formData)
      setSuccess(true)
      setTimeout(() => navigate("/login"), 2000)
    } catch (err) {
      setError(err.response?.data?.error || "Registration failed")
    }
    setLoading(false)
  }

  return (
    <Box sx={{ minHeight: "100vh", display: "flex", alignItems: "center", justifyContent: "center", p: 2 }}>
      <Paper sx={{ p: 4, maxWidth: 400, width: "100%" }}>
        <Typography variant="h4" gutterBottom align="center" fontWeight={700}>
          Create Account
        </Typography>
        <Typography variant="body2" color="text.secondary" align="center" mb={3}>
          Register for a new crypto wallet
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        {success && (
          <Alert severity="success" sx={{ mb: 2 }}>
            Registration successful! Redirecting to login...
          </Alert>
        )}

        <form onSubmit={handleSubmit}>
          <Input
            label="Full Name"
            name="fullName"
            value={formData.fullName}
            onChange={handleChange}
            required
            sx={{ mb: 2 }}
          />
          <Input
            label="Email Address"
            name="email"
            type="email"
            value={formData.email}
            onChange={handleChange}
            required
            sx={{ mb: 2 }}
          />
          <Input
            label="CNIC / National ID"
            name="cnic"
            value={formData.cnic}
            onChange={handleChange}
            required
            sx={{ mb: 2 }}
          />
          <Button type="submit" variant="contained" fullWidth loading={loading}>
            Register
          </Button>
        </form>

        <Typography variant="body2" align="center" mt={3}>
          Already have an account?{" "}
          <Link to="/login" style={{ color: "#6366f1" }}>
            Sign In
          </Link>
        </Typography>
      </Paper>
    </Box>
  )
}
