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
              <button
                @click="handleStartAudit"
                class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md"
              >
                Start Audit
              </button>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow">
            <div class="px-6 py-4 border-b border-gray-200">
              <h2 class="text-lg font-medium text-gray-800">Audit History</h2>
            </div>
            <div class="p-6">
              <p class="text-gray-500 text-center py-8">No audits yet</p>
            </div>
          </div>
        </div>
      </div>
    </main>
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

onMounted(async () => {
  project.value = await projectStore.fetchProject(route.params.id)
  loading.value = false
})

async function handleStartAudit() {
  try {
    const auditRun = await projectStore.startAudit(project.value.id)
    router.push(`/audits/${auditRun.id}`)
  } catch (error) {
    alert('Failed to start audit')
  }
}
</script>
