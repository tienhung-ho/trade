import { boot } from 'quasar/wrappers'
import axios from 'axios'
import { Cookies, Notify } from 'quasar'
// import { useAuthStore } from 'src/stores/auth'

// Create main axios instance
const api = axios.create({
  baseURL: `${import.meta.env.VITE_API_URL}/${import.meta.env.VITE_API_NAME}/${import.meta.env.VITE_API_VERSION}/auth`,
  withCredentials: true,
})

// Variable to track if token is being refreshed
let isRefreshing = false
// Queue to hold requests waiting for token refresh
let failedQueue = []

// Process the queued requests
const processQueue = (error, token = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(token)
    }
  })
  failedQueue = []
}

// Refresh token function
const refreshToken = async () => {
  const response = await api.post('/refresh-token')
  return response.data
}

export default boot(({ app, router }) => {
  // Destructure router from boot context
  // Add response interceptor
  api.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config

      if (error.response?.data?.status_code === 403) {
        // Create notification
        try {
          Notify.create({
            type: 'negative',
            message: 'Forbidden: You do not have permission to access this resource.',
            position: 'top',
            timeout: 3000,
          })

          console.log(error.response)

          router.push('/403') // Use router from boot context

          return Promise.reject(error)
        } catch (err) {
          return Promise.reject(err)
        }
      }

      if (error.response?.status === 401 && !originalRequest._retry) {
        if (isRefreshing) {
          // If token is being refreshed, add request to queue
          return new Promise((resolve, reject) => {
            failedQueue.push({ resolve, reject })
          })
            .then(() => {
              return api(originalRequest)
            })
            .catch((err) => Promise.reject(err))
        }

        originalRequest._retry = true
        isRefreshing = true

        try {
          const response = await refreshToken()
          const newAccessToken = response.access_token // Get access_token from response

          // const authStore = useAuthStore()
          // await authStore.updateAccessToken(newAccessToken) // Update store

          processQueue(null, newAccessToken)

          // Retry the original request
          return api(originalRequest)
        } catch (refreshError) {
          processQueue(refreshError, null)

          // If refresh token fails, logout user
          Cookies.remove('refresh_token') // or localStorage.removeItem('refresh_token')
          router.push('/auth/login-web2') // Use router from boot context

          Notify.create({
            type: 'negative',
            message: 'Your session has expired, please log in again.',
          })

          return Promise.reject(refreshError)
        } finally {
          isRefreshing = false
        }
      }

      // Handle other errors
      if (error.response) {
        switch (error.response.status) {
          case 401:
            Notify.create({
              type: 'negative',
              message: 'Unauthorized access.',
            })
            router.push('/auth/login-web2') // Use router from boot context
            break
          case 404:
            Notify.create({
              type: 'negative',
              message: 'Resource not found.',
            })
            break
          default:
            Notify.create({
              type: 'negative',
              message: error.response.data?.message || 'An error has occurred.',
            })
        }
      } else if (error.request) {
        Notify.create({
          type: 'negative',
          message: 'Unable to connect to the server.',
        })
      } else {
        Notify.create({
          type: 'negative',
          message: 'An error occurred while sending the request.',
        })
      }

      return Promise.reject(error)
    },
  )

  // Make Axios available globally
  app.config.globalProperties.$axios = api
})

export { api as axios }
