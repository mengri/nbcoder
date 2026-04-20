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
  repoUrl?: string
}

export interface CreateProjectDto {
  name: string
  description: string
  repoUrl?: string
  branchStrategy?: string
  techStack?: string
  codingConventions?: string
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
