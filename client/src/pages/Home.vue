<template>
  <div class="container">
    <h1><em class="bi bi-house-fill" /> ホーム</h1>
    <div class="m-auto">
      <div class="d-flex flex-wrap justify-content-center m-auto">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          fill="currentColor"
          class="bi bi-search mx-3 my-auto"
          viewBox="0 0 16 16">
          <path
            d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z" />
        </svg>
        <input
          v-model="state.searchTitle"
          type="searchTitle"
          class="form-control d-flex mx-1 my-1 w-25"
          name="searchTitle"
          placeholder="キーワードで検索"
          @Input=";(state.searchTitle = $event.target.value), getPolls()" />
        <input
          v-model="state.searchTag"
          type="searchTag"
          class="form-control d-flex mx-1 my-1 w-25"
          name="searchTag"
          placeholder="タグで検索"
          @Input=";(state.searchTag = $event.target.value), getPolls()" />
      </div>
      <div
        v-if="state.isLoading"
        class="spinner-border text-secondary"
        role="status"></div>
      <div v-else-if="state.PollSummaries.length === 0">
        <p>表示可能な質問がありません。</p>
      </div>
      <div v-else class="d-flex flex-wrap justify-content-center">
        <div
          v-for="PollSummary in state.PollSummaries"
          :key="PollSummary.pollId">
          <PollCardComponent
            :poll-id="PollSummary.pollId"
            :title="PollSummary.title"
            :type="PollSummary.type"
            :deadline="PollSummary.deadline"
            :tags="
              typeof PollSummary.tags !== 'undefined'
                ? PollSummary.tags
                : [{ id: '-1', name: '' }]
            "
            :question="PollSummary.question"
            :created-at="PollSummary.createdAt"
            :q-status="PollSummary.qStatus"
            :owner="PollSummary.owner"
            :user-status="PollSummary.userStatus"
            class="m-4">
          </PollCardComponent>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import PollCardComponent from '/@/components/PollCard.vue'
import apis, { PollSummary } from '/@/lib/apis'
import { reactive, onMounted } from 'vue'

interface State {
  PollSummaries: PollSummary[]
  isLoading: boolean
  searchLimit: number
  searchOffset: number
  searchTitle: string
  searchTag: string
}

export default {
  name: 'HomePage',
  components: { PollCardComponent },
  setup() {
    const state = reactive<State>({
      PollSummaries: [],
      isLoading: true,
      searchLimit: 100,
      searchOffset: 0,
      searchTitle: '',
      searchTag: ''
    })
    onMounted(async () => {
      await getPolls()
    })
    const getPolls = async () => {
      //state.isLoading = true
      try {
        state.PollSummaries = (
          await apis.getPolls(
            state.searchLimit,
            state.searchOffset,
            state.searchTitle
          )
        ).data
      } catch {
        state.PollSummaries = []
      }
      if (state.searchTag !== '') {
        state.PollSummaries = state.PollSummaries.filter(v => {
          return typeof v.tags !== 'undefined'
            ? v.tags.some(e => e.name === state.searchTag)
            : false
        })
      }
      state.isLoading = false
    }
    return {
      state,
      getPolls,
      PollCardComponent
    }
  }
}
</script>
