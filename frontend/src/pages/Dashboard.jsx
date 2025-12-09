"use client"

import { useState, useEffect, useContext } from "react"
import { Grid, Typography, Box, Chip, Skeleton } from "@mui/material"
import {
  AccountBalanceWallet as WalletIcon,
  TrendingUp as TrendingIcon,
  Receipt as ReceiptIcon,
  Pending as PendingIcon,
} from "@mui/icons-material"
import Card from "../components/ui/Card"
import { walletAPI, transactionAPI, miningAPI } from "../services/api"
import { AuthContext } from "../context/AuthContext"

export default function Dashboard() {
  const { user } = useContext(AuthContext)
  const [balance, setBalance] = useState(null)
  const [transactions, setTransactions] = useState([])
  const [miningStatus, setMiningStatus] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [balanceRes, txRes, miningRes] = await Promise.all([
          walletAPI.getBalance(),
          transactionAPI.getHistory(),
          miningAPI.getStatus(),
        ])
        setBalance(balanceRes.data.balance)
        setTransactions(txRes.data.transactions || [])
        setMiningStatus(miningRes.data)
      } catch (error) {
        console.error("Failed to fetch dashboard data:", error)
      }
      setLoading(false)
    }
    fetchData()
  }, [])

  const stats = [
    {
      title: "Balance",
      value: balance !== null ? `${balance.toFixed(4)} COIN` : "-",
      icon: <WalletIcon />,
      color: "#6366f1",
    },
    { title: "Transactions", value: transactions.length, icon: <ReceiptIcon />, color: "#22c55e" },
    { title: "Pending Tx", value: miningStatus?.pendingTransactions || 0, icon: <PendingIcon />, color: "#f59e0b" },
    {
      title: "Latest Block",
      value: `#${miningStatus?.latestBlockIndex || 0}`,
      icon: <TrendingIcon />,
      color: "#ef4444",
    },
  ]

  return (
    <Box>
      <Typography variant="h4" gutterBottom fontWeight={700}>
        Dashboard
      </Typography>
      <Typography variant="body1" color="text.secondary" mb={3}>
        Welcome back, {user?.fullName}
      </Typography>

      <Grid container spacing={3} mb={4}>
        {stats.map((stat, index) => (
          <Grid item xs={12} sm={6} md={3} key={index}>
            <Card sx={{ height: "100%" }}>
              <Box display="flex" alignItems="center" gap={2}>
                <Box sx={{ p: 1.5, borderRadius: 2, bgcolor: `${stat.color}20` }}>
                  <Box sx={{ color: stat.color }}>{stat.icon}</Box>
                </Box>
                <Box>
                  <Typography variant="body2" color="text.secondary">
                    {stat.title}
                  </Typography>
                  {loading ? (
                    <Skeleton width={80} />
                  ) : (
                    <Typography variant="h6" fontWeight={600}>
                      {stat.value}
                    </Typography>
                  )}
                </Box>
              </Box>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card title="Wallet Info">
            <Box mb={2}>
              <Typography variant="body2" color="text.secondary">
                Wallet ID
              </Typography>
              <Typography variant="body1" sx={{ wordBreak: "break-all" }}>
                {user?.walletId || "-"}
              </Typography>
            </Box>
            <Box>
              <Typography variant="body2" color="text.secondary">
                Public Key
              </Typography>
              <Typography variant="body2" sx={{ wordBreak: "break-all", fontSize: "0.75rem" }}>
                {user?.publicKey?.substring(0, 60)}...
              </Typography>
            </Box>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card title="Recent Transactions">
            {loading ? (
              [1, 2, 3].map((i) => <Skeleton key={i} height={40} sx={{ mb: 1 }} />)
            ) : transactions.length === 0 ? (
              <Typography color="text.secondary">No transactions yet</Typography>
            ) : (
              transactions.slice(0, 5).map((tx, index) => (
                <Box
                  key={index}
                  display="flex"
                  justifyContent="space-between"
                  alignItems="center"
                  py={1}
                  borderBottom={1}
                  borderColor="divider"
                >
                  <Box>
                    <Typography variant="body2">{tx.direction === "sent" ? "Sent" : "Received"}</Typography>
                    <Typography variant="caption" color="text.secondary">
                      {new Date(tx.timestamp).toLocaleDateString()}
                    </Typography>
                  </Box>
                  <Chip
                    label={`${tx.direction === "sent" ? "-" : "+"}${tx.amount.toFixed(4)}`}
                    color={tx.direction === "sent" ? "error" : "success"}
                    size="small"
                  />
                </Box>
              ))
            )}
          </Card>
        </Grid>
      </Grid>
    </Box>
  )
}
