import { ref, computed } from 'vue'
import { useCardStore } from '@/stores/card'
import type { Card, CreateCardDto, UpdateCardDto } from '@/types/card'

export function useCard() {
  const cardStore = useCardStore()
  const loading = ref(false)

  const currentCard = computed(() => cardStore.currentCard)
  const cards = computed(() => cardStore.cards)

  const loadCards = async (projectId: string, params?: { page?: number; size?: number; keyword?: string }) => {
    loading.value = true
    try {
      return await cardStore.loadCards(projectId, params || {})
    } finally {
      loading.value = false
    }
  }

  const createCard = async (projectId: string, dto: CreateCardDto) => {
    loading.value = true
    try {
      return await cardStore.createCard(projectId, dto)
    } finally {
      loading.value = false
    }
  }

  const updateCard = async (projectId: string, cardId: string, dto: UpdateCardDto) => {
    loading.value = true
    try {
      return await cardStore.updateCard(projectId, cardId, dto)
    } finally {
      loading.value = false
    }
  }

  const deleteCard = async (projectId: string, cardId: string) => {
    loading.value = true
    try {
      await cardStore.deleteCard(projectId, cardId)
    } finally {
      loading.value = false
    }
  }

  const getCard = async (projectId: string, cardId: string) => {
    loading.value = true
    try {
      return await cardStore.getCard(projectId, cardId)
    } finally {
      loading.value = false
    }
  }

  const confirmCard = async (projectId: string, cardId: string) => {
    loading.value = true
    try {
      await cardStore.confirmCard(projectId, cardId)
    } finally {
      loading.value = false
    }
  }

  const startCard = async (projectId: string, cardId: string) => {
    loading.value = true
    try {
      await cardStore.startCard(projectId, cardId)
    } finally {
      loading.value = false
    }
  }

  const completeCard = async (projectId: string, cardId: string) => {
    loading.value = true
    try {
      await cardStore.completeCard(projectId, cardId)
    } finally {
      loading.value = false
    }
  }

  const abandonCard = async (projectId: string, cardId: string) => {
    loading.value = true
    try {
      await cardStore.abandonCard(projectId, cardId)
    } finally {
      loading.value = false
    }
  }

  const addDependency = async (projectId: string, cardId: string, dependsOnId: string) => {
    loading.value = true
    try {
      await cardStore.addDependency(projectId, cardId, dependsOnId)
    } finally {
      loading.value = false
    }
  }

  const removeDependency = async (projectId: string, cardId: string, dependsOnId: string) => {
    loading.value = true
    try {
      await cardStore.removeDependency(projectId, cardId, dependsOnId)
    } finally {
      loading.value = false
    }
  }

  const setCurrentCard = (card: Card | null) => {
    cardStore.currentCard = card
  }

  return {
    loading,
    currentCard,
    cards,
    loadCards,
    createCard,
    updateCard,
    deleteCard,
    getCard,
    confirmCard,
    startCard,
    completeCard,
    abandonCard,
    addDependency,
    removeDependency,
    setCurrentCard
  }
}
