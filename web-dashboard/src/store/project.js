import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useProjectStore = defineStore('project', () => {
  const projects = ref([])
  const currentProject = ref(null)
  const loading = ref(false)

  async function fetchProjects() {
    loading.value = true
    try {
      const response = await axios.get('/api/v1/projects')
      projects.value = response.data.projects || []
    } catch (error) {
      console.error('Failed to fetch projects:', error)
    } finally {
      loading.value = false
    }
  }

  async function fetchProject(id) {
    loading.value = true
    try {
      const response = await axios.get(`/api/v1/projects/${id}`)
      currentProject.value = response.data.project
      return response.data.project
    } catch (error) {
      console.error('Failed to fetch project:', error)
      return null
    } finally {
      loading.value = false
    }
  }

  async function createProject(projectData) {
    try {
      const response = await axios.post('/api/v1/projects', projectData)
      projects.value.push(response.data.project)
      return response.data.project
    } catch (error) {
      console.error('Failed to create project:', error)
      throw error
    }
  }

  async function updateProject(id, projectData) {
    try {
      const response = await axios.put(`/api/v1/projects/${id}`, projectData)
      const index = projects.value.findIndex(p => p.id === id)
      if (index !== -1) {
        projects.value[index] = response.data.project
      }
      return response.data.project
    } catch (error) {
      console.error('Failed to update project:', error)
      throw error
    }
  }

  async function deleteProject(id) {
    try {
      await axios.delete(`/api/v1/projects/${id}`)
      projects.value = projects.value.filter(p => p.id !== id)
    } catch (error) {
      console.error('Failed to delete project:', error)
      throw error
    }
  }

  async function startAudit(projectId) {
    try {
      const response = await axios.post('/api/v1/audits/start', { project_id: projectId })
      return response.data.audit_run
    } catch (error) {
      console.error('Failed to start audit:', error)
      throw error
    }
  }

  return {
    projects,
    currentProject,
    loading,
    fetchProjects,
    fetchProject,
    createProject,
    updateProject,
    deleteProject,
    startAudit
  }
})
