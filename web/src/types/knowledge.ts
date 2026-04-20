import { BaseEntity } from './api'

export enum DocumentStatus {
  PENDING = 'PENDING',
  INDEXING = 'INDEXING',
  INDEXED = 'INDEXED',
  FAILED = 'FAILED'
}

export interface Document extends BaseEntity {
  projectId: string
  name: string
  path: string
  size: number
  status: DocumentStatus
  chunks: number
  indexedAt?: string
}

export interface CreateDocumentDto {
  projectId: string
  name: string
  path: string
}

export interface DocumentChunk {
  id: string
  documentId: string
  content: string
  embedding?: number[]
  metadata?: any
}

export interface IndexStatus {
  id: string
  projectId: string
  totalDocuments: number
  indexedDocuments: number
  totalChunks: number
  indexedChunks: number
  lastIndexedAt?: string
  status: 'IDLE' | 'INDEXING' | 'FAILED'
}

export interface SearchQuery {
  projectId: string
  query: string
  topK?: number
  filter?: any
}

export interface SearchResult {
  chunkId: string
  documentId: string
  documentName: string
  content: string
  score: number
  metadata?: any
}
