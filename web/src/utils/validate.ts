export const validators = {
  required: (message: string = '此字段为必填项') => {
    return { required: true, message, trigger: 'blur' }
  },

  email: (message: string = '请输入正确的邮箱地址') => {
    return {
      type: 'email',
      message,
      trigger: ['blur', 'change']
    }
  },

  url: (message: string = '请输入正确的URL地址') => {
    return {
      type: 'url',
      message,
      trigger: ['blur', 'change']
    }
  },

  minLength: (min: number, message?: string) => {
    return {
      min,
      message: message || `长度不能少于 ${min} 个字符`,
      trigger: ['blur', 'change']
    }
  },

  maxLength: (max: number, message?: string) => {
    return {
      max,
      message: message || `长度不能超过 ${max} 个字符`,
      trigger: ['blur', 'change']
    }
  },

  pattern: (pattern: RegExp, message: string) => {
    return {
      pattern,
      message,
      trigger: ['blur', 'change']
    }
  }
}

export const validateUrl = (url: string): boolean => {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

export const validateJson = (json: string): boolean => {
  try {
    JSON.parse(json)
    return true
  } catch {
    return false
  }
}
