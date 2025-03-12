import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'

module.exports.restoreWalletFromMnemonic = async (mnemonic) => {
  const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, { prefix: 'cosmos' })
  const [account] = await wallet.getAccounts()
  return { wallet, address: account.address }
}

module.exports.generateWallet = async () => {
  // Tạo ví 12 từ (hoặc 24 nếu muốn)
  const wallet = await DirectSecp256k1HdWallet.generate(24, {
    prefix: 'cosmos', // prefix của địa chỉ, ví dụ "cosmos", "osmo", "juno" v.v.
  })

  // Lấy mnemonic
  const mnemonic = wallet.mnemonic

  // Lấy danh sách tài khoản
  const [firstAccount] = await wallet.getAccounts()
  // Địa chỉ
  const address = firstAccount.address

  console.log('Mnemonic:', mnemonic)
  console.log('Address:', address)

  return { mnemonic, address, wallet }
}

module.exports.storeMnemonic = (mnemonic) => {
  // Ở đây là DEMO, bạn nên mã hoá mnemonic trước khi lưu
  localStorage.setItem('my_cosmos_mnemonic', mnemonic)
}

module.exports.getMnemonic = () => {
  return localStorage.getItem('my_cosmos_mnemonic')
}
