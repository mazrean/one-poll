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
      <li>
        <h4>
          <router-link class="nav-link link" :to="{ name: 'profile' }">
            <em class="bi bi-person-fill" /> プロフィール</router-link
          >
        </h4>
      </li>
      <li>
        <a href="#" class="nav-link link-dark">
          <div class="py-1">
            <button
              type="button"
              class="btn btn-lg btn-primary"
              data-bs-toggle="modal"
              data-bs-target="#newPollModal">
              新しい質問を作る
            </button>
          </div>
        </a>
      </li>
    </ul>
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
        <strong>@aya_se</strong>
      </a>
      <ul
        class="dropdown-menu text-small shadow"
        aria-labelledby="dropdownUser2">
        <li><a class="dropdown-item" href="#">プロフィール</a></li>
        <li><hr class="dropdown-divider" /></li>
        <li><a class="dropdown-item" href="#">サインアウト</a></li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'SideMenuComponent',
  isNarrow: false,
  components: {},
  setup() {
    return {}
  }
})
</script>
