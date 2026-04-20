import request from '@/utils/request'
import type { AgentTask, CreateAgentTaskDto, UpdateAgentTaskDto, PageParams, PageResult } from '@/types/agent'

export const agentApi = {
  listTasks: (projectId: string, params: PageParams) => {
    return request.get<PageResult<AgentTask>>(`/projects/${projectId}/agent-tasks`, { params })
  },

  getTask: (projectId: string, taskId: string) => {
    return request.get<AgentTask>(`/projects/${projectId}/agent-tasks/${taskId}`)
  },

  createTask: (projectId: string, dto: CreateAgentTaskDto) => {
    return request.post<AgentTask>(`/projects/${projectId}/agent-tasks`, dto)
  },

  updateTask: (projectId: string, taskId: string, dto: UpdateAgentTaskDto) => {
    return request.put<AgentTask>(`/projects/${projectId}/agent-tasks/${taskId}`, dto)
  },

  deleteTask: (projectId: string, taskId: string) => {
    return request.delete(`/projects/${projectId}/agent-tasks/${taskId}`)
  },

  retryTask: (projectId: string, taskId: string) => {
    return request.post(`/projects/${projectId}/agent-tasks/${taskId}/retry`)
  },

  cancelTask: (projectId: string, taskId: string) => {
    return request.post(`/projects/${projectId}/agent-tasks/${taskId}/cancel`)
  },

  getTaskLogs: (projectId: string, taskId: string) => {
    return request.get(`/projects/${projectId}/agent-tasks/${taskId}/logs`)
  },

  getTaskStatistics: (projectId: string) => {
    return request.get(`/projects/${projectId}/agent-tasks/statistics`)
  }
}

export default agentApi
