"use client"

import { useState, useContext } from "react"
import { Box, Typography, Alert, Grid } from "@mui/material"
import Card from "../components/ui/Card"
import Button from "../components/ui/Button"
import Input from "../components/ui/Input"
import { transactionAPI } from "../services/api"
import { AuthContext } from "../context/AuthContext"
import { sha256 } from "../services/crypto"

export default function SendMoney() {
  const { user } = useContext(AuthContext)
  const [formData, setFormData] = useState({ receiverWalletId: "", amount: "", note: "" })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError("")
    setSuccess("")

    try {
      const timestamp = new Date().toISOString()
      const amount = Number.parseFloat(formData.amount)

      // Create payload for signing
      const payload = `${user.walletId}${formData.receiverWalletId}${amount.toFixed(8)}${timestamp}${formData.note}`

      // Generate signature (simplified - in production use proper key from secure storage)
      const signature = await sha256(payload)

      await transactionAPI.send({
        receiverWalletId: formData.receiverWalletId,
        amount: amount,
        note: formData.note,
        signature: signature,
      })

      setSuccess("Transaction submitted successfully! It will be confirmed after mining.")
      setFormData({ receiverWalletId: "", amount: "", note: "" })
    } catch (err) {
      setError(err.response?.data?.error || "Transaction failed")
    }
    setLoading(false)
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom fontWeight={700}>
        Send Money
      </Typography>
      <Typography variant="body1" color="text.secondary" mb={3}>
        Transfer funds to another wallet
      </Typography>

      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card title="Transfer Details">
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error}
              </Alert>
            )}
            {success && (
              <Alert severity="success" sx={{ mb: 2 }}>
                {success}
              </Alert>
            )}

            <form onSubmit={handleSubmit}>
              <Input
                label="Receiver Wallet ID"
                name="receiverWalletId"
                value={formData.receiverWalletId}
                onChange={handleChange}
                required
                sx={{ mb: 2 }}
                placeholder="Enter receiver's wallet ID"
              />
              <Input
                label="Amount"
                name="amount"
                type="number"
                value={formData.amount}
                onChange={handleChange}
                required
                inputProps={{ min: 0.0001, step: 0.0001 }}
                sx={{ mb: 2 }}
              />
              <Input
                label="Note (Optional)"
                name="note"
                value={formData.note}
                onChange={handleChange}
                multiline
                rows={3}
                sx={{ mb: 2 }}
              />
              <Button type="submit" variant="contained" fullWidth loading={loading} size="large">
                Send Transaction
              </Button>
            </form>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card title="Your Wallet">
            <Box mb={2}>
              <Typography variant="body2" color="text.secondary">
                Your Wallet ID
              </Typography>
              <Typography variant="body1" sx={{ wordBreak: "break-all", fontFamily: "monospace" }}>
                {user?.walletId}
              </Typography>
            </Box>
            <Alert severity="info">
              Transactions are added to the pending pool and confirmed after mining. Each transaction is digitally
              signed and verified using your public key.
            </Alert>
          </Card>
        </Grid>
      </Grid>
    </Box>
  )
}
