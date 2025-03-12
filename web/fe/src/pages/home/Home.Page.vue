<template>
  <q-layout>
    <ProductsComponent />
  </q-layout>
</template>

<script>

import ProductsComponent from 'src/components/products/Products.Component.vue';


export default {
  name: "HomePage",
  components: {
    ProductsComponent
  },
}

</script>
