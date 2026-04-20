import request from '@/utils/request'
import type { Document, CreateDocumentDto, UpdateDocumentDto, IndexStatus, PageParams, PageResult } from '@/types/knowledge'

export const knowledgeApi = {
  listDocuments: (projectId: string, params: PageParams) => {
    return request.get<PageResult<Document>>(`/projects/${projectId}/documents`, { params })
  },

  getDocument: (projectId: string, documentId: string) => {
    return request.get<Document>(`/projects/${projectId}/documents/${documentId}`)
  },

  createDocument: (projectId: string, dto: CreateDocumentDto) => {
    return request.post<Document>(`/projects/${projectId}/documents`, dto)
  },

  updateDocument: (projectId: string, documentId: string, dto: UpdateDocumentDto) => {
    return request.put<Document>(`/projects/${projectId}/documents/${documentId}`, dto)
  },

  deleteDocument: (projectId: string, documentId: string) => {
    return request.delete(`/projects/${projectId}/documents/${documentId}`)
  },

  uploadDocument: (projectId: string, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return request.post<Document>(`/projects/${projectId}/documents/upload`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  downloadDocument: (projectId: string, documentId: string) => {
    return request.get(`/projects/${projectId}/documents/${documentId}/download`, {
      responseType: 'blob'
    })
  },

  getIndexStatus: (projectId: string) => {
    return request.get<IndexStatus>(`/projects/${projectId}/index-status`)
  },

  reindexDocument: (projectId: string, documentId: string) => {
    return request.post(`/projects/${projectId}/documents/${documentId}/reindex`)
  },

  reindexAll: (projectId: string) => {
    return request.post(`/projects/${projectId}/index-status/reindex-all`)
  },

  optimizeIndex: (projectId: string) => {
    return request.post(`/projects/${projectId}/index-status/optimize`)
  },

  searchDocuments: (projectId: string, query: string, params?: PageParams) => {
    return request.get<PageResult<Document>>(`/projects/${projectId}/documents/search`, {
      params: { query, ...params }
    })
  }
}

export default knowledgeApi
