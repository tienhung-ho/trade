<template>
  <q-layout>
    <q-page-container>
      <q-page padding class="flex flex-center">

        <div
          style="max-width: 40rem; width: 100%; box-shadow: rgba(99, 99, 99, 0.2) 0px 2px 8px 0px; background-color:#f5f5f5;"
          class="q-pa-lg">
          <!-- Network Selection -->
          <div class="q-my-md">
            <q-select filled v-model="selectedNetwork" :options="networks" option-label="name" emit-value map-options
              label="Select Blockchain Network" @update:model-value="updateNetwork" />
          </div>

          <!-- Mnemonic Input / Existing Wallet Restoration -->
          <div class="q-my-md">
            <q-input filled v-model="mnemonicInput" :type="showMnemonic ? 'text' : 'password'"
              label="Enter mnemonic (optional)" hint="If empty, the system will generate a new wallet"
              :disable="isConnecting">
              <template v-slot:append>
                <q-icon :name="showMnemonic ? 'visibility_off' : 'visibility'" class="cursor-pointer"
                  @click="toggleMnemonicVisibility" />
              </template>
            </q-input>
          </div>

          <!-- Buttons -->
          <div class="q-my-md flex flex-center">
            <!-- Connect existing wallet from mnemonic -->
            <q-btn color="black" class="q-mr-md" :loading="isConnecting" text-color="white" label="Connect Wallet"
              @click="connectWalletWithMnemonic" />

            <!-- Create new wallet -->
            <q-btn color="black" text-color="white" class="q-mr-md" :loading="isConnecting" label="Create New Wallet"
              @click="createNewWallet" />

            <!-- Authenticate Button -->
            <q-btn color="primary" text-color="white" :loading="isAuthenticating" :disable="!isConnected"
              label="Authenticate" @click="authenticateWallet" />
          </div>

          <!-- Authentication Status -->
          <q-banner v-if="authStatus" :class="authStatus.verified ? 'bg-positive' : 'bg-negative'" class="q-mt-md">
            <template v-slot:avatar>
              <q-icon :name="authStatus.verified ? 'check_circle' : 'error'" color="white" />
            </template>
            <div class="text-white">
              {{ authStatus.verified ? 'Authentication successful!' : 'Authentication failed!' }}
              {{ authStatus.message || '' }}
            </div>
          </q-banner>

          <!-- Display Wallet Info -->
          <q-card v-if="address" class="q-mt-lg bg-grey-1">
            <q-card-section>
              <div class="text-h6">Wallet Information</div>
              <div class="q-mt-sm">
                <div class="row items-center">
                  <div class="col-12 col-sm-3">
                    <strong>Address:</strong>
                  </div>
                  <div class="col-12 col-sm-9 ellipsis">
                    {{ address }}
                  </div>
                </div>

                <!-- Actual mnemonic is stored in 'mnemonic' but displayed as password. User can toggle visibility. -->
                <div v-if="mnemonic" class="row items-center q-mt-sm">
                  <div class="col-12 col-sm-3"><strong>Mnemonic:</strong></div>
                  <div class="col-12 col-sm-9">

                    <q-input filled v-model="mnemonic" :type="showMnemonic ? 'text' : 'password'" label="Mnemonic"
                      class="q-mb-md">
                      <template v-slot:append>
                        <q-icon :name="showMnemonic ? 'visibility_off' : 'visibility'" class="cursor-pointer"
                          @click="toggleMnemonicVisibility" />
                      </template>
                    </q-input>

                    <q-btn flat dense round icon="content_copy" class="q-ml-sm" @click="copyToClipboard(mnemonic)" />
                  </div>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- Balance info -->
          <q-card v-if="balances.length > 0" class="q-mt-md bg-grey-1">
            <q-card-section>
              <div class="text-h6">Wallet Balances</div>
              <q-list>
                <q-item v-for="(bal, idx) in balances" :key="idx">
                  <q-item-section>
                    <q-item-label>
                      {{ formatDenom(bal.denom) }}
                    </q-item-label>
                    <q-item-label caption>
                      {{ formatAmount(bal.amount, bal.denom) }}
                    </q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-card-section>
          </q-card>
        </div>

      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script>
import { ref, computed } from 'vue'
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing'
import { SigningStargateClient } from '@cosmjs/stargate'
import { useQuasar } from 'quasar'
import cosmosServices from 'src/services/cosmos/cosmos.services'

