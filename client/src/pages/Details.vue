<template>
  <div class="container">
    <PollCardComponent></PollCardComponent>
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
