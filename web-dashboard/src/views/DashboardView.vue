<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-bold text-gray-800">SEO Tech Platform</h1>
          </div>
          <div class="flex items-center space-x-4">
            <router-link to="/projects" class="text-gray-600 hover:text-gray-900">
              Projects
            </router-link>
            <button @click="handleLogout" class="text-red-600 hover:text-red-700">
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <h2 class="text-2xl font-bold text-gray-800 mb-6">Dashboard</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-gray-500 text-sm font-medium">Total Projects</h3>
            <p class="text-3xl font-bold text-gray-800 mt-2">{{ totalProjects }}</p>
          </div>
          
          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-gray-500 text-sm font-medium">Active Audits</h3>
            <p class="text-3xl font-bold text-blue-600 mt-2">{{ activeAudits }}</p>
          </div>
          
          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-gray-500 text-sm font-medium">Avg Score</h3>
            <p class="text-3xl font-bold text-green-600 mt-2">{{ averageScore }}</p>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow">
          <div class="px-6 py-4 border-b border-gray-200">
            <h3 class="text-lg font-medium text-gray-800">Recent Projects</h3>
          </div>
          <div class="p-6">
            <div v-if="projectStore.projects.length === 0" class="text-gray-500 text-center py-8">
              No recent projects. <router-link to="/projects" class="text-blue-600 hover:text-blue-700">Create your first project</router-link>
            </div>
            <div v-else class="space-y-4">
              <div
                v-for="project in projectStore.projects.slice(0, 5)"
                :key="project.id"
                class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer"
                @click="router.push(`/projects/${project.id}`)"
              >
                <div>
                  <h4 class="text-lg font-medium text-gray-800">{{ project.name }}</h4>
                  <p class="text-sm text-gray-500">{{ project.domain }}</p>
                </div>
                <div class="text-sm text-gray-400">
                  {{ new Date(project.created_at).toLocaleDateString() }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useProjectStore } from '@/store/project'

const router = useRouter()
const authStore = useAuthStore()
const projectStore = useProjectStore()

const totalProjects = ref(0)
const activeAudits = ref(0)
const averageScore = ref(0)

onMounted(async () => {
  await projectStore.fetchProjects()
  totalProjects.value = projectStore.projects.length
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
