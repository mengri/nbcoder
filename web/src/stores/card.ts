import { defineStore } from 'pinia'
import { ref } from 'vue'
import request from '@/utils/request'
import type { Card, CreateCardDto, UpdateCardDto, PageParams, PageResult } from '@/types/card'

export const useCardStore = defineStore('card', () => {
  const cards = ref<Card[]>([])
  const currentCard = ref<Card | null>(null)
  const loading = ref(false)

  const loadCards = async (projectId: string, params?: PageParams) => {
    loading.value = true
    try {
      const data = await request.get<PageResult<Card>>(`/projects/${projectId}/cards`, {
        params: params || { page: 1, size: 100 }
      })
      cards.value = data.items
      return data
    } catch (error) {
      console.error('Failed to load cards:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const createCard = async (dto: CreateCardDto) => {
    try {
      const data = await request.post<Card>('/cards', dto)
      cards.value.push(data)
      return data
    } catch (error) {
      console.error('Failed to create card:', error)
      throw error
    }
  }

  const updateCard = async (id: string, dto: UpdateCardDto) => {
    try {
      const data = await request.put<Card>(`/cards/${id}`, dto)
      const index = cards.value.findIndex(c => c.id === id)
      if (index !== -1) {
        cards.value[index] = data
      }
      if (currentCard.value?.id === id) {
        currentCard.value = data
      }
      return data
    } catch (error) {
      console.error('Failed to update card:', error)
      throw error
    }
  }

  const deleteCard = async (id: string) => {
    try {
      await request.delete(`/cards/${id}`)
      cards.value = cards.value.filter(c => c.id !== id)
      if (currentCard.value?.id === id) {
        currentCard.value = null
      }
    } catch (error) {
      console.error('Failed to delete card:', error)
      throw error
    }
  }

  const getCard = async (id: string) => {
    try {
      const data = await request.get<Card>(`/cards/${id}`)
      currentCard.value = data
      return data
    } catch (error) {
      console.error('Failed to get card:', error)
      throw error
    }
  }

  const confirmCard = async (id: string) => {
    return updateCard(id, { status: 'CONFIRMED' as any })
  }

  const startCard = async (id: string) => {
    return updateCard(id, { status: 'IN_PROGRESS' as any })
  }

  const completeCard = async (id: string) => {
    return updateCard(id, { status: 'COMPLETED' as any })
  }

  const abandonCard = async (id: string) => {
    return updateCard(id, { status: 'ABANDONED' as any })
  }

  return {
    cards,
    currentCard,
    loading,
    loadCards,
    createCard,
    updateCard,
    deleteCard,
    getCard,
    confirmCard,
    startCard,
    completeCard,
    abandonCard
  }
})
