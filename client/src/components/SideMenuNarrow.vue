<template>
  <div
    class="d-flex flex-column flex-shrink-0 px-2 py-3 bg-light vh-100"
    style="min-width: 60px; width: 10%; position: sticky; top: 0">
    <h1 class="display-5 px-0 m-0">
      <router-link class="link" :to="{ name: 'home' }">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="currentColor"
          class="bi bi-tree-fill"
          viewBox="0 0 16 16">
          <path
            d="M8.416.223a.5.5 0 0 0-.832 0l-3 4.5A.5.5 0 0 0 5 5.5h.098L3.076 8.735A.5.5 0 0 0 3.5 9.5h.191l-1.638 3.276a.5.5 0 0 0 .447.724H7V16h2v-2.5h4.5a.5.5 0 0 0 .447-.724L12.31 9.5h.191a.5.5 0 0 0 .424-.765L10.902 5.5H11a.5.5 0 0 0 .416-.777l-3-4.5z" />
        </svg>
      </router-link>
    </h1>
    <hr />
    <ul class="nav nav-pills flex-column mb-auto gap-3">
      <li class="nav-item p-0 m-0 w-100">
        <h4 class="m-0">
          <router-link class="nav-link link p-0 m-0" :to="{ name: 'home' }">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="currentColor"
              class="bi bi-house-fill"
              viewBox="0 0 16 16">
              <path
                fill-rule="evenodd"
                d="M8 3.293l6 6V13.5a1.5 1.5 0 0 1-1.5 1.5h-9A1.5 1.5 0 0 1 2 13.5V9.293l6-6zm5-.793V6l-2-2V2.5a.5.5 0 0 1 .5-.5h1a.5.5 0 0 1 .5.5z" />
              <path
                fill-rule="evenodd"
                d="M7.293 1.5a1 1 0 0 1 1.414 0l6.647 6.646a.5.5 0 0 1-.708.708L8 2.207 1.354 8.854a.5.5 0 1 1-.708-.708L7.293 1.5z" />
            </svg>
          </router-link>
        </h4>
      </li>
      <li v-if="userID" class="nav-item p-0 m-0 w-100">
        <h4 class="m-0">
          <router-link class="nav-link link p-0 m-0" :to="{ name: 'profile' }">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="currentColor"
              class="bi bi-person-fill"
              viewBox="0 0 16 16">
              <path
                d="M3 14s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1H3zm5-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6z" />
            </svg>
          </router-link>
        </h4>
      </li>
      <li class="nav-item p-0 m-0 w-100">
        <a href="#" class="nav-link link-dark px-0">
          <div class="py-1">
            <button
              v-if="userID"
              type="button"
              class="btn btn-lg btn-primary w-100 p-1"
              data-bs-toggle="modal"
              data-bs-target="#newPollModal">
              <em class="bi bi-pencil-fill" />
            </button>
            <router-link
              v-else
              class="btn btn-lg btn-primary"
              :to="{ name: 'signin' }"
              ><em class="bi bi-person-fill"
            /></router-link>
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
            fill="currentColor"
            class="bi bi-person-fill rounded-circle me-2"
            viewBox="0 0 16 16">
            <path
              d="M3 14s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1H3zm5-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6z" />
          </svg>
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
  name: 'SideMenuNarrowComponent',
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
