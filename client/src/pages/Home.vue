<template>
  <div class="container">
    <h1><em class="bi bi-house-fill" /> ホーム</h1>
    <div class="m-auto">
      <div class="d-flex flex-wrap justify-content-center m-auto">
        <svg
          v-show="false"
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
          class="form-control d-flex mx-3 my-1"
          name="searchTitle"
          placeholder="キーワードで検索 (Enterで更新)"
          @Input="state.searchTitle = $event.target.value"
          @keydown.enter="onKeyword()" />
        <input
          v-model="state.searchTag"
          type="searchTag"
          class="form-control d-flex mx-3 my-1"
          name="searchTag"
          placeholder="タグで検索 (タグ候補選択で更新)"
          @Input="calculateFilter()" />
        <ul v-for="v in state.autocompletes" :key="v" class="list-group">
          <button
            class="list-group-item list-group-item-action p-1"
            type="button"
            @click="onAutocomplete(v.name)">
            <em class="bi bi-tags-fill" /> {{ v.name }}
          </button>
        </ul>
      </div>
      <div
        v-if="state.isLoading"
        class="spinner-border text-secondary m-3"
        role="status"></div>
      <div v-else-if="state.PollSummaries.length === 0" class="m-3">
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
            :deadline="
              typeof PollSummary.deadline !== 'undefined'
                ? PollSummary.deadline
                : '-1'
            "
            :tags="
              typeof PollSummary.tags !== 'undefined'
                ? PollSummary.tags
                : [{ id: '-1', name: '' }]
            "
            :question="PollSummary.question"
            :created-at="PollSummary.createdAt"
            :q-status="PollSummary.qStatus"
            :owner="PollSummary.owner"
            :user-status="PollSummary.userStatus">
          </PollCardComponent>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import PollCardComponent from '/@/components/PollCard.vue'
import apis, { PollSummary, PollTag } from '/@/lib/apis'
import { reactive, onMounted, onUnmounted } from 'vue'

interface State {
  PollSummaries: PollSummary[]
  PollSummaries_origin: PollSummary[]
  searchTitle: string
  searchTag: string
  tags: PollTag[]
  autocompletes: PollTag[]
}

export default {
  name: 'HomePage',
  components: { PollCardComponent },
  setup() {
    const state = reactive<State>({
      PollSummaries: [],
      PollSummaries_origin: [],
      searchTitle: '',
      searchTag: '',
      tags: [],
      autocompletes: []
    })

    let limit = 5
    let offset = 0
    let isLoading = false
    let isEnd = false

    const getPolls = async () => {
      if (isLoading || isEnd) return []

      isLoading = true
      let newPollSummaries: PollSummary[]
      try {
        newPollSummaries = (
          await apis.getPolls(limit, offset, state.searchTitle)
        ).data
      } catch {
        newPollSummaries = []
      }

      if (newPollSummaries.length < limit) {
        isEnd = true
      }
      offset += newPollSummaries.length

      state.PollSummaries_origin = [
        ...state.PollSummaries_origin,
        ...newPollSummaries
      ]
      if (state.searchTag !== '') {
        state.PollSummaries = [
          ...state.PollSummaries,
          ...newPollSummaries.filter(v => {
            return typeof v.tags !== 'undefined'
              ? v.tags.some(e => e.name === state.searchTag)
              : false
          })
        ]
      } else {
        state.PollSummaries = [...state.PollSummaries, ...newPollSummaries]
      }
      isLoading = false
    }
    const getTags = async () => {
      try {
        state.tags = (await apis.getTags()).data
      } catch {
        state.tags = []
      }
    }
    getPolls()
    getTags()

    const calculateFilter = async () => {
      if (state.searchTag.length === 0) {
        state.autocompletes = []
      } else {
        state.autocompletes = state.tags
          .filter((v: PollTag) => {
            return v.name.indexOf(state.searchTag) === 0
          })
          .slice(0, 5)
      }
      state.PollSummaries = state.PollSummaries_origin
    }
    const onAutocomplete = async (str: string) => {
      state.autocompletes = []
      if (str.length === 0) return
      state.searchTag = str
      state.PollSummaries = []
      state.PollSummaries_origin = []
      offset = 0
      isEnd = false
      await getPolls()
    }
    const onKeyword = async () => {
      state.PollSummaries = []
      state.PollSummaries_origin = []
      offset = 0
      isEnd = false
      await getPolls()
    }

    const scrollHandler = async () => {
      const hasReached =
        window.innerHeight + window.scrollY >= document.body.offsetHeight
      if (hasReached) {
        await getPolls()
      }
    }

    onMounted(() => {
      window.addEventListener('scroll', scrollHandler)
    })

    onUnmounted(() => {
      window.removeEventListener('scroll', scrollHandler)
    })

    return {
      state,
      getPolls,
      calculateFilter,
      onAutocomplete,
      onKeyword,
      PollCardComponent
    }
  }
}
</script>
