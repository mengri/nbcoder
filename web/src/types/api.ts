export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PageParams {
  page: number
  size: number
  keyword?: string
}

export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  size: number
  totalPages: number
}

export interface BaseEntity {
  id: string
  createdAt: string
  updatedAt: string
}
