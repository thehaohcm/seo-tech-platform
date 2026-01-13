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
            class="bg-white rounded-lg shadow hover:shadow-lg transition-shadow relative"
          >
            <div class="p-6 cursor-pointer" @click="goToProject(project.id)">
              <h3 class="text-lg font-semibold text-gray-800">{{ project.name }}</h3>
              <p class="text-sm text-gray-500 mt-1">{{ project.domain }}</p>
              <div class="mt-4 flex items-center text-xs text-gray-400">
                <span>Created {{ formatDate(project.created_at) }}</span>
              </div>
            </div>
            <div class="absolute top-4 right-4 flex gap-2">
              <button
                @click.stop="confirmDelete(project)"
                class="p-2 text-red-600 hover:bg-red-50 rounded-full transition"
                title="Delete project"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
              </button>
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

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-md">
        <div class="flex items-center gap-3 mb-4">
          <div class="flex-shrink-0 w-10 h-10 rounded-full bg-red-100 flex items-center justify-center">
            <svg class="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Delete Project</h3>
          </div>
        </div>
        
        <p class="text-sm text-gray-600 mb-6">
          Are you sure you want to delete <strong class="text-gray-900">{{ projectToDelete?.name }}</strong>? 
          This action cannot be undone and will permanently delete all audit data associated with this project.
        </p>

        <div class="flex justify-end space-x-3">
          <button
            type="button"
            @click="showDeleteModal = false"
            :disabled="deleting"
            class="px-4 py-2 text-gray-700 hover:text-gray-900 disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            type="button"
            @click="handleDeleteProject"
            :disabled="deleting"
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:bg-red-400 flex items-center gap-2"
          >
            <svg v-if="deleting" class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ deleting ? 'Deleting...' : 'Delete Project' }}
          </button>
        </div>
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
const showDeleteModal = ref(false)
const projectToDelete = ref(null)
const deleting = ref(false)
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

function confirmDelete(project) {
  projectToDelete.value = project
  showDeleteModal.value = true
}

async function handleDeleteProject() {
  if (!projectToDelete.value) return
  
  deleting.value = true
  try {
    await projectStore.deleteProject(projectToDelete.value.id)
    showDeleteModal.value = false
    projectToDelete.value = null
  } catch (error) {
    alert('Failed to delete project: ' + (error.response?.data?.error || error.message))
  } finally {
    deleting.value = false
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

function formatDate(date) {
  return new Date(date).toLocaleDateString()
}
</script>
