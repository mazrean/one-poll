<template>
  <div
    class="d-flex flex-column flex-shrink-0 p-3 bg-light vh-100"
    style="width: 280px; position: sticky; top: 0">
    <h1 class="display-5 py-3">
      <router-link
        class="link"
        :to="{ name: 'home' }"
        style="font-weight: normal"
        >One Poll</router-link
      >
    </h1>
    <hr />
    <ul class="nav nav-pills flex-column mb-auto">
      <li class="nav-item">
        <h4>
          <router-link class="nav-link link" :to="{ name: 'home' }">
            <em class="bi bi-house-fill" /> ホーム</router-link
          >
        </h4>
      </li>
      <li v-if="userID">
        <h4>
          <router-link class="nav-link link" :to="{ name: 'profile' }">
            <em class="bi bi-person-fill" /> プロフィール
          </router-link>
        </h4>
      </li>
      <li>
        <a href="#" class="nav-link link-dark">
          <div class="py-1">
            <button
              v-if="userID"
              type="button"
              class="btn btn-lg btn-primary"
              data-bs-toggle="modal"
              data-bs-target="#newPollModal">
              <strong>新しい質問を作る</strong>
            </button>
            <router-link
              v-else
              class="btn btn-lg btn-primary"
              :to="{ name: 'signin' }">
              <strong>サインインする</strong>
            </router-link>
          </div>
        </a>
      </li>
    </ul>
    <div v-if="userID">
      <hr />
      <div class="dropdown">
        <a
          id="dropdownUser2"
          href="#"
          class="d-flex align-items-center link dropdown-toggle"
          data-bs-toggle="dropdown"
          aria-expanded="false">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="32"
            height="32"
            fill="currentColor"
            class="bi bi-person-fill rounded-circle me-2"
            viewBox="0 0 16 16">
            <path
              d="M3 14s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1H3zm5-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6z" />
          </svg>
          <strong>@{{ userID }}</strong>
        </a>
        <ul
          class="dropdown-menu text-small shadow"
          aria-labelledby="dropdownUser2">
          <li>
            <a class="dropdown-item" @click="onSignOut()">サインアウト</a>
          </li>
          <li><hr class="dropdown-divider" /></li>
          <li><a class="dropdown-item" @click="onDeleteUser()">退会</a></li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue'
import { useMainStore } from '/@/store/index'
import api from '/@/lib/apis'

export default defineComponent({
  name: 'SideMenuComponent',
  components: {},
  setup() {
    const store = useMainStore()
    const userID = computed(() => store.userID)
    const onSignOut = async () => {
      await api.postUsersSignout()
      store.setUserID()
      location.href = '/'
    }
    const onDeleteUser = async () => {
      await api.deleteUsersMe()
      await store.setUserID()
      location.href = '/'
    }
    return {
      userID,
      onSignOut,
      onDeleteUser
    }
  }
})
</script>
