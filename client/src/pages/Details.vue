<template>
  <div class="container">
    <Routes />
    <h1>ディテール</h1>
    <p>ここは詳細画面</p>
    <!-- 以下にスクリプトを記述 -->
    <poll-card-component></poll-card-component>
  </div>
</template>

<script lang="ts">
import PollCardComponent from '/@/components/PollCard.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'DetailsPage',
  components: { PollCardComponent }
})
</script>
