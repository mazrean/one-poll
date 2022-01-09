<template>
  <div class="card">
    <div class="card-header text-start">
      <h2 class="card-title">{{ title }}</h2>
      <div class="card-tags text-muted d-flex justify-content-start flex-wrap">
        <div v-for="tag in tags" :key="tag.id" class="me-1 bi bi-tags-fill">
          {{ tag.name }}
        </div>
      </div>
    </div>
    <div class="card-body">
      <div>
        <button
          v-for="q in question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1"
          @click="increment()">
          {{ q.choice }}
        </button>
      </div>
      <div>
        <PollResultComponent />
        <div class="d-flex justify-content-around">
          <div class="card-link"><a href="#">詳細を見る</a></div>
        </div>
      </div>
    </div>
    <div class="footer d-flex justify-content-around">
      <div>投票数 : {{ count }}</div>
      <div>
        作成者 : <a href="#">{{ owner }}</a>
      </div>
      <div>作成日 : {{ createdAt }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import PollResultComponent from '/@/components/PollResult.vue'

export default defineComponent({
  components: { PollResultComponent },
  props: {
    pollId: {
      type: String,
      default: 'poll_id',
      required: false
    },
    title: {
      type: String,
      default: 'Title',
      required: true
    },
    type: {
      type: String,
      default: 'type',
      required: false
    },
    deadline: {
      type: String,
      default: 'deadline',
      required: false
    },
    tags: {
      type: Array,
      default() {
        return [
          {
            id: -1,
            name: 'name'
          }
        ]
      },
      required: true
    },
    question: {
      type: Array,
      default() {
        return [
          {
            id: -1,
            choice: 'choice'
          }
        ]
      },
      required: true
    },
    createdAt: {
      type: String,
      default: 'created_at',
      required: true
    },
    owner: {
      type: String,
      default: 'owner',
      required: true
    }
  },
  setup() {
    const count = ref(0)
    const increment = () => {
      count.value++
    }
    return {
      count,
      increment,
      PollResultComponent
    }
  }
})
</script>

<style>
.card {
  width: 34rem;
}
.vote-button {
  width: 32rem;
}
</style>
