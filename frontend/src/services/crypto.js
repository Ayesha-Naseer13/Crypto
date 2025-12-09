import { ec as EC } from "elliptic"

// Choose the curve matching your backend. Common options:
// "p256" (prime256v1) or "secp256k1" (Ethereum/Bitcoin style)
const ec = new EC("p256")

// Async SHA-256 hash function
export const sha256 = async (message) => {
  const msgBuffer = new TextEncoder().encode(message)
  const hashBuffer = await crypto.subtle.digest("SHA-256", msgBuffer)
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  return hashArray.map((b) => b.toString(16).padStart(2, "0")).join("")
}

// Sign transaction (async)
export const signTransaction = async (privateKeyHex, payload) => {
  try {
    const key = ec.keyFromPrivate(privateKeyHex, "hex")
    const msgHash = await sha256(payload)  // ✅ Await hash properly
    const signature = key.sign(msgHash)

    // Pad r and s to 32 bytes each
    const r = signature.r.toString("hex").padStart(64, "0")
    const s = signature.s.toString("hex").padStart(64, "0")

    // Return concatenated r+s hex
    return r + s

    // OR, if your backend expects DER format:
    // return signature.toDER("hex")
  } catch (error) {
    console.error("Signing error:", error)
    throw error
  }
}

// Create transaction payload for signing
export const createSignPayload = (senderID, receiverID, amount, timestamp, note = "") => {
  return `${senderID}${receiverID}${amount.toFixed(8)}${timestamp}${note}`
}
