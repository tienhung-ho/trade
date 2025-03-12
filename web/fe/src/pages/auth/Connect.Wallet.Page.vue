<template>
 <q-layout>
  <q-page-container>
    <q-page>
       <ConnectWalletComponent></ConnectWalletComponent>
     </q-page>

  </q-page-container>

 </q-layout>
</template>

<script>
import ConnectWalletComponent from 'src/components/auth/Connect.Wallet.Component.vue';


export default {
  name: "ConnectWallet",
  components: ConnectWalletComponent,

}


</script>
