<template>
  <div class="card"> <!-- v-for="card in cards" v-bind:key="card" -->
    <div class="card-header text-start">
      <h2 class="card-title">{{ card.title }}</h2>
      <div class="card-tags text-muted d-flex justify-content-start flex-wrap">
        <div class="me-1 bi bi-tags-fill" v-for="tag in card.tags" v-bind:key="tag">{{ tag }}</div>
      </div>
    </div>
    <div class="card-body">
      <p class="card-text">{{ card.question }}</p>
      <div v-if="!user.voted">
        <button
          type="button" class="vote-button btn btn-outline-secondary mb-1"
          v-for="(choice, i) in card.choices" v-bind:key="choice"
          @click="vote(i)">
          {{ choice }}
        </button>
      </div>
      <div v-else>
        <render-chart class="chart" :chart-data="chartData"></render-chart>
        <div class="d-flex justify-content-around">
          <button
            type="button" style="width: 12rem;"
            @click="reset()">
            投票前の状態に戻す
          </button>
          <div class="card-link"><a href="#">詳細を見る</a></div>
        </div>
      </div>

    </div>
    <div class="footer d-flex justify-content-around">
      <div>投票数 : {{ card.count }}</div>
      <div>作成者 : <a href="#">{{ card.author }}</a></div>
      <div>作成日 : {{ card.date }}</div>
    </div>
  </div>
</template>

<style>
.card { width: 36rem; }
.vote-button { width: 34rem; }
</style>

<script lang="ts">
import RenderChart from '/@/components/RenderChart.vue'
import { defineComponent, reactive } from 'vue'

export default defineComponent({
  components: { RenderChart },
  setup() {
    const card = reactive({
      title: 'Title',
      tags: ['tag1', 'tag2', 'tag3'],
      question: 'Question',
      choices: ['choice1', 'choice2', 'choice3', 'choice4'],
      count: 140,
      author: 'moka',
      date: '2022/1/4 14:07',
    });
    const user = reactive({
      voted: false
    });
    const data = reactive([0, 0, 0, 0])
    const vote = (i: number) => {
      data[i]++;
      user.voted = true;
      // test
      console.log(data[0]);
      console.log(data[1]);
      console.log(data[2]);
      console.log(data[3]);
    };
    const reset = () => {
      user.voted = false;
    };
    return {
      card, user, vote, reset
    };
  }
})
</script>
