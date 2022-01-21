<template>
  <div v-for="(choice, i) in result" :key="choice.id" class="poll-choice mb-1">
    <div
      class="poll-bar position-absolute bg-secondary bg-opacity-25"
      :style="{ width: arr[i] + '%' }"></div>
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
import { defineComponent, PropType, watchEffect } from 'vue'
import { Result } from '../lib/apis'

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
      type: Array as PropType<Result[]>,
      required: true
    }
  },
  setup(props) {
    const arr: number[] = []
    watchEffect(() =>
      props.result.forEach(el => {
        if (el.count == 0 || props.count == 0) {
          arr.push(1)
        } else {
          arr.push((el.count * 90) / props.count)
        }
      })
    )
    return {
      arr
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
