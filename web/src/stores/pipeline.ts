import { defineStore } from 'pinia'
import { ref } from 'vue'
import request from '@/utils/request'
import type { Pipeline, CreatePipelineDto, PageParams, PageResult } from '@/types/pipeline'

export const usePipelineStore = defineStore('pipeline', () => {
  const pipelines = ref<Pipeline[]>([])
  const currentPipeline = ref<Pipeline | null>(null)
  const loading = ref(false)

  const loadPipelines = async (projectId: string, params?: PageParams) => {
    loading.value = true
    try {
      const data = await request.get<PageResult<Pipeline>>(`/projects/${projectId}/pipelines`, {
        params: params || { page: 1, size: 100 }
      })
      pipelines.value = data.items
      return data
    } catch (error) {
      console.error('Failed to load pipelines:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const createPipeline = async (dto: CreatePipelineDto) => {
    try {
      const data = await request.post<Pipeline>('/pipelines', dto)
      pipelines.value.push(data)
      return data
    } catch (error) {
      console.error('Failed to create pipeline:', error)
      throw error
    }
  }

  const getPipeline = async (id: string) => {
    try {
      const data = await request.get<Pipeline>(`/pipelines/${id}`)
      currentPipeline.value = data
      return data
    } catch (error) {
      console.error('Failed to get pipeline:', error)
      throw error
    }
  }

  const pausePipeline = async (id: string) => {
    try {
      const data = await request.post<Pipeline>(`/pipelines/${id}/pause`)
      return data
    } catch (error) {
      console.error('Failed to pause pipeline:', error)
      throw error
    }
  }

  const resumePipeline = async (id: string) => {
    try {
      const data = await request.post<Pipeline>(`/pipelines/${id}/resume`)
      return data
    } catch (error) {
      console.error('Failed to resume pipeline:', error)
      throw error
    }
  }

  const cancelPipeline = async (id: string) => {
    try {
      const data = await request.post<Pipeline>(`/pipelines/${id}/cancel`)
      return data
    } catch (error) {
      console.error('Failed to cancel pipeline:', error)
      throw error
    }
  }

  const getPipelineLogs = async (id: string) => {
    try {
      const data = await request.get<any[]>(`/pipelines/${id}/logs`)
      return data
    } catch (error) {
      console.error('Failed to get pipeline logs:', error)
      throw error
    }
  }

  return {
    pipelines,
    currentPipeline,
    loading,
    loadPipelines,
    createPipeline,
    getPipeline,
    pausePipeline,
    resumePipeline,
    cancelPipeline,
    getPipelineLogs
  }
})
