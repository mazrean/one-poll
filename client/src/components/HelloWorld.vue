<template>
  <h1>{{ msg }}</h1>
  <button @click="increment">count: {{ count }}</button>
  <p>
    Edit <code>components/HelloWorld.vue</code> to test hot module
    <span :class="$style.re">replacement</span>.
  </p>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue'
import { useMainStore } from '/@/store'
import { storeToRefs } from 'pinia'

export default defineComponent({
  name: 'HelloWorld',
  props: {
    msg: {
      type: String,
      required: true
    }
  },
  setup() {
    const main = useMainStore()
    const { countString } = storeToRefs(main)

    return {
      count: computed(() => countString.value),
      increment: main.increment
    }
  }
})
</script>

<style module>
.re {
  font-weight: bold;
}
</style>
