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
        <div class="bg-white rounded-lg shadow p-6">
          <h1 class="text-2xl font-bold text-gray-800 mb-4">Audit Results</h1>
          <p class="text-gray-600">Audit ID: {{ $route.params.id }}</p>
          <div class="mt-6">
            <p class="text-gray-500">Loading audit results...</p>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
// Audit detail view - to be implemented
</script>
