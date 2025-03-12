<template>
  <q-page class="bg-white">
    <!-- Tiêu đề trang -->
    <div class="text-h4 text-center q-mt-xl q-mb-md">
      Top Products
    </div>

    <!-- Grid hiển thị sản phẩm -->
    <div class="row q-col-gutter-md q-px-xl q-pb-xl">
      <!-- Hiển thị sản phẩm từ API (đã lọc thông tin cần hiển thị) -->
      <div v-for="(item, index) in pagedProducts" :key="item.product_id || index"
        class="col-12 col-sm-6 col-md-4 col-lg-3 q-mb-md">
        <q-card class="bg-dark text-white relative-position card--product">
          <!-- Ảnh sản phẩm: nếu có trường image, nếu không thì dùng placeholder -->
          <q-img :src="item.image" :alt="item.alt_text" ratio="1" style="object-fit: contain;" class="q-pa-md" />

          <!-- Thông tin sản phẩm -->
          <q-card-section>
            <div class="text-h6 q-mb-xs">${{ item.price }}</div>
            <div class="text-subtitle1">{{ item.name }}</div>
            <!-- Hiển thị tên category nếu cần -->
            <div class="text-caption">Category: {{ item.categoryName }}</div>
          </q-card-section>

          <!-- Các nút hành động -->
          <q-card-actions align="around" class="q-pt-none q-pb-md">
            <q-btn icon="info" color="white" flat @click="showInfo(item)" />
            <q-btn label="Add" color="primary" @click="addToCart(item)" />
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- Phân trang -->
    <div class="row justify-center q-my-md" v-if="pageCount > 1">
      <q-pagination v-model="currentPage" :max="pageCount" color="primary" boundary-links
        @update:model-value="onPageChange" />
    </div>
  </q-page>
</template>

<script>
import productServices from 'src/services/product/product.services'

export default {
  name: 'ProductComponent',
  data() {
    return {
      products: [], // Dữ liệu sản phẩm đã được lọc
      currentPage: 1,
      itemsPerPage: 6,
      total: 0 // Tổng số sản phẩm để tính số trang
    }
  },
  computed: {
    // Tính số trang dựa trên tổng số sản phẩm
    pageCount() {
      return Math.ceil(this.total / this.itemsPerPage)
    },
    // Trong trường hợp API đã phân trang, bạn có thể dùng products trực tiếp
    pagedProducts() {
      return this.products
    }
  },
  methods: {
    async fetchProducts() {
      try {
        // Gọi API với các tham số page và limit
        const response = await productServices.listProduct({
          page: this.currentPage,
          limit: this.itemsPerPage
        })

        const data = response.data
        const paging = response.paging
        this.total = paging.total

        // Lọc thông tin cần hiển thị từ từng sản phẩm
        this.products = data.map(item => {
          return {
            product_id: item.product_id,
            name: item.name,
            price: item.price,
            description: item.description,
            status: item.status,
            // Nếu images tồn tại và có phần tử thì lấy url, nếu không dùng placeholder
            image: item.images && item.images.length > 0 ? item.images[0].url : 'https://via.placeholder.com/150',
            alt_text: item.images && item.images.length > 0 ? item.images[0].alt_text : "",
            categoryName: item.category ? item.category.name : ''
          }
        })
      } catch (error) {
        console.error('Error fetching products:', error)
      }
    },

    onPageChange(page) {
      this.currentPage = page
      this.fetchProducts()
    },

    // Hàm hiển thị dialog với HTML thông tin sản phẩm
    showInfo(item) {
      this.$q.dialog({
        title: `<strong>${item.name}</strong>`,
        message: `
          <div>
            <p><strong>Price:</strong> $${item.price}</p>
            <p><strong>Description:</strong> ${item.description}</p>
            <p><strong>Status:</strong> ${item.status}</p>
            <p><strong>Category:</strong> ${item.categoryName}</p>
          </div>
        `,
        html: true
      }).onOk(() => {
        console.log('OK clicked')
      }).onCancel(() => {
        console.log('Cancel clicked')
      }).onDismiss(() => {
        console.log('Dialog dismissed')
      })
    },

    addToCart(item) {
      this.$q.notify({
        message: `Đã thêm "${item.name}" vào giỏ hàng!`,
        color: 'green',
        position: 'top'
      })
    }
  },
  mounted() {
    this.fetchProducts()
  }
}
</script>

<style scoped>
.card--product {
  background-blend-mode: multiply;
  background-size: cover;
  background-position: center;
  border-radius: 8px;
}
</style>
