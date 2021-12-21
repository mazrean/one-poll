<template>
  <div
    class="d-flex flex-column flex-shrink-0 p-3 bg-light vh-100"
    style="width: 90px; position: sticky; top: 0">
    <h1 class="display-5 py-3 px-0">
      <router-link
        class="link"
        :to="{ name: 'home' }"
        style="font-weight: normal"
        ><svg
          xmlns="http://www.w3.org/2000/svg"
          width="32"
          height="32"
          fill="currentColor"
          class="bi bi-tree-fill"
          viewBox="0 0 16 16">
          <path
            d="M8.416.223a.5.5 0 0 0-.832 0l-3 4.5A.5.5 0 0 0 5 5.5h.098L3.076 8.735A.5.5 0 0 0 3.5 9.5h.191l-1.638 3.276a.5.5 0 0 0 .447.724H7V16h2v-2.5h4.5a.5.5 0 0 0 .447-.724L12.31 9.5h.191a.5.5 0 0 0 .424-.765L10.902 5.5H11a.5.5 0 0 0 .416-.777l-3-4.5z" /></svg
      ></router-link>
    </h1>
    <hr />
    <ul class="nav nav-pills flex-column mb-auto">
      <li class="nav-item">
        <h4>
          <router-link class="nav-link link px-0" :to="{ name: 'home' }">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              fill="currentColor"
              class="bi bi-house-fill"
              viewBox="0 0 16 16">
              <path
                fill-rule="evenodd"
                d="M8 3.293l6 6V13.5a1.5 1.5 0 0 1-1.5 1.5h-9A1.5 1.5 0 0 1 2 13.5V9.293l6-6zm5-.793V6l-2-2V2.5a.5.5 0 0 1 .5-.5h1a.5.5 0 0 1 .5.5z" />
              <path
                fill-rule="evenodd"
                d="M7.293 1.5a1 1 0 0 1 1.414 0l6.647 6.646a.5.5 0 0 1-.708.708L8 2.207 1.354 8.854a.5.5 0 1 1-.708-.708L7.293 1.5z" /></svg
          ></router-link>
        </h4>
      </li>
      <li>
        <h4>
          <router-link class="nav-link link px-0" :to="{ name: 'profile' }">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="32"
              height="32"
              fill="currentColor"
              class="bi bi-person-fill"
              viewBox="0 0 16 16">
              <path
                d="M3 14s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1H3zm5-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6z" /></svg
          ></router-link>
        </h4>
      </li>
      <li>
        <a href="#" class="nav-link link-dark px-0">
          <div class="py-1">
            <button
              type="button"
              class="btn btn-lg btn-primary"
              data-bs-toggle="modal"
              data-bs-target="#newPollModal">
              <em class="bi bi-pencil-fill" />
            </button>
          </div>
        </a>
      </li>
    </ul>
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
