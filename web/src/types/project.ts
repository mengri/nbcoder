import { BaseEntity } from './api'

export enum ProjectStatus {
  ACTIVE = 'ACTIVE',
  ARCHIVED = 'ARCHIVED',
  DELETED = 'DELETED'
}

export interface Project extends BaseEntity {
  name: string
  description: string
  status: ProjectStatus
  gitRepoUrl?: string
  gitBranch?: string
  standards?: string
  branchPolicy?: string
}

export interface CreateProjectDto {
  name: string
  description: string
  gitRepoUrl?: string
  gitBranch?: string
  standards?: string
  branchPolicy?: string
}

export interface UpdateProjectDto extends Partial<CreateProjectDto> {
  status?: ProjectStatus
}

export interface ProjectConfig {
  id: string
  projectId: string
  configKey: string
  configValue: string
  scope: 'GLOBAL' | 'PROJECT'
}
