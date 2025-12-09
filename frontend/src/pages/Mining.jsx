"use client"

import { useState, useEffect } from "react"
import { Box, Typography, Grid, Alert, CircularProgress, LinearProgress } from "@mui/material"
import Card from "../components/ui/Card"
import Button from "../components/ui/Button"
import { miningAPI, transactionAPI } from "../services/api"

export default function Mining() {
  const [status, setStatus] = useState(null)
  const [pendingTxs, setPendingTxs] = useState([])
  const [mining, setMining] = useState(false)
  const [result, setResult] = useState(null)
  const [loading, setLoading] = useState(true)

  const fetchStatus = async () => {
    try {
      const [statusRes, pendingRes] = await Promise.all([miningAPI.getStatus(), transactionAPI.getPending()])
      setStatus(statusRes.data)
      setPendingTxs(pendingRes.data.pendingTransactions || [])
    } catch (error) {
      console.error("Failed to fetch mining status:", error)
    }
    setLoading(false)
  }

  useEffect(() => {
    fetchStatus()
    const interval = setInterval(fetchStatus, 10000)
    return () => clearInterval(interval)
  }, [])

  const handleMine = async () => {
    setMining(true)
    setResult(null)
    try {
      const { data } = await miningAPI.mine()
      setResult({ success: true, data })
      fetchStatus()
    } catch (error) {
      setResult({ success: false, error: error.response?.data?.error || "Mining failed" })
    }
    setMining(false)
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom fontWeight={700}>
        Mining
      </Typography>
      <Typography variant="body1" color="text.secondary" mb={3}>
        Mine pending transactions to earn rewards
      </Typography>

      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <Card title="Mining Status">
            {loading ? (
              <CircularProgress />
            ) : (
              <>
                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    Status
                  </Typography>
                  <Typography variant="h6" color={status?.isMining ? "warning.main" : "success.main"}>
                    {status?.isMining ? "Mining..." : "Ready"}
                  </Typography>
                </Box>
                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    Pending Transactions
                  </Typography>
                  <Typography variant="h6">{status?.pendingTransactions || 0}</Typography>
                </Box>
                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    Difficulty
                  </Typography>
                  <Typography variant="h6">{status?.currentDifficulty} zeros</Typography>
                </Box>
                <Box>
                  <Typography variant="body2" color="text.secondary">
                    Latest Block
                  </Typography>
                  <Typography variant="body1">#{status?.latestBlockIndex}</Typography>
                  <Typography variant="caption" sx={{ fontFamily: "monospace" }}>
                    {status?.latestBlockHash?.substring(0, 20)}...
                  </Typography>
                </Box>
              </>
            )}
          </Card>
        </Grid>

        <Grid item xs={12} md={8}>
          <Card title="Start Mining">
            {result && (
              <Alert severity={result.success ? "success" : "error"} sx={{ mb: 2 }}>
                {result.success ? (
                  <>
                    Block mined successfully!
                    <br />
                    Block #{result.data.block?.index} | Hash: {result.data.block?.hash?.substring(0, 20)}...
                  </>
                ) : (
                  result.error
                )}
              </Alert>
            )}

            {mining && (
              <Box mb={2}>
                <Typography variant="body2" mb={1}>
                  Mining in progress... (Proof of Work)
                </Typography>
                <LinearProgress />
              </Box>
            )}

            <Button
              variant="contained"
              size="large"
              onClick={handleMine}
              loading={mining}
              disabled={pendingTxs.length === 0}
              fullWidth
            >
              {pendingTxs.length === 0 ? "No Pending Transactions" : `Mine ${pendingTxs.length} Transaction(s)`}
            </Button>

            <Alert severity="info" sx={{ mt: 2 }}>
              Mining uses Proof-of-Work (PoW) with SHA-256. The hash must start with {status?.currentDifficulty || 5}{" "}
              zeros. Each mined block confirms all pending transactions.
            </Alert>
          </Card>
        </Grid>

        <Grid item xs={12}>
          <Card title="Pending Transactions Pool">
            {pendingTxs.length === 0 ? (
              <Typography color="text.secondary">No pending transactions to mine</Typography>
            ) : (
              pendingTxs.map((tx, index) => (
                <Box key={index} p={2} mb={1} bgcolor="background.default" borderRadius={2}>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} sm={4}>
                      <Typography variant="caption" color="text.secondary">
                        TX ID
                      </Typography>
                      <Typography variant="body2" sx={{ fontFamily: "monospace", fontSize: "0.7rem" }}>
                        {tx.txId?.substring(0, 24)}...
                      </Typography>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                      <Typography variant="caption" color="text.secondary">
                        Amount
                      </Typography>
                      <Typography variant="body1" fontWeight={600}>
                        {tx.amount?.toFixed(4)}
                      </Typography>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                      <Typography variant="caption" color="text.secondary">
                        Type
                      </Typography>
                      <Typography variant="body1">{tx.type}</Typography>
                    </Grid>
                    <Grid item xs={12} sm={2}>
                      <Typography variant="caption" color="text.secondary">
                        Status
                      </Typography>
                      <Typography variant="body1" color="warning.main">
                        Pending
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
