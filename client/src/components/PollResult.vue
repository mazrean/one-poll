<template>
  <div v-for="(choice, i) in result" :key="choice.id" class="poll-result mb-1">
    <div class="poll-choice position-relative">
      <div
        class="poll-bar position-absolute top-50 start-0 translate-middle-y bg-secondary bg-opacity-25"
        :style="{ width: bg_width[i] + '%' }"></div>
      <div class="position-absolute top-50 start-50 translate-middle">
        {{ choice.choice }}
      </div>
      <div class="position-absolute top-50 start-100 translate-middle-y">
        {{ percentage[i] + '%' }}
      </div>
    </div>
  </div>
  <div class="d-flex justify-content-around">
    <div>{{ count }} 票</div>
    <div class="card-link"><a href="#">詳細を見る</a></div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, watchEffect } from 'vue'
import { Result } from '/@/lib/apis'

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
    const bg_width: number[] = []
    const percentage: number[] = []
    watchEffect(() =>
      props.result.forEach(rslt => {
        if (rslt.count === 0 || props.count === 0) {
          bg_width.push(1)
          percentage.push(0)
        } else {
          bg_width.push((rslt.count * 96) / props.count)
          percentage.push(Math.round((rslt.count * 100) / props.count))
        }
      })
    )
    return {
      bg_width,
      percentage
    }
  }
})
</script>

<style>
.poll-result {
  width: 420px;
}
.poll-choice {
  width: 420px;
  height: 30px;
}
.poll-bar {
  height: 30px;
}
</style>
