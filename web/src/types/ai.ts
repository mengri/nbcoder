import { BaseEntity } from './api'

export enum ProviderStatus {
  ACTIVE = 'ACTIVE',
  INACTIVE = 'INACTIVE',
  UNAVAILABLE = 'UNAVAILABLE'
}

export interface AIProvider extends BaseEntity {
  name: string
  type: string
  status: ProviderStatus
  apiKey: string
  baseUrl?: string
  models: AIModel[]
}

export interface AIModel {
  id: string
  providerId: string
  name: string
  type: 'CHAT' | 'COMPLETION' | 'EMBEDDING'
  status: ProviderStatus
  maxTokens: number
  inputPrice: number
  outputPrice: number
}

export interface ModelChain {
  id: string
  name: string
  description: string
  models: ChainModel[]
  fallbackEnabled: boolean
}

export interface ChainModel {
  modelId: string
  order: number
  weight: number
}

export interface ModelCall {
  id: string
  providerId: string
  modelId: string
  modelName: string
  promptTokens: number
  completionTokens: number
  totalTokens: number
  cost: number
  latency: number
  timestamp: string
  success: boolean
  error?: string
}

export interface CreateProviderDto {
  name: string
  type: string
  apiKey: string
  baseUrl?: string
}

export interface Statistics {
  totalCalls: number
  totalTokens: number
  totalCost: number
  todayCalls: number
  modelDistribution: ModelStats[]
  dailyTrend: DailyStats[]
}

export interface ModelStats {
  modelName: string
  calls: number
  tokens: number
  cost: number
}

export interface DailyStats {
  date: string
  calls: number
  tokens: number
  cost: number
}
