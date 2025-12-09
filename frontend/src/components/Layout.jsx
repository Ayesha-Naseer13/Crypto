"use client"

import { useState, useContext } from "react"
import { Outlet, useNavigate, useLocation } from "react-router-dom"
import {
  Box,
  Drawer,
  AppBar,
  Toolbar,
  Typography,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  IconButton,
  Avatar,
  Menu,
  MenuItem,
  Divider,
} from "@mui/material"
import {
  Dashboard as DashboardIcon,
  Send as SendIcon,
  History as HistoryIcon,
  AccountBalanceWallet as WalletIcon,
  Construction as MiningIcon,
  Explore as ExploreIcon,
  VolunteerActivism as ZakatIcon,
  Description as LogsIcon,
  Menu as MenuIcon,
  Logout as LogoutIcon,
} from "@mui/icons-material"
import { AuthContext } from "../context/AuthContext"

const drawerWidth = 240

const menuItems = [
  { text: "Dashboard", icon: <DashboardIcon />, path: "/" },
  { text: "Send Money", icon: <SendIcon />, path: "/send" },
  { text: "Transactions", icon: <HistoryIcon />, path: "/transactions" },
  { text: "Wallet", icon: <WalletIcon />, path: "/wallet" },
  { text: "Mining", icon: <MiningIcon />, path: "/mining" },
  { text: "Block Explorer", icon: <ExploreIcon />, path: "/blocks" },
  { text: "Zakat", icon: <ZakatIcon />, path: "/zakat" },
  { text: "Logs", icon: <LogsIcon />, path: "/logs" },
]

export default function Layout() {
  const [mobileOpen, setMobileOpen] = useState(false)
  const [anchorEl, setAnchorEl] = useState(null)
  const { user, logout } = useContext(AuthContext)
  const navigate = useNavigate()
  const location = useLocation()

  const handleDrawerToggle = () => setMobileOpen(!mobileOpen)
  const handleMenuOpen = (event) => setAnchorEl(event.currentTarget)
  const handleMenuClose = () => setAnchorEl(null)
  const handleLogout = () => {
    logout()
    navigate("/login")
  }

  const drawer = (
    <Box>
      <Toolbar>
        <Typography variant="h6" sx={{ fontWeight: 700, color: "primary.main" }}>
          CryptoWallet
        </Typography>
      </Toolbar>
      <Divider />
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.text} disablePadding>
            <ListItemButton
              selected={location.pathname === item.path}
              onClick={() => {
                navigate(item.path)
                setMobileOpen(false)
              }}
            >
              <ListItemIcon sx={{ color: location.pathname === item.path ? "primary.main" : "inherit" }}>
                {item.icon}
              </ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  )

  return (
    <Box sx={{ display: "flex" }}>
      <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
        <Toolbar>
          <IconButton color="inherit" edge="start" onClick={handleDrawerToggle} sx={{ mr: 2, display: { sm: "none" } }}>
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
            Decentralized Crypto Wallet
          </Typography>
          <IconButton onClick={handleMenuOpen} sx={{ p: 0 }}>
            <Avatar sx={{ bgcolor: "primary.main" }}>{user?.fullName?.charAt(0) || "U"}</Avatar>
          </IconButton>
          <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleMenuClose}>
            <MenuItem disabled>
              <Typography variant="body2">{user?.email}</Typography>
            </MenuItem>
            <Divider />
            <MenuItem onClick={handleLogout}>
              <ListItemIcon>
                <LogoutIcon fontSize="small" />
              </ListItemIcon>
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
      <Box component="nav" sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}>
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{ keepMounted: true }}
          sx={{
            display: { xs: "block", sm: "none" },
            "& .MuiDrawer-paper": { boxSizing: "border-box", width: drawerWidth },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: "none", sm: "block" },
            "& .MuiDrawer-paper": { boxSizing: "border-box", width: drawerWidth },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>
      <Box component="main" sx={{ flexGrow: 1, p: 3, width: { sm: `calc(100% - ${drawerWidth}px)` } }}>
        <Toolbar />
        <Outlet />
      </Box>
    </Box>
  )
}
