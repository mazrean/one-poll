import { defineStore } from 'pinia'
import api from '/@/lib/apis'

export const useMainStore = defineStore('main', {
  state: () => {
    return {
      userID: 'default_user'
    }
  },

  getters: {
    getUserID(state) {
      return state.userID
    }
  },

  actions: {
    setUserID() {
      this.userID = api.getUsersMe.name
    }
  }
})
