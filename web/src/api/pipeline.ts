import request from '@/utils/request'
import type { Pipeline, CreatePipelineDto, UpdatePipelineDto, PageParams, PageResult } from '@/types/pipeline'

export const pipelineApi = {
  list: (projectId: string, params: PageParams) => {
    return request.get<PageResult<Pipeline>>(`/projects/${projectId}/pipelines`, { params })
  },

  get: (projectId: string, pipelineId: string) => {
    return request.get<Pipeline>(`/projects/${projectId}/pipelines/${pipelineId}`)
  },

  create: (projectId: string, dto: CreatePipelineDto) => {
    return request.post<Pipeline>(`/projects/${projectId}/pipelines`, dto)
  },

  update: (projectId: string, pipelineId: string, dto: UpdatePipelineDto) => {
    return request.put<Pipeline>(`/projects/${projectId}/pipelines/${pipelineId}`, dto)
  },

  delete: (projectId: string, pipelineId: string) => {
    return request.delete(`/projects/${projectId}/pipelines/${pipelineId}`)
  },

  start: (projectId: string, pipelineId: string) => {
    return request.post(`/projects/${projectId}/pipelines/${pipelineId}/start`)
  },

  pause: (projectId: string, pipelineId: string) => {
    return request.post(`/projects/${projectId}/pipelines/${pipelineId}/pause`)
  },

  resume: (projectId: string, pipelineId: string) => {
    return request.post(`/projects/${projectId}/pipelines/${pipelineId}/resume`)
  },

  cancel: (projectId: string, pipelineId: string) => {
    return request.post(`/projects/${projectId}/pipelines/${pipelineId}/cancel`)
  },

  getStages: (projectId: string, pipelineId: string) => {
    return request.get(`/projects/${projectId}/pipelines/${pipelineId}/stages`)
  }
}

export default pipelineApi
