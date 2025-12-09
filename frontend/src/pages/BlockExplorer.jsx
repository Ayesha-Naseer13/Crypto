"use client"

import { useState, useEffect } from "react"
import { Box, Typography, Accordion, AccordionSummary, AccordionDetails, Grid, Chip, Skeleton } from "@mui/material"
import { ExpandMore as ExpandMoreIcon } from "@mui/icons-material"
import Card from "../components/ui/Card"
import { blockAPI } from "../services/api"

export default function BlockExplorer() {
  const [blocks, setBlocks] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchBlocks = async () => {
      try {
        const { data } = await blockAPI.getAll()
        setBlocks(data.blocks || [])
      } catch (error) {
        console.error("Failed to fetch blocks:", error)
      }
      setLoading(false)
    }
    fetchBlocks()
  }, [])

  const formatDate = (date) => new Date(date).toLocaleString()

  return (
    <Box>
      <Typography variant="h4" gutterBottom fontWeight={700}>
        Block Explorer
      </Typography>
      <Typography variant="body1" color="text.secondary" mb={3}>
        Explore the blockchain - view all blocks and transactions
      </Typography>

      <Card sx={{ mb: 3 }}>
        <Grid container spacing={2}>
          <Grid item xs={6} sm={3}>
            <Typography variant="body2" color="text.secondary">
              Total Blocks
            </Typography>
            <Typography variant="h5" fontWeight={600}>
              {blocks.length}
            </Typography>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Typography variant="body2" color="text.secondary">
              Latest Block
            </Typography>
            <Typography variant="h5" fontWeight={600}>
              #{blocks[0]?.index || 0}
            </Typography>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Typography variant="body2" color="text.secondary">
              Difficulty
            </Typography>
            <Typography variant="h5" fontWeight={600}>
              {blocks[0]?.difficulty || 5}
            </Typography>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Typography variant="body2" color="text.secondary">
              Status
            </Typography>
            <Chip label="Synced" color="success" size="small" />
          </Grid>
        </Grid>
      </Card>

      {loading ? (
        [1, 2, 3].map((i) => <Skeleton key={i} height={100} sx={{ mb: 2 }} />)
      ) : blocks.length === 0 ? (
        <Card>
          <Typography color="text.secondary" align="center">
            No blocks found
          </Typography>
        </Card>
      ) : (
        blocks.map((block, index) => (
          <Accordion key={index} sx={{ mb: 1 }}>
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
              <Grid container alignItems="center" spacing={2}>
                <Grid item>
                  <Chip label={`Block #${block.index}`} color="primary" />
                </Grid>
                <Grid item xs>
                  <Typography variant="body2" sx={{ fontFamily: "monospace", fontSize: "0.75rem" }}>
                    {block.hash?.substring(0, 40)}...
                  </Typography>
                </Grid>
                <Grid item>
                  <Typography variant="body2" color="text.secondary">
                    {block.transactions?.length || 0} tx
                  </Typography>
                </Grid>
              </Grid>
            </AccordionSummary>
            <AccordionDetails>
              <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                  <Typography variant="body2" color="text.secondary">
                    Hash
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{ fontFamily: "monospace", fontSize: "0.7rem", wordBreak: "break-all" }}
                  >
                    {block.hash}
                  </Typography>
                </Grid>
                <Grid item xs={12} md={6}>
                  <Typography variant="body2" color="text.secondary">
                    Previous Hash
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{ fontFamily: "monospace", fontSize: "0.7rem", wordBreak: "break-all" }}
                  >
                    {block.previousHash}
                  </Typography>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Typography variant="body2" color="text.secondary">
                    Nonce
                  </Typography>
                  <Typography variant="body1">{block.nonce}</Typography>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Typography variant="body2" color="text.secondary">
                    Difficulty
                  </Typography>
                  <Typography variant="body1">{block.difficulty}</Typography>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Typography variant="body2" color="text.secondary">
                    Miner
                  </Typography>
                  <Typography variant="body2" sx={{ fontFamily: "monospace", fontSize: "0.7rem" }}>
                    {block.miner?.substring(0, 12)}...
                  </Typography>
                </Grid>
                <Grid item xs={6} sm={3}>
                  <Typography variant="body2" color="text.secondary">
                    Timestamp
                  </Typography>
                  <Typography variant="body2">{formatDate(block.timestamp)}</Typography>
                </Grid>
                <Grid item xs={12}>
                  <Typography variant="body2" color="text.secondary">
                    Merkle Root
                  </Typography>
                  <Typography variant="body2" sx={{ fontFamily: "monospace", fontSize: "0.7rem" }}>
                    {block.merkleRoot}
                  </Typography>
                </Grid>
                {block.transactions?.length > 0 && (
                  <Grid item xs={12}>
                    <Typography variant="body2" color="text.secondary" mb={1}>
                      Transactions ({block.transactions.length})
                    </Typography>
                    {block.transactions.map((tx, txIndex) => (
                      <Box key={txIndex} p={1} mb={1} bgcolor="background.default" borderRadius={1}>
                        <Typography variant="caption" sx={{ fontFamily: "monospace" }}>
                          {tx.txId?.substring(0, 32)}... | {tx.amount?.toFixed(4)} COIN | {tx.type}
                        </Typography>
                      </Box>
                    ))}
                  </Grid>
                )}
              </Grid>
            </AccordionDetails>
          </Accordion>
        ))
      )}
    </Box>
  )
}
