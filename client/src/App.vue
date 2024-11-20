<template>
  <div class="d-flex min-vh-100">
    <SideMenu v-if="width >= 768" class="z-1" />
    <SideMenuNarrow v-else class="z-1" />
    <main class="py-2 min-vh-100 w-80">
      <router-view />
    </main>
  </div>
  <NewPoll />
</template>

<script lang="ts">
import SideMenu from '/@/components/SideMenu.vue'
import SideMenuNarrow from '/@/components/SideMenuNarrow.vue'
import NewPoll from '/@/components/NewPoll.vue'
import { defineComponent, ref, onMounted, onBeforeUnmount } from 'vue'

export default defineComponent({
  name: 'App',
  components: { SideMenu, SideMenuNarrow, NewPoll },
  setup() {
    const width = ref(window.innerWidth)
    const handleResize = () => {
      width.value = window.innerWidth
    }
    onMounted(() => {
      window.addEventListener('resize', handleResize)
    })
    onBeforeUnmount(() => {
      window.removeEventListener('resize', handleResize)
    })
    return { width, handleResize }
  }
})
</script>
