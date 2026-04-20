import { ref, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'

export function useWebSocket(url: string) {
  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)
  const messageHandler = ref<((data: any) => void) | null>(null)
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5
  const reconnectDelay = 3000

  const connect = () => {
    try {
      ws.value = new WebSocket(url)

      ws.value.onopen = () => {
        connected.value = true
        reconnectAttempts.value = 0
        console.log('WebSocket connected')
      }

      ws.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          if (messageHandler.value) {
            messageHandler.value(data)
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      ws.value.onerror = (error) => {
        console.error('WebSocket error:', error)
      }

      ws.value.onclose = () => {
        connected.value = false
        console.log('WebSocket disconnected')

        if (reconnectAttempts.value < maxReconnectAttempts) {
          reconnectAttempts.value++
          setTimeout(() => {
            console.log(`Reconnecting... Attempt ${reconnectAttempts.value}`)
            connect()
          }, reconnectDelay)
        } else {
          ElMessage.error('WebSocket 连接失败，请刷新页面重试')
        }
      }
    } catch (error) {
      console.error('Failed to connect WebSocket:', error)
      ElMessage.error('WebSocket 连接失败')
    }
  }

  const send = (data: any) => {
    if (ws.value && connected.value) {
      ws.value.send(JSON.stringify(data))
    } else {
      console.warn('WebSocket is not connected')
    }
  }

  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    connected.value = false
  }

  const onMessage = (handler: (data: any) => void) => {
    messageHandler.value = handler
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    connect,
    send,
    disconnect,
    onMessage
  }
}
