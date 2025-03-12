import { axios } from '../../boot/axios'

class ProductServices {
  constructor() {
    this.baseURL = `${import.meta.env.VITE_API_URL}/${import.meta.env.VITE_API_NAME}/${import.meta.env.VITE_API_VERSION}`
  }

  async listProduct(params = {}) {
    const response = await axios.get(`${this.baseURL}/product/list`, {
      params,
    })
    return response.data
  }
}

export default new ProductServices()
