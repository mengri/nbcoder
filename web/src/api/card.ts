import request from '@/utils/request'
import type { Card, CreateCardDto, UpdateCardDto, PageParams, PageResult } from '@/types/card'

export const cardApi = {
  list: (projectId: string, params: PageParams) => {
    return request.get<PageResult<Card>>(`/projects/${projectId}/cards`, { params })
  },

  get: (projectId: string, cardId: string) => {
    return request.get<Card>(`/projects/${projectId}/cards/${cardId}`)
  },

  create: (projectId: string, dto: CreateCardDto) => {
    return request.post<Card>(`/projects/${projectId}/cards`, dto)
  },

  update: (projectId: string, cardId: string, dto: UpdateCardDto) => {
    return request.put<Card>(`/projects/${projectId}/cards/${cardId}`, dto)
  },

  delete: (projectId: string, cardId: string) => {
    return request.delete(`/projects/${projectId}/cards/${cardId}`)
  },

  confirm: (projectId: string, cardId: string) => {
    return request.post(`/projects/${projectId}/cards/${cardId}/confirm`)
  },

  start: (projectId: string, cardId: string) => {
    return request.post(`/projects/${projectId}/cards/${cardId}/start`)
  },

  complete: (projectId: string, cardId: string) => {
    return request.post(`/projects/${projectId}/cards/${cardId}/complete`)
  },

  abandon: (projectId: string, cardId: string) => {
    return request.post(`/projects/${projectId}/cards/${cardId}/abandon`)
  },

  addDependency: (projectId: string, cardId: string, dependsOnId: string) => {
    return request.post(`/projects/${projectId}/cards/${cardId}/dependencies`, { dependsOnId })
  },

  removeDependency: (projectId: string, cardId: string, dependsOnId: string) => {
    return request.delete(`/projects/${projectId}/cards/${cardId}/dependencies/${dependsOnId}`)
  }
}

export default cardApi
