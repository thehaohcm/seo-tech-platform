<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/" class="text-xl font-bold text-gray-800">
              SEO Tech Platform
            </router-link>
          </div>
          <div class="flex items-center space-x-4">
            <router-link to="/projects" class="text-gray-900 font-medium">
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
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold text-gray-800">Projects</h2>
          <button
            @click="showCreateModal = true"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md"
          >
            + New Project
          </button>
        </div>

        <div v-if="loading" class="text-center py-12">
          <p class="text-gray-500">Loading projects...</p>
        </div>

        <div v-else-if="projects.length === 0" class="bg-white rounded-lg shadow p-12 text-center">
          <p class="text-gray-500 mb-4">No projects yet</p>
          <button
            @click="showCreateModal = true"
            class="text-blue-600 hover:text-blue-700"
          >
            Create your first project
          </button>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="project in projects"
            :key="project.id"
            class="bg-white rounded-lg shadow hover:shadow-lg transition-shadow cursor-pointer"
            @click="goToProject(project.id)"
          >
            <div class="p-6">
              <h3 class="text-lg font-semibold text-gray-800">{{ project.name }}</h3>
              <p class="text-sm text-gray-500 mt-1">{{ project.domain }}</p>
              <div class="mt-4 flex items-center text-xs text-gray-400">
                <span>Created {{ formatDate(project.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Create Project Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
      <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-md">
        <h3 class="text-lg font-semibold mb-4">Create New Project</h3>
        <form @submit.prevent="handleCreateProject">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Project Name</label>
              <input
                v-model="newProject.name"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Domain</label>
              <input
                v-model="newProject.domain"
                type="url"
                required
                placeholder="https://example.com"
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
              />
            </div>
          </div>
          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 text-gray-700 hover:text-gray-900"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Create
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useProjectStore } from '@/store/project'

const router = useRouter()
const authStore = useAuthStore()
const projectStore = useProjectStore()

const showCreateModal = ref(false)
const newProject = ref({
  name: '',
  domain: '',
  settings: {}
})

const projects = computed(() => projectStore.projects)
const loading = computed(() => projectStore.loading)

onMounted(() => {
  projectStore.fetchProjects()
})

async function handleCreateProject() {
  try {
    await projectStore.createProject(newProject.value)
    showCreateModal.value = false
    newProject.value = { name: '', domain: '', settings: {} }
  } catch (error) {
    alert('Failed to create project')
  }
}

function goToProject(id) {
  router.push(`/projects/${id}`)
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

function formatDate(date) {
  return new Date(date).toLocaleDateString()
}
</script>
