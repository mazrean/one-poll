<template>
  <div class="card" v-for="card in cards" v-bind:key="card">
    <div class="card-body">
      <h5 class="card-title">{{ card.title }}</h5>
      <div class="card-header">
        {{ card.question }}
      </div>
      <div v-if="card.beforeVote">
        <button
          type="button"
          class="btn btn-light"
          v-for="choice in card.choices"
          v-bind:key="choice.item">
          {{ choice.item }}
        </button>
      </div>
      <div v-else>
        <SampleChart class="chart" :chart-data="chartData" />
      </div>
    </div>
  </div>
</template>

<style>
  .card { width: 35rem; }
  .card button { width: 33rem; }
</style>

<script lang="ts">
import SampleChart from '../components/SampleChart.vue'
import { defineComponent } from 'vue'

var cards = [
  {
    title: 'Title 1',
    question: 'Question 1',
    beforeVote: false,
    choices: [
      {item: 'Choice 1', count: 0},
      {item: 'Choice 2', count: 0},
      {item: 'Choice 3', count: 0},
    ]
  },
  {
    title: 'Title 2',
    question: 'Question 2',
    beforeVote: true,
    choices: [
      {item: 'Choice 1', count: 0},
      {item: 'Choice 2', count: 0},
      {item: 'Choice 3', count: 0},
    ]
  }
]

export default defineComponent({
  name: 'PollCardComponent',
  data() {
    return {
      cards: cards
    }
  },
  components: { SampleChart },
  setup() {
    return {}
  }
})
</script>
