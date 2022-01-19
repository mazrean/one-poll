import { defineStore } from 'pinia'
import api, { User } from '/@/lib/apis'
import { ref } from 'vue'

export const useMainStore = defineStore('main', () => {
  const userID = ref('')
  const getUserID = () => {
    return userID.value
  }
  const setUserID = async () => {
    try {
      const user: User = (await api.getUsersMe()).data
      userID.value = user.name
    } catch (err) {
      userID.value = ''
    }
  }
  return {
    userID,
    getUserID,
    setUserID
  }
})
