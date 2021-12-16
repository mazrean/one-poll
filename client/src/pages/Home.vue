<template>
  <div class="container">
    <Routes />
    <h1>ホーム</h1>
    <p>ここはホーム画面</p>
    <!-- 以下にスクリプトを記述 -->
  </div>
</template>

<script lang="ts">
import Routes from '../components/Routes.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Home',
  components: { Routes }
})
</script>
