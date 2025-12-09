import axios from "axios"

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api"

const api = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token")
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Auth API
export const authAPI = {
  register: (data) => api.post("/auth/register", data),
  login: (email) => api.post("/auth/login", { email }),
  verifyOTP: (email, otp) => api.post("/auth/verify-otp", { email, otp }),
  getProfile: () => api.get("/auth/profile"),
  updateProfile: (data) => api.put("/auth/profile", data),
}

// Wallet API
export const walletAPI = {
  getWallet: () => api.get("/wallet"),
  getBalance: () => api.get("/wallet/balance"),
  getUTXOs: () => api.get("/wallet/utxos"),
  addBeneficiary: (data) => api.post("/wallet/beneficiaries", data),
  removeBeneficiary: (id) => api.delete(`/wallet/beneficiaries/${id}`),
  getBeneficiaries: () => api.get("/wallet/beneficiaries"),
}

// Transaction API
export const transactionAPI = {
  send: (data) => api.post("/transactions/send", data),
  getHistory: () => api.get("/transactions/history"),
  getPending: () => api.get("/transactions/pending"),
}

// Mining API
export const miningAPI = {
  mine: () => api.post("/mining/mine"),
  getStatus: () => api.get("/mining/status"),
}

// Block Explorer API
export const blockAPI = {
  getAll: () => api.get("/blocks"),
  getLatest: () => api.get("/blocks/latest"),
  getByHash: (hash) => api.get(`/blocks/${hash}`),
}

// Zakat API
export const zakatAPI = {
  getHistory: () => api.get("/zakat/history"),
  process: () => api.post("/zakat/process"),
}

// Logs API
export const logsAPI = {
  getSystem: () => api.get("/logs/system"),
  getTransactions: () => api.get("/logs/transactions"),
}

export default api
