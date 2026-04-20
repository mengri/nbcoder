import { BaseEntity } from './api'

export enum AgentTaskStatus {
  PENDING = 'PENDING',
  ASSIGNED = 'ASSIGNED',
  RUNNING = 'RUNNING',
  COMPLETED = 'COMPLETED',
  FAILED = 'FAILED',
  CANCELLED = 'CANCELLED'
}

export enum AgentType {
  CODE_GENERATION = 'CODE_GENERATION',
  TEST = 'TEST',
  DOCUMENTATION = 'DOCUMENTATION',
  REVIEW = 'REVIEW',
  REFACTORING = 'REFACTORING',
  CUSTOM = 'CUSTOM'
}

export interface AgentTask extends BaseEntity {
  projectId: string
  type: AgentType
  status: AgentTaskStatus
  cardId?: string
  pipelineId?: string
  stageId?: string
  stepId?: string
  cloneInstanceId?: string
  input: any
  output?: any
  error?: string
  skill?: string
  startedAt?: string
  completedAt?: string
  duration?: number
}

export interface AgentTaskLog {
  id: string
  taskId: string
  level: 'INFO' | 'WARN' | 'ERROR'
  message: string
  timestamp: string
  data?: any
}
