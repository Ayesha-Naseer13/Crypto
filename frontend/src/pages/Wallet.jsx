"use client"

import { useState, useEffect, useContext } from "react"
import { Box, Typography, Grid, Divider, Chip, Skeleton } from "@mui/material"
import Card from "../components/ui/Card"
import { walletAPI, authAPI } from "../services/api"
import { AuthContext } from "../context/AuthContext"

export default function Wallet() {
  const { user } = useContext(AuthContext)
  const [wallet, setWallet] = useState(null)
  const [utxos, setUtxos] = useState([])
  const [profile, setProfile] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [walletRes, utxoRes, profileRes] = await Promise.all([
          walletAPI.getWallet(),
          walletAPI.getUTXOs(),
          authAPI.getProfile(),
        ])
        setWallet(walletRes.data)
        setUtxos(utxoRes.data.utxos || [])
        setProfile(profileRes.data)
      } catch (error) {
        console.error("Failed to fetch wallet data:", error)
      }
      setLoading(false)
    }
    fetchData()
  }, [])

  return (
    <Box>
      <Typography variant="h4" gutterBottom fontWeight={700}>
        Wallet Profile
      </Typography>
      <Typography variant="body1" color="text.secondary" mb={3}>
        View your wallet details and UTXOs
      </Typography>

      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card title="Profile Information">
            {loading ? (
              [1, 2, 3, 4].map((i) => <Skeleton key={i} height={40} sx={{ mb: 1 }} />)
            ) : (
              <>
                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    Full Name
                  </Typography>
                  <Typography variant="body1">{profile?.fullName}</Typography>
                </Box>
                <Divider sx={{ my: 2 }} />
                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    Email
                  </Typography>
                  <Typography variant="body1">{profile?.email}</Typography>
                </Box>
                <Divider sx={{ my: 2 }} />
                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    CNIC
                  </Typography>
                  <Typography variant="body1">{profile?.cnic}</Typography>
                </Box>
                <Divider sx={{ my: 2 }} />
                <Box>
                  <Typography variant="body2" color="text.secondary">
                    Account Status
                  </Typography>
                  <Chip
                    label={profile?.isVerified ? "Verified" : "Unverified"}
                    color={profile?.isVerified ? "success" : "warning"}
                    size="small"
                  />
                </Box>
              </>
            )}
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card title="Wallet Details">
            {loading ? (
              [1, 2, 3].map((i) => <Skeleton key={i} height={40} sx={{ mb: 1 }} />)
            ) : (
              <>
                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    Wallet ID
                  </Typography>
                  <Typography
                    variant="body1"
                    sx={{ wordBreak: "break-all", fontFamily: "monospace", fontSize: "0.85rem" }}
                  >
                    {wallet?.walletId}
                  </Typography>
                </Box>
                <Divider sx={{ my: 2 }} />
                <Box mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    Public Key
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{ wordBreak: "break-all", fontFamily: "monospace", fontSize: "0.7rem" }}
                  >
                    {wallet?.publicKey}
                  </Typography>
                </Box>
                <Divider sx={{ my: 2 }} />
                <Box>
                  <Typography variant="body2" color="text.secondary">
                    Cached Balance
                  </Typography>
                  <Typography variant="h5" fontWeight={600} color="primary">
                    {wallet?.cachedBalance?.toFixed(4) || "0.0000"} COIN
                  </Typography>
                </Box>
              </>
            )}
          </Card>
        </Grid>

        <Grid item xs={12}>
          <Card title="Unspent Transaction Outputs (UTXOs)">
            {loading ? (
              [1, 2, 3].map((i) => <Skeleton key={i} height={60} sx={{ mb: 1 }} />)
            ) : utxos.length === 0 ? (
              <Typography color="text.secondary">No UTXOs available</Typography>
            ) : (
              utxos.map((utxo, index) => (
                <Box key={index} p={2} mb={1} bgcolor="background.default" borderRadius={2}>
                  <Grid container spacing={2}>
                    <Grid item xs={12} sm={6}>
                      <Typography variant="body2" color="text.secondary">
                        Transaction ID
                      </Typography>
                      <Typography variant="body2" sx={{ fontFamily: "monospace", fontSize: "0.75rem" }}>
                        {utxo.txId?.substring(0, 32)}...
                      </Typography>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                      <Typography variant="body2" color="text.secondary">
                        Output Index
                      </Typography>
                      <Typography variant="body1">{utxo.outputIndex}</Typography>
                    </Grid>
                    <Grid item xs={6} sm={3}>
                      <Typography variant="body2" color="text.secondary">
                        Amount
                      </Typography>
                      <Typography variant="body1" fontWeight={600} color="secondary">
                        {utxo.amount?.toFixed(4)} COIN
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
