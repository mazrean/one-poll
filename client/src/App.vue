<template>
  <div class="d-flex min-vh-100">
    <SideMenu v-if="width >= 1000" />
    <SideMenuNarrow v-else />
    <main class="p-3 min-vh-100">
      <router-view />
    </main>
  </div>
  <NewPoll />
</template>

<script lang="ts">
import SideMenu from './components/SideMenu.vue'
import SideMenuNarrow from './components/SideMenuNarrow.vue'
import NewPoll from './components/NewPoll.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'App',
  components: { SideMenu, SideMenuNarrow, NewPoll },
  setup() {
    return {}
  },
  data: function () {
    return {
      width: window.innerWidth,
      height: window.innerHeight
    }
  },
  mounted: function () {
    window.addEventListener('resize', this.handleResize)
  },
  beforeUnmount: function () {
    window.removeEventListener('resize', this.handleResize)
  },
  methods: {
    handleResize: function () {
      // resizeのたびにこいつが発火するので、ここでやりたいことをやる
      this.width = window.innerWidth
      this.height = window.innerHeight
    }
  }
})
</script>
