import { BaseEntity } from './api'

export enum CardStatus {
  DRAFT = 'DRAFT',
  PENDING = 'PENDING',
  CONFIRMED = 'CONFIRMED',
  IN_PROGRESS = 'IN_PROGRESS',
  COMPLETED = 'COMPLETED',
  SUPPLANTED = 'SUPPLANTED',
  ABANDONED = 'ABANDONED'
}

export enum CardPriority {
  LOW = 'LOW',
  MEDIUM = 'MEDIUM',
  HIGH = 'HIGH',
  CRITICAL = 'CRITICAL'
}

export interface Card extends BaseEntity {
  projectId: string
  title: string
  description: string
  status: CardStatus
  priority: CardPriority
  rawInput: string
  structuredOutput?: any
  dependencies: string[]
  pipelineId?: string
  agentTaskId?: string
}

export interface CreateCardDto {
  projectId: string
  title: string
  description: string
  priority: CardPriority
  rawInput: string
}

export interface UpdateCardDto extends Partial<CreateCardDto> {
  status?: CardStatus
  dependencies?: string[]
  structuredOutput?: any
}

export interface CardDependency {
  cardId: string
  dependsOn: string
}
