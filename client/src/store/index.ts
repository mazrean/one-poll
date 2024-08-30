import { defineStore } from 'pinia'
import api, { User } from '/@/lib/apis'
import { ref } from 'vue'

export const useMainStore = defineStore('main', () => {
  const userID = ref(localStorage.getItem('sessions'))

  const getUserID = () => {
    return userID.value
  }
  const setUserID = async () => {
    if (!document.cookie.split('; ').some(row => row.startsWith('sessions='))) {
      userID.value = null
      return
    }

    try {
      const localStorageUserID = localStorage.getItem('sessions')
      if (localStorageUserID) {
        userID.value = localStorageUserID
        return
      }

      const user: User = (await api.getUsersMe()).data
      userID.value = user.name
      localStorage.setItem('sessions', user.name)
    } catch {
      userID.value = null
    }
  }
  const setUserIDValue = (value: string) => {
    userID.value = value
    localStorage.setItem('sessions', value)
  }
  return {
    userID,
    getUserID,
    setUserID,
    setUserIDValue
  }
})
