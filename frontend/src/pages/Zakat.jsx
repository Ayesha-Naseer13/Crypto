"use client"

import { useState, useEffect } from "react"
import { Box, Typography, Grid, Alert, Skeleton } from "@mui/material"
import Card from "../components/ui/Card"
import Button from "../components/ui/Button"
import { zakatAPI } from "../services/api"

export default function Zakat() {
  const [history, setHistory] = useState([])
  const [loading, setLoading] = useState(true)
  const [processing, setProcessing] = useState(false)
  const [result, setResult] = useState(null)

  useEffect(() => {
    const fetchHistory = async () => {
      try {
        const { data } = await zakatAPI.getHistory()
        setHistory(data.zakatHistory || [])
      } catch (error) {
        console.error("Failed to fetch zakat history:", error)
      }
      setLoading(false)
    }
    fetchHistory()
  }, [])

  const handleProcess = async () => {
    setProcessing(true)
    setResult(null)
    try {
      await zakatAPI.process()
      setResult({ success: true, message: "Zakat processing completed" })
      const { data } = await zakatAPI.getHistory()
      setHistory(data.zakatHistory || [])
    } catch (error) {
      setResult({ success: false, message: error.response?.data?.error || "Processing failed" })
    }
    setProcessing(false)
  }

  const totalZakat = history.reduce((sum, record) => sum + record.amount, 0)
  const formatDate = (date) => new Date(date).toLocaleDateString()

  return (
    <Box>
      <Typography variant="h4" gutterBottom fontWeight={700}>
        Zakat Deductions
      </Typography>
      <Typography variant="body1" color="text.secondary" mb={3}>
        View and manage your monthly zakat (2.5%) deductions
      </Typography>

      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <Card title="Zakat Summary">
            <Box mb={2}>
              <Typography variant="body2" color="text.secondary">
                Rate
              </Typography>
              <Typography variant="h5" fontWeight={600}>
                2.5%
              </Typography>
            </Box>
            <Box mb={2}>
              <Typography variant="body2" color="text.secondary">
                Frequency
              </Typography>
              <Typography variant="h6">Monthly</Typography>
            </Box>
            <Box mb={2}>
              <Typography variant="body2" color="text.secondary">
                Total Deducted
              </Typography>
              <Typography variant="h5" fontWeight={600} color="secondary">
                {totalZakat.toFixed(4)} COIN
              </Typography>
            </Box>
            <Box>
              <Typography variant="body2" color="text.secondary">
                Deductions Count
              </Typography>
              <Typography variant="h6">{history.length}</Typography>
            </Box>
          </Card>
        </Grid>

        <Grid item xs={12} md={8}>
          <Card title="Manual Processing (Admin)">
            {result && (
              <Alert severity={result.success ? "success" : "error"} sx={{ mb: 2 }}>
                {result.message}
              </Alert>
            )}
            <Alert severity="info" sx={{ mb: 2 }}>
              Zakat is automatically processed on the 1st of each month. Use this button to manually trigger processing
              (admin only).
            </Alert>
            <Button variant="contained" onClick={handleProcess} loading={processing}>
              Process Zakat Now
            </Button>
          </Card>
        </Grid>

        <Grid item xs={12}>
          <Card title="Deduction History">
            {loading ? (
              [1, 2, 3].map((i) => <Skeleton key={i} height={60} sx={{ mb: 1 }} />)
            ) : history.length === 0 ? (
              <Typography color="text.secondary">No zakat deductions yet</Typography>
            ) : (
              history.map((record, index) => (
                <Box key={index} p={2} mb={1} bgcolor="background.default" borderRadius={2}>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item xs={6} sm={3}>
                      <Typography variant="body2" color="text.secondary">
                        Date
                      </Typography>
                      <Typography variant="body1">{formatDate(record.date)}</Typography>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                      <Typography variant="body2" color="text.secondary">
                        Amount
                      </Typography>
                      <Typography variant="body1" fontWeight={600} color="error">
                        -{record.amount?.toFixed(4)} COIN
                      </Typography>
                    </Grid>
                    <Grid item xs={12} sm={6}>
                      <Typography variant="body2" color="text.secondary">
                        Transaction ID
                      </Typography>
                      <Typography variant="body2" sx={{ fontFamily: "monospace", fontSize: "0.75rem" }}>
                        {record.txId?.substring(0, 32)}...
                      </Typography>
                    </Grid>
                  </Grid>
                </Box>
              ))
            )}
          </Card>
        </Grid>
      </Grid>
    </Box>
  )
}
