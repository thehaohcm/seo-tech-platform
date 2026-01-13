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
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <!-- Audit Info -->
        <div class="bg-white rounded-lg shadow p-6 mb-6">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h1 class="text-2xl font-bold text-gray-800">Audit Results</h1>
              <p class="text-gray-600">Audit ID: {{ $route.params.id }}</p>
            </div>
            <div v-if="auditRun" class="flex items-center space-x-4">
              <span
                :class="{
                  'px-3 py-1 rounded-full text-sm font-medium': true,
                  'bg-yellow-100 text-yellow-800': auditRun.status === 'queued' || auditRun.status === 'running',
                  'bg-green-100 text-green-800': auditRun.status === 'completed',
                  'bg-red-100 text-red-800': auditRun.status === 'failed'
                }"
              >
                {{ auditRun.status }}
              </span>
              <button
                v-if="auditRun.status === 'running' || auditRun.status === 'queued'"
                @click="fetchAuditData"
                class="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
              >
                Refresh
              </button>
            </div>
          </div>

          <div v-if="auditRun" class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p class="text-sm text-gray-500">Started At</p>
              <p class="font-medium">{{ formatDate(auditRun.started_at) }}</p>
            </div>
            <div v-if="auditRun.finished_at">
              <p class="text-sm text-gray-500">Finished At</p>
              <p class="font-medium">{{ formatDate(auditRun.finished_at) }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Total Pages</p>
              <p class="font-medium">{{ pageAudits.length }}</p>
            </div>
            <div v-if="auditRun.overall_score">
              <p class="text-sm text-gray-500">Overall Score</p>
              <p class="font-medium text-green-600">{{ auditRun.overall_score }}</p>
            </div>
          </div>
        </div>

        <!-- Loading State -->
        <div v-if="loading && pageAudits.length === 0" class="bg-white rounded-lg shadow p-12 text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p class="text-gray-500">Loading audit results...</p>
        </div>

        <!-- Page Audits List -->
        <div v-else-if="pageAudits.length > 0" class="space-y-4">
          <div
            v-for="page in pageAudits"
            :key="page.id"
            class="bg-white rounded-lg shadow p-6"
          >
            <div class="flex items-start justify-between mb-4">
              <div class="flex-1">
                <h3 class="text-lg font-semibold text-gray-800 mb-2">{{ page.title || 'Untitled Page' }}</h3>
                <a
                  :href="page.url"
                  target="_blank"
                  class="text-blue-600 hover:text-blue-700 text-sm break-all"
                >
                  {{ page.url }}
                </a>
              </div>
              <span
                :class="{
                  'px-3 py-1 rounded-full text-sm font-medium': true,
                  'bg-green-100 text-green-800': page.status_code === 200,
                  'bg-red-100 text-red-800': page.status_code !== 200
                }"
              >
                {{ page.status_code || 'N/A' }}
              </span>
            </div>

            <!-- Performance Metrics -->
            <div class="grid grid-cols-2 md:grid-cols-5 gap-4 mb-4">
              <div>
                <p class="text-xs text-gray-500">LCP</p>
                <p class="text-sm font-medium">{{ formatScore(page.lcp_score) }}s</p>
              </div>
              <div>
                <p class="text-xs text-gray-500">FCP</p>
                <p class="text-sm font-medium">{{ formatScore(page.fcp_score) }}s</p>
              </div>
              <div>
                <p class="text-xs text-gray-500">CLS</p>
                <p class="text-sm font-medium">{{ formatScore(page.cls_score) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-500">Load Time</p>
                <p class="text-sm font-medium">{{ page.load_time_ms || 'N/A' }}ms</p>
              </div>
              <div>
                <p class="text-xs text-gray-500">Created</p>
                <p class="text-sm font-medium">{{ formatTime(page.created_at) }}</p>
              </div>
            </div>

            <!-- Meta Description -->
            <div v-if="page.meta_description" class="mb-4">
              <p class="text-xs text-gray-500 mb-1">Meta Description</p>
              <p class="text-sm text-gray-700">{{ page.meta_description }}</p>
            </div>

            <!-- AI Suggestions -->
            <div v-if="page.ai_suggestions && !page.ai_suggestions.includes('Error')" class="bg-blue-50 p-4 rounded">
              <p class="text-xs text-blue-600 font-medium mb-2">AI Suggestions</p>
              <p class="text-sm text-gray-700">{{ page.ai_suggestions }}</p>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="bg-white rounded-lg shadow p-12 text-center">
          <p class="text-gray-500">No page audits found yet. The crawler is processing pages...</p>
          <button
            @click="fetchAuditData"
            class="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Refresh
          </button>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const auditRun = ref(null)
const pageAudits = ref([])
const loading = ref(true)
let pollingInterval = null

const fetchAuditData = async () => {
  try {
    const auditId = route.params.id
    const token = localStorage.getItem('token')
    
    // Fetch audit run
    const runResponse = await axios.get(`/api/v1/audits/${auditId}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    auditRun.value = runResponse.data.audit_run

    // Fetch page audits
    const pagesResponse = await axios.get(`/api/v1/audits/${auditId}/pages`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    pageAudits.value = pagesResponse.data.pages || []

    loading.value = false

    // Stop polling if audit is completed or failed
    if (auditRun.value.status === 'completed' || auditRun.value.status === 'failed') {
      stopPolling()
    }
  } catch (error) {
    console.error('Error fetching audit data:', error)
    loading.value = false
  }
}

const startPolling = () => {
  // Poll every 3 seconds
  pollingInterval = setInterval(fetchAuditData, 3000)
}

const stopPolling = () => {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleString()
}

const formatTime = (dateString) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleTimeString()
}

const formatScore = (score) => {
  if (!score && score !== 0) return 'N/A'
  return typeof score === 'number' ? score.toFixed(2) : score
}

onMounted(() => {
  fetchAuditData()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>
