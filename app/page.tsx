import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

export default function Page() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="text-center space-y-4">
          <Badge variant="secondary" className="mb-4">
            Downloadable Project
          </Badge>
          <h1 className="text-4xl font-bold text-white">Decentralized Cryptocurrency Wallet</h1>
          <p className="text-slate-400 text-lg">Full-stack blockchain application with Go backend and React frontend</p>
        </div>

        <Card className="bg-slate-800/50 border-slate-700">
          <CardHeader>
            <CardTitle className="text-white">📦 Project Structure</CardTitle>
            <CardDescription>This project contains two separate applications</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid md:grid-cols-2 gap-4">
              <div className="p-4 bg-slate-900/50 rounded-lg">
                <h3 className="font-semibold text-emerald-400 mb-2">backend/</h3>
                <ul className="text-sm text-slate-300 space-y-1">
                  <li>• Go + Gin Framework</li>
                  <li>• MongoDB Atlas</li>
                  <li>• Blockchain with UTXO model</li>
                  <li>• PoW Mining (SHA-256)</li>
                  <li>• ECDSA Digital Signatures</li>
                  <li>• JWT Authentication</li>
                  <li>• Zakat Scheduler (2.5%)</li>
                </ul>
              </div>
              <div className="p-4 bg-slate-900/50 rounded-lg">
                <h3 className="font-semibold text-blue-400 mb-2">frontend/</h3>
                <ul className="text-sm text-slate-300 space-y-1">
                  <li>• React + Vite</li>
                  <li>• Material UI</li>
                  <li>• React Router</li>
                  <li>• Dashboard & Wallet</li>
                  <li>• Mining Interface</li>
                  <li>• Block Explorer</li>
                  <li>• Transaction History</li>
                </ul>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-slate-800/50 border-slate-700">
          <CardHeader>
            <CardTitle className="text-white">🚀 How to Run</CardTitle>
            <CardDescription>Download and run locally</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              <h3 className="font-semibold text-white mb-2">1. Download the project</h3>
              <p className="text-slate-400 text-sm">Click the three dots (⋮) in the top right → "Download ZIP"</p>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-2">2. Setup Backend</h3>
              <pre className="bg-slate-900 p-4 rounded-lg text-sm text-slate-300 overflow-x-auto">
                {`cd backend

# Create .env file with:
MONGODB_URI=mongodb+srv://user:pass@cluster.mongodb.net/db
JWT_SECRET=your-secret-key
AES_KEY=32-byte-secret-key-for-aes-encr

# Install dependencies and run
go mod download
go run main.go`}
              </pre>
            </div>

            <div>
              <h3 className="font-semibold text-white mb-2">3. Setup Frontend</h3>
              <pre className="bg-slate-900 p-4 rounded-lg text-sm text-slate-300 overflow-x-auto">
                {`cd frontend

# Create .env file with:
VITE_API_URL=http://localhost:8080/api

# Install dependencies and run
npm install
npm run dev`}
              </pre>
            </div>

            <div className="p-4 bg-amber-500/10 border border-amber-500/20 rounded-lg">
              <p className="text-amber-400 text-sm">
                <strong>Note:</strong> This project requires Go 1.21+, Node.js 18+, and a MongoDB Atlas account. See
                README.md for detailed instructions.
              </p>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-slate-800/50 border-slate-700">
          <CardHeader>
            <CardTitle className="text-white">✨ Features</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
              {[
                "Blockchain",
                "UTXO Model",
                "PoW Mining",
                "Digital Signatures",
                "Encrypted Keys",
                "Zakat (2.5%)",
                "Block Explorer",
                "REST API",
                "JWT Auth",
                "OTP Login",
                "System Logs",
                "MongoDB Atlas",
              ].map((feature) => (
                <Badge key={feature} variant="outline" className="justify-center py-2 text-slate-300 border-slate-600">
                  {feature}
                </Badge>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
