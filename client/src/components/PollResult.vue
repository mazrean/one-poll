<template>
  <div v-for="(choice, i) in result" :key="choice.id" class="poll-choice mb-1">
    <div
      class="poll-bar position-absolute bg-secondary bg-opacity-25"
      :style="{ width: w[i] + 'rem' }"></div>
    <div class="position-relative top-50 start-50 translate-middle">
      {{ choice.choice }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'

interface Choice {
  id: number
  choice: string
  count: number
}

/*
interface Result {
  choice: Array<Choice>
}
*/

export default defineComponent({
  props: {
    result: {
      type: Array as PropType<Choice[]>,
      default() {
        return [
          {
            id: -1,
            choice: 'choice',
            count: 0
          }
        ]
      },
      required: false
    }
  },
  setup() {
    const w = [4, 8, 12, 6]
    return {
      w
    }
  }
})
</script>

<style>
.poll-choice {
  width: 30rem;
  height: 2.7rem;
}
.poll-bar {
  height: 2.7rem;
}
</style>
