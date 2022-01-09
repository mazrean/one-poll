<template>
  <div>
    {{ data }}
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    data: {
      type: Number,
      default: 0,
      required: true
    }
  },
  setup() {
    return {}
  }
})
</script>
