<template>
  <routes />
  <div>Hello</div>
  <HelloWorld msg="Hello Vue 3.0 + Vite" />
  <!--Bootstrap test-->
  <div class="dropdown">
    <button
      id="dropdownMenuButton1"
      class="btn btn-secondary dropdown-toggle"
      type="button"
      data-bs-toggle="dropdown"
      aria-expanded="false"
    >
      Check Bootstrap
    </button>
    <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton1">
      <li>
        <a class="dropdown-item" href="#">Action</a>
      </li>
      <li>
        <a class="dropdown-item" href="#">Another action</a>
      </li>
      <li>
        <a class="dropdown-item" href="#">Something else here</a>
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import Routes from '/@/components/Routes.vue'
import HelloWorld from '/@/components/HelloWorld.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Hello',
  components: {
    Routes,
    HelloWorld
  }
})
</script>
