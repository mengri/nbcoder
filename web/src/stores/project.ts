import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import request from '@/utils/request'
import type { Project, CreateProjectDto, UpdateProjectDto, PageParams, PageResult } from '@/types/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const loading = ref(false)

  const loadProjects = async () => {
    loading.value = true
    try {
      const data = await request.get<PageResult<Project>>('/projects', {
        params: { page: 1, size: 100 }
      })
      projects.value = data.items
    } catch (error) {
      console.error('Failed to load projects:', error)
    } finally {
      loading.value = false
    }
  }

  const createProject = async (dto: CreateProjectDto) => {
    try {
      const data = await request.post<Project>('/projects', dto)
      projects.value.push(data)
      return data
    } catch (error) {
      console.error('Failed to create project:', error)
      throw error
    }
  }

  const updateProject = async (id: string, dto: UpdateProjectDto) => {
    try {
      const data = await request.put<Project>(`/projects/${id}`, dto)
      const index = projects.value.findIndex(p => p.id === id)
      if (index !== -1) {
        projects.value[index] = data
      }
      if (currentProject.value?.id === id) {
        currentProject.value = data
      }
      return data
    } catch (error) {
      console.error('Failed to update project:', error)
      throw error
    }
  }

  const deleteProject = async (id: string) => {
    try {
      await request.delete(`/projects/${id}`)
      projects.value = projects.value.filter(p => p.id !== id)
      if (currentProject.value?.id === id) {
        currentProject.value = null
      }
    } catch (error) {
      console.error('Failed to delete project:', error)
      throw error
    }
  }

  const setCurrentProject = (project: Project | null) => {
    currentProject.value = project
  }

  const getProject = async (id: string) => {
    try {
      const data = await request.get<Project>(`/projects/${id}`)
      return data
    } catch (error) {
      console.error('Failed to get project:', error)
      throw error
    }
  }

  const searchProjects = async (params: PageParams) => {
    try {
      const data = await request.get<PageResult<Project>>('/projects', { params })
      return data
    } catch (error) {
      console.error('Failed to search projects:', error)
      throw error
    }
  }

  return {
    projects,
    currentProject,
    loading,
    loadProjects,
    createProject,
    updateProject,
    deleteProject,
    setCurrentProject,
    getProject,
    searchProjects
  }
})