export default {
  name: 'SelfCustodyWallet',

  setup() {
    const $q = useQuasar()

    // List of networks
    const networks = [
      {
        name: 'mytoken',
        rpcEndpoint: 'http://localhost:26657',
        prefix: 'cosmos',
        denomDefault: 'citcoint',
        denomName: 'CITCOINT',
        denomDecimals: 6,
      }
    ]
    const selectedNetwork = ref(networks[0])

    // States
    const mnemonic = ref('')
    const mnemonicInput = ref('')
    const address = ref('')
    const client = ref(null)
    const wallet = ref(null)
    const balances = ref([])
    const isConnecting = ref(false)
    const isLoadingBalance = ref(false)
    const showMnemonic = ref(false)

    // Authentication state
    const isAuthenticating = ref(false)
    const authStatus = ref(null)

    const isConnected = computed(() => !!client.value && !!address.value)

    // Switch network
    function updateNetwork(network) {
      selectedNetwork.value = network
      // reset
      client.value = null
      wallet.value = null
      address.value = ''
      balances.value = []
      authStatus.value = null
    }

    // Toggle mnemonic password/text
    function toggleMnemonicVisibility() {
      showMnemonic.value = !showMnemonic.value
    }

    // 1) Create a new wallet
    async function createNewWallet() {
      if (isConnecting.value) return
      isConnecting.value = true

      try {
        const newWallet = await DirectSecp256k1HdWallet.generate(24, {
          prefix: selectedNetwork.value.prefix,
        })
        mnemonic.value = newWallet.mnemonic
        wallet.value = newWallet

        const [account] = await newWallet.getAccounts()
        address.value = account.address

        client.value = await SigningStargateClient.connectWithSigner(
          selectedNetwork.value.rpcEndpoint,
          newWallet
        )

        $q.notify({
          color: 'positive',
          message: 'Successfully created a new wallet!',
          icon: 'check'
        })

        // Call faucet
        await faucetAndNotify(address.value)

        // Retrieve balance
        getBalance()
      } catch (err) {
        console.error(err)
        $q.notify({
          color: 'negative',
          message: 'Error while creating a new wallet: ' + err.message,
          icon: 'error'
        })
      } finally {
        isConnecting.value = false
      }
    }

    // 2) Connect an existing wallet (restoring from user mnemonic)
    async function connectWalletWithMnemonic() {
      if (isConnecting.value) return
      isConnecting.value = true

      try {
        const input = mnemonicInput.value.trim()
        if (!input) {
          $q.notify({
            color: 'warning',
            message: 'No mnemonic entered!',
            icon: 'warning'
          })
          return
        }

        const existingWallet = await DirectSecp256k1HdWallet.fromMnemonic(input, {
          prefix: selectedNetwork.value.prefix,
        })
        mnemonic.value = input
        wallet.value = existingWallet

        const [account] = await existingWallet.getAccounts()
        address.value = account.address

        client.value = await SigningStargateClient.connectWithSigner(
          selectedNetwork.value.rpcEndpoint,
          existingWallet
        )

        $q.notify({
          color: 'positive',
          message: 'Successfully connected/imported wallet!',
          icon: 'check'
        })

        // Retrieve balance
        getBalance()
      } catch (err) {
        console.error(err)
        $q.notify({
          color: 'negative',
          message: 'Error restoring wallet: ' + err.message,
          icon: 'error'
        })
      } finally {
        isConnecting.value = false
      }
    }

    // Authenticate wallet
    async function authenticateWallet() {
      if (!wallet.value || !address.value) {
        $q.notify({
          color: 'warning',
          message: 'Please connect a wallet first!',
          icon: 'warning'
        })
        return
      }

      isAuthenticating.value = true
      authStatus.value = null

      try {
        // Call authentication function from service
        const result = await cosmosServices.authenticateWallet(wallet.value)

        // Display result
        authStatus.value = {
          verified: true,
          message: `Address ${result.address} verified successfully!`
        }

        console.log(result);

        $q.notify({
          color: 'positive',
          message: 'Authentication successful!',
          icon: 'check'
        })
      } catch (err) {
        console.error('Authentication error:', err)

        authStatus.value = {
          verified: false,
          message: err.message
        }

        $q.notify({
          color: 'negative',
          message: 'Authentication failed: ' + err.message,
          icon: 'error'
        })
      } finally {
        isAuthenticating.value = false
      }
    }

    // Call faucet (optional)
    async function faucetAndNotify(addr) {
      try {
        const faucetResult = await cosmosServices.faucet({ address: addr })
        console.log('Faucet result:', faucetResult)
        $q.notify({
          color: 'positive',
          message: `Faucet success! TxHash: ${faucetResult.txHash || ''}`,
          icon: 'check'
        })
      } catch (err) {
        console.error('Faucet error:', err)
        $q.notify({
          color: 'negative',
          message: 'Faucet error: ' + err.message,
          icon: 'error'
        })
      }
    }

    // Get balances
    async function getBalance() {
      if (!client.value || !address.value) {
        $q.notify({
          color: 'warning',
          message: 'No client or address!',
          icon: 'warning'
        })
        return
      }
      isLoadingBalance.value = true
      try {
        const allBalances = await client.value.getAllBalances(address.value)
        balances.value = allBalances
        console.log('Balances:', allBalances)

        if (allBalances.length === 0) {
          $q.notify({
            color: 'info',
            message: 'This wallet has no tokens.',
            icon: 'info'
          })
        }
      } catch (err) {
        console.error(err)
        $q.notify({
          color: 'negative',
          message: 'Error fetching balance: ' + err.message,
          icon: 'error'
        })
      } finally {
        isLoadingBalance.value = false
      }
    }

    // Format denom
    function formatDenom(denom) {
      if (denom.startsWith('u')) {
        return denom.substring(1).toUpperCase()
      }
      return denom.toUpperCase()
    }

    function formatAmount(amount) {
      const decimals = 6
      const value = parseFloat(amount) / Math.pow(10, decimals)
      return value.toLocaleString('en-US', { maximumFractionDigits: decimals })
    }

    function copyToClipboard(text) {
      navigator.clipboard.writeText(text).then(() => {
        $q.notify({
          color: 'positive',
          message: 'Copied to clipboard',
          icon: 'content_copy'
        })
      })
    }

    return {
      // Networks
      networks,
      selectedNetwork,

      // States
      mnemonic,
      mnemonicInput,
      address,
      client,
      balances,
      showMnemonic,
      isConnecting,
      isLoadingBalance,
      isConnected,

      // Authentication state
      isAuthenticating,
      authStatus,

      updateNetwork,
      toggleMnemonicVisibility,
      createNewWallet,
      connectWalletWithMnemonic,
      getBalance,
      faucetAndNotify,
      authenticateWallet,

      formatDenom,
      formatAmount,
      copyToClipboard,
    }
  }
}
</script>
