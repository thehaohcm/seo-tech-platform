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
                <div class="flex items-center gap-2 mb-2">
                  <h3 class="text-lg font-semibold text-gray-800">{{ page.title || 'Untitled Page' }}</h3>
                  <span
                    :class="{
                      'px-3 py-1 rounded-full text-sm font-bold': true,
                      'bg-green-100 text-green-700': getSEOScore(page).rating === 'Good',
                      'bg-yellow-100 text-yellow-700': getSEOScore(page).rating === 'Average',
                      'bg-red-100 text-red-700': getSEOScore(page).rating === 'Poor'
                    }"
                  >
                    {{ getSEOScore(page).rating }}
                  </span>
                </div>
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

            <!-- Auto Test Section -->
            <div class="border-t pt-4 mb-4">
              <div class="flex items-center justify-between mb-3">
                <h4 class="text-sm font-semibold text-gray-700">Automated Testing</h4>
                <button
                  @click="generateAutoTest(page)"
                  :disabled="testStatus[page.id]?.loading"
                  class="px-4 py-2 bg-purple-600 text-white text-sm rounded hover:bg-purple-700 disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <svg v-if="testStatus[page.id]?.loading" class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ testStatus[page.id]?.loading ? 'Testing...' : 'Generate Auto Test' }}
                </button>
              </div>
              
              <!-- Test Results -->
              <div v-if="testStatus[page.id]?.result" class="bg-gray-50 rounded-lg p-4">
                <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                  <div>
                    <p class="text-xs text-gray-500">Test Status</p>
                    <span
                      :class="{
                        'px-2 py-1 rounded text-xs font-medium inline-block': true,
                        'bg-green-100 text-green-800': testStatus[page.id].result.status === 'passed',
                        'bg-red-100 text-red-800': testStatus[page.id].result.status === 'failed',
                        'bg-yellow-100 text-yellow-800': testStatus[page.id].result.status === 'warning'
                      }"
                    >
                      {{ testStatus[page.id].result.status }}
                    </span>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Tests Run</p>
                    <p class="text-sm font-medium">{{ testStatus[page.id].result.total_tests }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Passed</p>
                    <p class="text-sm font-medium text-green-600">{{ testStatus[page.id].result.passed }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Failed</p>
                    <p class="text-sm font-medium text-red-600">{{ testStatus[page.id].result.failed }}</p>
                  </div>
                </div>

                <!-- Screenshot -->
                <div v-if="testStatus[page.id].result.screenshot_url" class="mt-4">
                  <p class="text-xs text-gray-500 mb-2">Screenshot</p>
                  <img 
                    :src="testStatus[page.id].result.screenshot_url" 
                    alt="Page Screenshot"
                    class="w-full max-w-2xl border rounded cursor-pointer hover:shadow-lg transition"
                    @click="openScreenshot(testStatus[page.id].result.screenshot_url)"
                  />
                </div>

                <!-- Test Details -->
                <div v-if="testStatus[page.id].result.test_details" class="mt-4">
                  <p class="text-xs text-gray-500 mb-2">Test Details</p>
                  <div class="space-y-2">
                    <div
                      v-for="(test, idx) in testStatus[page.id].result.test_details"
                      :key="idx"
                      class="flex items-start gap-2 text-sm"
                    >
                      <span v-if="test.passed" class="text-green-600">✓</span>
                      <span v-else class="text-red-600">✗</span>
                      <span>{{ test.name }}: {{ test.message }}</span>
                    </div>
                  </div>
                </div>
                
                <!-- Download Python Code Button -->
                <div v-if="testStatus[page.id].result.python_code" class="mt-4 pt-4 border-t border-gray-200">
                  <div class="flex items-center justify-between">
                    <div>
                      <p class="text-xs text-gray-500 mb-1">Playwright Page Object</p>
                      <p class="text-xs text-gray-400">Auto-generated Python code for test automation</p>
                    </div>
                    <button
                      @click="downloadPythonCode(page.id, page.url)"
                      class="px-4 py-2 bg-green-600 text-white text-sm rounded hover:bg-green-700 transition flex items-center gap-2"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
                      </svg>
                      Download Python Code
                    </button>
                  </div>
                </div>
              </div>

              <div v-else-if="testStatus[page.id]?.error" class="bg-red-50 p-4 rounded">
                <p class="text-sm text-red-600">{{ testStatus[page.id].error }}</p>
              </div>
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
          <div class="flex flex-col items-center">
            <!-- Animated Spinner -->
            <svg class="animate-spin h-12 w-12 text-blue-600 mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            
            <!-- Progress Message -->
            <p class="text-gray-700 font-medium mb-2">Processing Pages...</p>
            <p class="text-gray-500 text-sm mb-6">The crawler is analyzing pages and will update results automatically</p>
            
            <!-- Animated Progress Bar -->
            <div class="w-full max-w-xs h-2 bg-gray-200 rounded-full overflow-hidden mb-6">
              <div class="h-full bg-blue-600 rounded-full animate-progress"></div>
            </div>
            
            <button
              @click="fetchAuditData"
              class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition flex items-center gap-2"
            >
              <svg class="h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Refresh Now
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const auditRun = ref(null)
const pageAudits = ref([])
const loading = ref(true)
const testStatus = ref({})
let pollingInterval = null

const goBack = () => {
  router.back()
}

const generateAutoTest = async (page) => {
  const pageId = page.id
  
  // Set loading state
  testStatus.value[pageId] = { loading: true, result: null, error: null }
  
  try {
    const token = localStorage.getItem('token')
    await axios.post(
      `/api/v1/pages/${pageId}/generate-test`,
      { url: page.url },
      { headers: { Authorization: `Bearer ${token}` } }
    )
    
    // Poll for test results
    const pollInterval = setInterval(async () => {
      try {
        const resultResponse = await axios.get(
          `/api/v1/pages/${pageId}/test-result`,
          { headers: { Authorization: `Bearer ${token}` } }
        )
        
        if (resultResponse.data) {
          clearInterval(pollInterval)
          testStatus.value[pageId] = {
            loading: false,
            result: resultResponse.data,
            error: null
          }
        }
      } catch (pollError) {
        if (pollError.response?.status === 404) {
          // Still processing, continue polling
          return
        }
        
        console.error('Error polling test results:', pollError)
        clearInterval(pollInterval)
        testStatus.value[pageId] = {
          loading: false,
          result: null,
          error: 'Failed to retrieve test results'
        }
      }
    }, 2000) // Poll every 2 seconds
    
    // Stop polling after 2 minutes
    setTimeout(() => {
      clearInterval(pollInterval)
      if (testStatus.value[pageId]?.loading) {
        testStatus.value[pageId] = {
          loading: false,
          result: null,
          error: 'Test execution timeout'
        }
      }
    }, 120000)
    
  } catch (error) {
    console.error('Error generating auto test:', error)
    testStatus.value[pageId] = {
      loading: false,
      result: null,
      error: error.response?.data?.error || 'Failed to generate auto test'
    }
  }
}

const downloadPythonCode = (pageId, url) => {
  const pythonCode = testStatus.value[pageId]?.result?.python_code
  
  if (!pythonCode) {
    console.error('No Python code available')
    return
  }
  
  // Create blob from Python code
  const blob = new Blob([pythonCode], { type: 'text/x-python' })
  
  // Generate filename from URL
  const urlObj = new URL(url)
  const hostname = urlObj.hostname.replace(/\./g, '_')
  const pathname = urlObj.pathname.replace(/\//g, '_').replace(/\.$/, '') || 'index'
  const filename = `page_object_${hostname}${pathname}.py`
  
  // Create download link
  const downloadUrl = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = downloadUrl
  link.download = filename
  
  // Trigger download
  document.body.appendChild(link)
  link.click()
  
  // Cleanup
  document.body.removeChild(link)
  URL.revokeObjectURL(downloadUrl)
}

const openScreenshot = (url) => {
  window.open(url, '_blank')
}

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

const getSEOScore = (page) => {
  // Score thresholds based on Core Web Vitals
  const lcpScore = page.lcp_score || 0
  const fcpScore = page.fcp_score || 0
  const clsScore = page.cls_score || 0
  
  let goodCount = 0
  let poorCount = 0
  
  // LCP (Largest Contentful Paint): < 2.5s = good, 2.5-4s = average, > 4s = poor
  if (lcpScore < 2.5) goodCount++
  else if (lcpScore > 4) poorCount++
  
  // FCP (First Contentful Paint): < 1.8s = good, 1.8-3s = average, > 3s = poor
  if (fcpScore < 1.8) goodCount++
  else if (fcpScore > 3) poorCount++
  
  // CLS (Cumulative Layout Shift): < 0.1 = good, 0.1-0.25 = average, > 0.25 = poor
  if (clsScore < 0.1) goodCount++
  else if (clsScore > 0.25) poorCount++
  
  // Check for SEO issues
  const hasSEOIssues = page.seo_issues?.items?.length > 0 || false
  if (hasSEOIssues) poorCount++
  
  // Determine overall rating
  if (poorCount >= 2) {
    return { rating: 'Poor', score: 0 }
  } else if (goodCount >= 2) {
    return { rating: 'Good', score: 100 }
  } else {
    return { rating: 'Average', score: 50 }
  }
}

onMounted(() => {
  fetchAuditData()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>
