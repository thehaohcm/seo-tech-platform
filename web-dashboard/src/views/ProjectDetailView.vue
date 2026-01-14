<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center space-x-4">
            <router-link to="/" class="text-xl font-bold text-gray-800">
              SEO Tech Platform
            </router-link>
            <span class="text-gray-400">/</span>
            <router-link to="/projects" class="text-gray-600 hover:text-gray-900">
              Projects
            </router-link>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <!-- Back Button -->
        <button
          @click="goBack"
          class="mb-4 flex items-center gap-2 text-gray-600 hover:text-gray-900 transition"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z" clip-rule="evenodd" />
          </svg>
          <span class="font-medium">Back</span>
        </button>
        
        <div v-if="loading" class="text-center py-12">
          <p class="text-gray-500">Loading project...</p>
        </div>

        <div v-else-if="project">
          <div class="bg-white rounded-lg shadow p-6 mb-6">
            <div class="flex justify-between items-start">
              <div>
                <h1 class="text-2xl font-bold text-gray-800">{{ project.name }}</h1>
                <p class="text-gray-600 mt-1">{{ project.domain }}</p>
              </div>
              <div class="flex gap-3">
                <button
                  @click="confirmDelete"
                  class="px-4 py-2 text-red-600 border border-red-600 hover:bg-red-50 rounded-md transition"
                  title="Delete project"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 inline-block" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </button>
                <button
                  @click="handleStartAudit"
                  class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md"
                >
                  Start Full Audit
                </button>
              </div>
            </div>
            
            <!-- Single URL Audit Section -->
            <div class="mt-6 pt-6 border-t border-gray-200">
              <h3 class="text-sm font-medium text-gray-700 mb-3">Audit Specific URL</h3>
              <div class="flex gap-3">
                <input
                  v-model="singleUrl"
                  type="url"
                  placeholder="https://example.com/specific-page"
                  class="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <button
                  @click="handleSingleUrlAudit"
                  :disabled="!singleUrl || auditingSingle"
                  class="bg-green-600 hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white px-6 py-2 rounded-md transition"
                >
                  {{ auditingSingle ? 'Auditing...' : 'Audit URL' }}
                </button>
              </div>
              <p class="text-xs text-gray-500 mt-2">Enter a specific URL to audit just that page instead of crawling the entire site.</p>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow">
            <div class="px-6 py-4 border-b border-gray-200">
              <h2 class="text-lg font-medium text-gray-800">Audit History</h2>
            </div>
            <div class="p-6">
              <!-- Loading State -->
              <div v-if="loadingAudits" class="text-center py-8">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-2"></div>
                <p class="text-gray-500 text-sm">Loading audits...</p>
              </div>
              
              <!-- Audit List -->
              <div v-else-if="auditRuns.length > 0" class="space-y-3">
                <div
                  v-for="audit in auditRuns"
                  :key="audit.id"
                  class="border border-gray-200 rounded-lg p-4 hover:border-blue-300 transition cursor-pointer"
                  @click="viewAudit(audit.id)"
                >
                  <div class="flex items-center justify-between">
                    <div>
                      <p class="font-medium text-gray-800">Audit #{{ audit.id }}</p>
                      <p class="text-sm text-gray-500">{{ formatDate(audit.started_at) }}</p>
                    </div>
                    <span
                      :class="{
                        'px-3 py-1 rounded-full text-sm font-medium': true,
                        'bg-yellow-100 text-yellow-800': audit.status === 'queued' || audit.status === 'running',
                        'bg-green-100 text-green-800': audit.status === 'completed',
                        'bg-red-100 text-red-800': audit.status === 'failed'
                      }"
                    >
                      {{ audit.status }}
                    </span>
                  </div>
                </div>
              </div>
              
              <!-- Empty State -->
              <p v-else class="text-gray-500 text-center py-8">No audits yet</p>
            </div>
          </div>
        </div>
      </div>
    </main>

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
          Are you sure you want to delete <strong class="text-gray-900">{{ project?.name }}</strong>? 
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
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50 flex items-center gap-2"
          >
            <svg v-if="deleting" class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ deleting ? 'Deleting...' : 'Delete' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/store/project'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()

const project = ref(null)
const loading = ref(true)
const showDeleteModal = ref(false)
const deleting = ref(false)
const singleUrl = ref('')
const auditingSingle = ref(false)
const auditRuns = ref([])
const loadingAudits = ref(false)

const goBack = () => {
  router.back()
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleString()
}

const viewAudit = (auditId) => {
  router.push(`/audits/${auditId}`)
}

async function fetchAuditHistory() {
  if (!project.value) return
  
  loadingAudits.value = true
  try {
    const token = localStorage.getItem('token')
    console.log('Fetching audit history for project:', project.value.id)
    const response = await projectStore.fetchAuditHistory(project.value.id)
    console.log('Audit history response:', response)
    auditRuns.value = response || []
    console.log('Audit runs set to:', auditRuns.value)
  } catch (error) {
    console.error('Failed to fetch audit history:', error)
    auditRuns.value = []
  } finally {
    loadingAudits.value = false
  }
}

onMounted(async () => {
  project.value = await projectStore.fetchProject(route.params.id)
  loading.value = false
  
  // Fetch audit history after project is loaded
  if (project.value) {
    await fetchAuditHistory()
  }
})

async function handleStartAudit() {
  try {
    const auditRun = await projectStore.startAudit(project.value.id)
    router.push(`/audits/${auditRun.id}`)
  } catch (error) {
    alert('Failed to start audit')
  }
}

async function handleSingleUrlAudit() {
  if (!singleUrl.value) {
    alert('Please enter a URL')
    return
  }
  
  auditingSingle.value = true
  try {
    const auditRun = await projectStore.startSingleUrlAudit(project.value.id, singleUrl.value)
    router.push(`/audits/${auditRun.id}`)
  } catch (error) {
    alert('Failed to start audit: ' + (error.response?.data?.error || error.message))
  } finally {
    auditingSingle.value = false
  }
}

function confirmDelete() {
  showDeleteModal.value = true
}

async function handleDeleteProject() {
  if (!project.value) return
  
  deleting.value = true
  try {
    await projectStore.deleteProject(project.value.id)
    router.push('/projects')
  } catch (error) {
    alert('Failed to delete project: ' + (error.response?.data?.error || error.message))
  } finally {
    deleting.value = false
  }
}
</script>
