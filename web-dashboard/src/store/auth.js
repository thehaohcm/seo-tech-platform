import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)

  const isAuthenticated = computed(() => !!token.value)

  async function login(email, password) {
    try {
      const response = await axios.post('/api/v1/auth/login', { email, password })
      user.value = response.data.user
      token.value = response.data.token
      localStorage.setItem('token', response.data.token)
      
      // Set axios default header
      axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.token}`
      
      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  }

  async function register(email, password, fullName) {
    try {
      const response = await axios.post('/api/v1/auth/register', { 
        email, 
        password, 
        full_name: fullName 
      })
      user.value = response.data.user
      token.value = response.data.token
      localStorage.setItem('token', response.data.token)
      
      // Set axios default header
      axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.token}`
      
      return true
    } catch (error) {
      console.error('Registration failed:', error)
      return false
    }
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    delete axios.defaults.headers.common['Authorization']
  }

  // Initialize axios header if token exists
  if (token.value) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return {
    user,
    token,
    isAuthenticated,
    login,
    register,
    logout
  }
})
