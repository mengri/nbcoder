import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const userName = ref<string>('Admin')
  const userId = ref<string>('1')
  const permissions = ref<string[]>(['*'])

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const clearToken = () => {
    token.value = null
    localStorage.removeItem('token')
  }

  const setUserInfo = (name: string, id: string) => {
    userName.value = name
    userId.value = id
  }

  const hasPermission = (permission: string) => {
    if (permissions.value.includes('*')) {
      return true
    }
    return permissions.value.includes(permission)
  }

  const isAuthenticated = computed(() => !!token.value)

  return {
    token,
    userName,
    userId,
    permissions,
    setToken,
    clearToken,
    setUserInfo,
    hasPermission,
    isAuthenticated
  }
})
