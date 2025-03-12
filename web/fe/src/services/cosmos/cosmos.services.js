import { axios } from '../../boot/axios'

class CosmosServices {
  constructor() {
    this.baseURL = `${import.meta.env.VITE_API_URL}/${import.meta.env.VITE_API_NAME}/${import.meta.env.VITE_API_VERSION}`
  }

  async faucet(params = {}) {
    const response = await axios.post(`${this.baseURL}/cosmos/faucet`, {
      params,
    })
    return response.data
  }

  // Request nonce from server
  async requestNonce(walletAddress) {
    try {
      const response = await axios.post(
        `${this.baseURL}/auth/request-nonce?wallet=${walletAddress}`,
      )
      return response.data
    } catch (error) {
      console.error('Error requesting nonce:', error)
      throw error
    }
  }

  // Helper function to convert array-like objects to base64
  arrayToBase64(arr) {
    if (!arr) return ''

    // Convert to Uint8Array if needed
    const uint8Arr = arr instanceof Uint8Array ? arr : new Uint8Array(arr)

    // Convert to base64
    let binary = ''
    uint8Arr.forEach((byte) => {
      binary += String.fromCharCode(byte)
    })
    return window.btoa(binary)
  }

  // Convert base64 to hex
  base64ToHex(base64Str) {
    const raw = window.atob(base64Str)
    let result = ''
    for (let i = 0; i < raw.length; i++) {
      const hex = raw.charCodeAt(i).toString(16)
      result += hex.length === 2 ? hex : '0' + hex
    }
    return result
  }

  // Sign message using Cosmos wallet with sign_amino instead of signDirect
  async signWithWallet(message, wallet) {
    try {
      if (!wallet) {
        throw new Error('No wallet provided')
      }

      const [account] = await wallet.getAccounts()
      const address = account.address

      // Prepare the signDoc for sign_amino
      const signDoc = {
        chain_id: wallet.chainId || '',
        account_number: '0',
        sequence: '0',
        fee: {
          amount: [],
          gas: '0',
        },
        msgs: [
          {
            type: 'sign/MsgSignData',
            value: {
              signer: address,
              data: Buffer.from(message).toString('base64'),
            },
          },
        ],
        memo: '',
      }

      // Call sign_amino method
      const signResponse = await wallet.signDirect(address, signDoc)

      console.log('Sign amino response:', JSON.stringify(signResponse, null, 2))

      // Extract signature and public key
      let signature, pubKey

      if (signResponse.signature && signResponse.signature.signature) {
        signature = signResponse.signature.signature // Base64 string
      } else {
        throw new Error('Signature not found in response')
      }

      if (
        signResponse.signature &&
        signResponse.signature.pub_key &&
        signResponse.signature.pub_key.value
      ) {
        pubKey = signResponse.signature.pub_key.value // Base64 string
      } else if (account.pubkey) {
        pubKey = this.arrayToBase64(account.pubkey) // Fallback from account
      } else {
        throw new Error('Public key not found in response')
      }

      // Stringify the signed document for verification
      const signedDoc = JSON.stringify(signResponse.signed)

      return {
        signature,
        pubKey,
        signDoc: signedDoc,
      }
    } catch (error) {
      console.error('Error signing with wallet (amino):', error)
      throw error
    }
  }

  // Verify signature with server
  async verifySignature(walletAddress, signature, pubKey, signDoc) {
    try {
      console.log('Sending to backend:', { walletAddress, signature, pubKey, signDoc })

      const response = await axios.post(
        `${this.baseURL}/auth/verify-signature`,
        {
          walletAddress,
          signature,
          pubKey,
          signDoc,
          signMode: 'SIGN_MODE_AMINO', // Add signMode to indicate amino signing
        },
        {
          headers: {
            'Content-Type': 'application/json',
          },
        },
      )
      return response.data
    } catch (error) {
      console.error('Error verifying signature:', error)
      throw error
    }
  }

  // Complete wallet authentication
  async authenticateWallet(wallet) {
    if (!wallet) {
      throw new Error('No wallet provided')
    }
    const [account] = await wallet.getAccounts()
    const address = account.address

    // Request nonce from server
    const { nonce, message } = await this.requestNonce(address)
    console.log('Received nonce:', nonce)

    // Sign message using sign_amino
    const { signature, pubKey, signDoc } = await this.signWithWallet(message, wallet)

    // Send verification to backend
    return await this.verifySignature(address, signature, pubKey, signDoc)
  }

  

}

export default new CosmosServices()
