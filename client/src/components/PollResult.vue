<template>
  <div v-for="(choice, i) in result" :key="choice.id" class="poll-choice mb-1">
    <div
      class="poll-bar position-absolute bg-secondary bg-opacity-25"
      :style="{ width: w[i] + 'rem' }"></div>
    <div class="position-relative top-50 start-50 translate-middle">
      {{ choice.choice }}
    </div>
  </div>
  <div class="d-flex justify-content-around">
    <div>{{ count }} 票</div>
    <div class="card-link"><a href="#">詳細を見る</a></div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { Choice } from '../lib/apis'

export default defineComponent({
  props: {
    pollId: {
      type: String,
      required: true
    },
    type: {
      type: String,
      required: true
    },
    count: {
      type: Number,
      required: true
    },
    result: {
      type: Array as PropType<Choice[]>,
      required: true
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
  width: 100%;
  height: 2.7rem;
}
.poll-bar {
  height: 2.7rem;
}
</style>
