import { BaseEntity } from './api'

export enum PipelineStatus {
  PENDING = 'PENDING',
  RUNNING = 'RUNNING',
  PAUSED = 'PAUSED',
  COMPLETED = 'COMPLETED',
  FAILED = 'FAILED',
  CANCELLED = 'CANCELLED'
}

export enum StageStatus {
  PENDING = 'PENDING',
  RUNNING = 'RUNNING',
  REVIEW = 'REVIEW',
  COMPLETED = 'COMPLETED',
  FAILED = 'FAILED',
  SKIPPED = 'SKIPPED'
}

export interface Pipeline extends BaseEntity {
  projectId: string
  cardId: string
  name: string
  status: PipelineStatus
  currentStage?: string
  stages: PipelineStage[]
}

export interface PipelineStage {
  id: string
  name: string
  status: StageStatus
  order: number
  requiresReview: boolean
  steps: PipelineStep[]
  startTime?: string
  endTime?: string
}

export interface PipelineStep {
  id: string
  name: string
  status: StageStatus
  order: number
  agentTaskId?: string
  startTime?: string
  endTime?: string
  output?: any
  error?: string
}

export interface CreatePipelineDto {
  projectId: string
  cardId: string
  name: string
  stages: Omit<PipelineStage, 'id' | 'status' | 'startTime' | 'endTime'>[]
}
