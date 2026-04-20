import request from '@/utils/request'
import type { Provider, CreateProviderDto, UpdateProviderDto, Model, CreateModelDto, PageParams, PageResult } from '@/types/ai'

export const aiRuntimeApi = {
  listProviders: (params?: PageParams) => {
    return request.get<PageResult<Provider>>('/ai/providers', { params })
  },

  getProvider: (providerId: string) => {
    return request.get<Provider>(`/ai/providers/${providerId}`)
  },

  createProvider: (dto: CreateProviderDto) => {
    return request.post<Provider>('/ai/providers', dto)
  },

  updateProvider: (providerId: string, dto: UpdateProviderDto) => {
    return request.put<Provider>(`/ai/providers/${providerId}`, dto)
  },

  deleteProvider: (providerId: string) => {
    return request.delete(`/ai/providers/${providerId}`)
  },

  testProvider: (providerId: string) => {
    return request.post(`/ai/providers/${providerId}/test`)
  },

  listModels: (providerId?: string, params?: PageParams) => {
    const url = providerId ? `/ai/providers/${providerId}/models` : '/ai/models'
    return request.get<PageResult<Model>>(url, { params })
  },

  getModel: (modelId: string) => {
    return request.get<Model>(`/ai/models/${modelId}`)
  },

  createModel: (dto: CreateModelDto) => {
    return request.post<Model>('/ai/models', dto)
  },

  updateModel: (modelId: string, dto: Partial<CreateModelDto>) => {
    return request.put<Model>(`/ai/models/${modelId}`, dto)
  },

  deleteModel: (modelId: string) => {
    return request.delete(`/ai/models/${modelId}`)
  },

  getStatistics: () => {
    return request.get('/ai/statistics')
  },

  getUsage: (startDate?: string, endDate?: string) => {
    return request.get('/ai/usage', {
      params: { startDate, endDate }
    })
  },

  getCost: (startDate?: string, endDate?: string) => {
    return request.get('/ai/cost', {
      params: { startDate, endDate }
    })
  }
}

export default aiRuntimeApi
