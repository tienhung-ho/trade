import { axios } from '../../boot/axios'

class AuthServices {
  constructor() {
    this.baseURL = `${import.meta.env.VITE_API_URL}/${import.meta.env.VITE_API_NAME}/${import.meta.env.VITE_API_VERSION}`
  }

  async login(data) {
    try {
      const response = await axios.post(`${this.baseURL}/auth/login-web-2`, data)
      return response.data
    } catch (error) {
      if (error.response) {
        return error.response
      }
      throw error
    }
  }
}

export default new AuthServices()
