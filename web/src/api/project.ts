import request from '@/utils/request'
import type { Project, CreateProjectDto, UpdateProjectDto, PageParams, PageResult } from '@/types/project'

export const projectApi = {
  list: (params: PageParams) => {
    return request.get<PageResult<Project>>('/projects', { params })
  },

  get: (id: string) => {
    return request.get<Project>(`/projects/${id}`)
  },

  create: (dto: CreateProjectDto) => {
    return request.post<Project>('/projects', dto)
  },

  update: (id: string, dto: UpdateProjectDto) => {
    return request.put<Project>(`/projects/${id}`, dto)
  },

  delete: (id: string) => {
    return request.delete(`/projects/${id}`)
  },

  search: (params: PageParams & { keyword?: string }) => {
    return request.get<PageResult<Project>>('/projects/search', { params })
  }
}

export default projectApi
