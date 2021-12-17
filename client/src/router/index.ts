import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import Signup from '../pages/Signup.vue'
import Signin from '../pages/Signin.vue'
import Profile from '../pages/Profile.vue'

export const routerHistory = createWebHistory()

export default createRouter({
  history: routerHistory,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/signup',
      name: 'signup',
      component: Signup
    },
    {
      path: '/signin',
      name: 'signin',
      component: Signin
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile
    }
  ]
})
