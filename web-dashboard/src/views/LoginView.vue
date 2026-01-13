<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-800">SEO Tech Platform</h1>
        <p class="text-gray-600 mt-2">{{ isRegister ? 'Create an account' : 'Sign in to your account' }}</p>
      </div>

      <form @submit.prevent="isRegister ? handleRegister() : handleLogin()" class="space-y-6">
        <div v-if="isRegister">
          <label class="block text-sm font-medium text-gray-700">Full Name</label>
          <input
            v-model="fullName"
            type="text"
            required
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700">Email</label>
          <input
            v-model="email"
            type="email"
            required
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700">Password</label>
          <input
            v-model="password"
            type="password"
            required
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <div v-if="error" class="text-red-600 text-sm">
          {{ error }}
        </div>

        <button
          type="submit"
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          {{ isRegister ? 'Register' : 'Sign In' }}
        </button>

        <div class="text-center mt-4">
          <p class="text-sm text-gray-600">
            {{ isRegister ? 'Already have an account?' : "Don't have an account?" }}
            <button type="button" @click="isRegister = !isRegister" class="text-blue-600 hover:text-blue-700 font-medium">
              {{ isRegister ? 'Sign in here' : 'Register here' }}
            </button>
          </p>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()

const isRegister = ref(false)
const fullName = ref('')
const email = ref('')
const password = ref('')
const error = ref('')

async function handleLogin() {
  error.value = ''
  const success = await authStore.login(email.value, password.value)
  
  if (success) {
    router.push('/')
  } else {
    error.value = 'Invalid email or password'
  }
}

async function handleRegister() {
  error.value = ''
  const success = await authStore.register(email.value, password.value, fullName.value)
  
  if (success) {
    router.push('/')
  } else {
    error.value = 'Registration failed. Email may already exist.'
  }
}
</script>
