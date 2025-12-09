"use client"

import { useState, useEffect } from "react"
import {
  Box,
  Typography,
  Tabs,
  Tab,
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
import { logsAPI } from "../services/api"

export default function Logs() {
  const [tab, setTab] = useState(0)
  const [systemLogs, setSystemLogs] = useState([])
  const [txLogs, setTxLogs] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const [systemRes, txRes] = await Promise.all([logsAPI.getSystem(), logsAPI.getTransactions()])
        setSystemLogs(systemRes.data.logs || [])
        setTxLogs(txRes.data.logs || [])
      } catch (error) {
        console.error("Failed to fetch logs:", error)
      }
      setLoading(false)
    }
    fetchLogs()
  }, [])

  const formatDate = (date) => new Date(date).toLocaleString()

  return (
    <Box>
      <Typography variant="h4" gutterBottom fontWeight={700}>
        System Logs
      </Typography>
      <Typography variant="body1" color="text.secondary" mb={3}>
        View system and transaction activity logs
      </Typography>

      <Card>
        <Tabs value={tab} onChange={(e, v) => setTab(v)} sx={{ mb: 2 }}>
          <Tab label="System Logs" />
          <Tab label="Transaction Logs" />
        </Tabs>

        {tab === 0 && (
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Action</TableCell>
                  <TableCell>Details</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>IP Address</TableCell>
                  <TableCell>Timestamp</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {loading ? (
                  [1, 2, 3, 4, 5].map((i) => (
                    <TableRow key={i}>
                      {[1, 2, 3, 4, 5].map((j) => (
                        <TableCell key={j}>
                          <Skeleton />
                        </TableCell>
                      ))}
                    </TableRow>
                  ))
                ) : systemLogs.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} align="center">
                      <Typography color="text.secondary" py={4}>
                        No system logs
                      </Typography>
                    </TableCell>
                  </TableRow>
                ) : (
                  systemLogs.map((log, index) => (
                    <TableRow key={index}>
                      <TableCell>{log.action}</TableCell>
                      <TableCell sx={{ maxWidth: 200, overflow: "hidden", textOverflow: "ellipsis" }}>
                        {log.details}
                      </TableCell>
                      <TableCell>
                        <Chip label={log.status} color={log.status === "success" ? "success" : "error"} size="small" />
                      </TableCell>
                      <TableCell>{log.ipAddress || "-"}</TableCell>
                      <TableCell>{formatDate(log.timestamp)}</TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        )}

        {tab === 1 && (
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Action</TableCell>
                  <TableCell>TX ID</TableCell>
                  <TableCell align="right">Amount</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Timestamp</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {loading ? (
                  [1, 2, 3, 4, 5].map((i) => (
                    <TableRow key={i}>
                      {[1, 2, 3, 4, 5].map((j) => (
                        <TableCell key={j}>
                          <Skeleton />
                        </TableCell>
                      ))}
                    </TableRow>
                  ))
                ) : txLogs.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} align="center">
                      <Typography color="text.secondary" py={4}>
                        No transaction logs
                      </Typography>
                    </TableCell>
                  </TableRow>
                ) : (
                  txLogs.map((log, index) => (
                    <TableRow key={index}>
                      <TableCell>
                        <Chip
                          label={log.action}
                          color={log.action === "sent" ? "error" : "success"}
                          size="small"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell sx={{ fontFamily: "monospace", fontSize: "0.75rem" }}>
                        {log.txId?.substring(0, 16)}...
                      </TableCell>
                      <TableCell align="right" fontWeight={600}>
                        {log.amount?.toFixed(4)}
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={log.status}
                          color={log.status === "confirmed" ? "success" : "warning"}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>{formatDate(log.timestamp)}</TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Card>
    </Box>
  )
}
