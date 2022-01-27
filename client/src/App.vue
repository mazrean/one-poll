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
import SideMenu from '/@/components/SideMenu.vue'
import SideMenuNarrow from '/@/components/SideMenuNarrow.vue'
import NewPoll from '/@/components/NewPoll.vue'
import { defineComponent, ref, onMounted, onBeforeUnmount } from 'vue'
import { useMainStore } from '/@/store/index'

export default defineComponent({
  name: 'App',
  components: { SideMenu, SideMenuNarrow, NewPoll },
  setup() {
    const store = useMainStore()
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
