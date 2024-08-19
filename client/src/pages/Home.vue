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
          v-model="state.inputSearchTitle"
          type="searchTitle"
          class="form-control d-flex mx-3 my-1"
          name="searchTitle"
          placeholder="キーワードで検索 (Enterで更新)"
          @Input="state.inputSearchTitle = $event.target.value"
          @keydown.enter="state.searchTitle = state.inputSearchTitle" />
        <input
          v-model="state.inputSearchTag"
          type="searchTag"
          class="form-control d-flex mx-3 my-1"
          name="searchTag"
          placeholder="タグで検索 (タグ候補選択で更新)"
          @Input="onInputTag()" />
        <ul v-for="v in state.autocompletes" :key="v.id" class="list-group">
          <button
            class="list-group-item list-group-item-action p-1"
            type="button"
            @click="onSubmitTag(v.name)">
            <em class="bi bi-tags-fill" /> {{ v.name }}
          </button>
        </ul>
      </div>
      <div
        v-if="state.isLoading"
        class="spinner-border text-secondary m-3"
        role="status"></div>
      <div v-else-if="state.pollSummaries.length === 0" class="m-3">
        <p>表示可能な質問がありません。</p>
      </div>
      <div v-else class="d-flex flex-wrap justify-content-center">
        <div
          v-for="pollSummary in state.pollSummaries"
          :key="pollSummary.pollId">
          <PollCardComponent :poll="pollSummary" />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { watchEffect } from 'vue'
import PollCardComponent from '/@/components/PollCard.vue'
import apis, { PollSummary, PollTag } from '/@/lib/apis'
import { reactive, onMounted, onUnmounted } from 'vue'

interface State {
  pollSummaries: PollSummary[]
  isLoading: boolean
  searchTitle: string
  inputSearchTitle: string
  searchTag: string
  inputSearchTag: string
  autocompletes: PollTag[]
}

export default {
  name: 'HomePage',
  components: { PollCardComponent },
  setup() {
    const state = reactive<State>({
      pollSummaries: [],
      isLoading: false,
      searchTitle: '',
      inputSearchTitle: '',
      searchTag: '',
      inputSearchTag: '',
      autocompletes: []
    })

    const tagsPromise = (async () => {
      return (await apis.getTags()).data
    })()
    const onInputTag = async () => {
      if (!state.inputSearchTag) {
        state.autocompletes = []
        return
      }

      state.autocompletes = (await tagsPromise)
        .filter(v => v.name.startsWith(state.inputSearchTag))
        .slice(0, 5)
    }
    const onSubmitTag = async (str: string) => {
      state.autocompletes = []
      state.inputSearchTag = str
      state.searchTag = str
    }

    const limit = 10
    let offset = 0
    let isEnd = false
    const getPolls = async (searchTag: string, searchTitle: string) => {
      if (isEnd) return []

      let newPollSummaries: PollSummary[]
      try {
        newPollSummaries = (
          await apis.getPolls(
            limit,
            offset,
            state.inputSearchTitle || undefined
          )
        ).data
      } catch {
        newPollSummaries = []
      }
      if (newPollSummaries.length < limit) {
        isEnd = true
      }
      offset += newPollSummaries.length

      if (searchTag) {
        newPollSummaries = newPollSummaries.filter(
          v => v.tags && v.tags.some(e => e.name === state.inputSearchTag)
        )
      }

      if (searchTitle) {
        newPollSummaries = newPollSummaries.filter(v =>
          v.title.includes(state.inputSearchTitle)
        )
      }

      return newPollSummaries
    }

    const scrollHandler = async () => {
      const hasReached =
        window.innerHeight + Math.ceil(window.scrollY) >=
        document.body.offsetHeight
      if (hasReached) {
        state.isLoading = true
        state.pollSummaries = [
          ...state.pollSummaries,
          ...(await getPolls(state.inputSearchTag, state.inputSearchTitle))
        ]
        state.isLoading = false
      }
    }

    watchEffect(() => {
      state.pollSummaries = state.pollSummaries
        .filter(v => !state.searchTitle || v.title.includes(state.searchTitle))
        .filter(
          v => !state.searchTag || v.tags?.some(e => e.name === state.searchTag)
        )
      scrollHandler()
    })

    onMounted(() => {
      window.addEventListener('wheel', scrollHandler)
      window.addEventListener('touchmove', scrollHandler)
      window.addEventListener('scroll', scrollHandler)
    })

    onUnmounted(() => {
      window.removeEventListener('wheel', scrollHandler)
      window.removeEventListener('touchmove', scrollHandler)
      window.removeEventListener('scroll', scrollHandler)
    })

    return {
      state,
      getPolls,
      onInputTag,
      onSubmitTag,
      PollCardComponent
    }
  }
}
</script>
