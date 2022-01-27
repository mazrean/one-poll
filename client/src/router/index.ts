import { createRouter, createWebHistory } from 'vue-router'
import Home from '/@/pages/Home.vue'
import Signup from '/@/pages/Signup.vue'
import Signin from '/@/pages/Signin.vue'
import Profile from '/@/pages/Profile.vue'
import Details from '/@/pages/Details.vue'

export const routerHistory = createWebHistory()

const router = createRouter({
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
      component: Profile,
      meta: { requiresAuth: true }
    },
    {
      path: '/details/:pollId',
      name: 'details',
      component: Details,
      meta: { requiresAuth: true }
    }
  ]
})

import { useMainStore } from '/@/store/index'

router.beforeEach(async (to, _, next) => {
  const store = useMainStore()
  await store.setUserID()
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // requiresAuthがtrueなら評価
    if (store.getUserID() === '') {
      // 未ログインならログインページへ
      next('/signin')
    } else {
      next() // スルー
    }
  } else {
    next() // スルー
  }
})

export default router
