"use client"

import { useState, useEffect, useContext } from "react"
import {
  Box,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Skeleton,
} from "@mui/material"
import Card from "../components/ui/Card"
import { transactionAPI } from "../services/api"
import { AuthContext } from "../context/AuthContext"

export default function Transactions() {
  const { user } = useContext(AuthContext)
  const [transactions, setTransactions] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        const { data } = await transactionAPI.getHistory()
        setTransactions(data.transactions || [])
      } catch (error) {
        console.error("Failed to fetch transactions:", error)
      }
      setLoading(false)
    }
    fetchTransactions()
  }, [])

  const formatDate = (date) => new Date(date).toLocaleString()
  const truncate = (str, len = 12) => (str ? `${str.substring(0, len)}...` : "-")

  return (
    <Box>
      <Typography variant="h4" gutterBottom fontWeight={700}>
        Transaction History
      </Typography>
      <Typography variant="body1" color="text.secondary" mb={3}>
        View all your confirmed transactions
      </Typography>

      <Card>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Type</TableCell>
                <TableCell>From/To</TableCell>
                <TableCell align="right">Amount</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Date</TableCell>
                <TableCell>Block</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                [1, 2, 3, 4, 5].map((i) => (
                  <TableRow key={i}>
                    {[1, 2, 3, 4, 5, 6].map((j) => (
                      <TableCell key={j}>
                        <Skeleton />
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : transactions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} align="center">
                    <Typography color="text.secondary" py={4}>
                      No transactions found
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                transactions.map((tx, index) => (
                  <TableRow key={index}>
                    <TableCell>
                      <Chip label={tx.direction} color={tx.direction === "sent" ? "error" : "success"} size="small" />
                    </TableCell>
                    <TableCell sx={{ fontFamily: "monospace", fontSize: "0.8rem" }}>
                      {tx.direction === "sent" ? truncate(tx.receiverWalletId) : truncate(tx.senderWalletId)}
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 600 }}>
                      {tx.direction === "sent" ? "-" : "+"}
                      {tx.amount.toFixed(4)}
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={tx.status}
                        color={tx.status === "confirmed" ? "success" : "warning"}
                        size="small"
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell>{formatDate(tx.timestamp)}</TableCell>
                    <TableCell sx={{ fontFamily: "monospace", fontSize: "0.75rem" }}>
                      {truncate(tx.blockHash, 8)}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Card>
    </Box>
  )
}
