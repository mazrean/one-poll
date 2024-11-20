<template>
  <div class="container">
    <h1><em class="bi bi-house-fill" /> ホーム</h1>
    <div class="d-flex flex-wrap justify-content-center m-auto gap-1">
      <svg
        v-show="false"
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        fill="currentColor"
        class="bi bi-search mx-1 my-auto"
        viewBox="0 0 16 16">
        <path
          d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z" />
      </svg>
      <input
        type="searchTitle"
        class="form-control d-flex m-0"
        name="searchTitle"
        placeholder="キーワードで検索"
        @keydown.enter="
          searchTitle = ($event.target as HTMLInputElement | null)?.value ?? ''
        " />
      <input
        v-model="inputTag"
        type="searchTag"
        class="form-control d-flex m-0"
        name="searchTag"
        placeholder="タグで検索" />
      <ul v-for="v in autocompletes" :key="v.id" class="list-group">
        <button
          class="list-group-item list-group-item-action p-1"
          type="button"
          @click="
            () => {
              searchTag = v.name
              inputTag = v.name
            }
          ">
          <em class="bi bi-tags-fill" /> {{ v.name }}
        </button>
      </ul>
      <div
        v-if="pollSummaries != null && pollSummaries.length === 0"
        class="m-3">
        <p>表示可能な質問がありません。</p>
      </div>
      <div v-else class="d-flex flex-wrap justify-content-center gap-1 w-100">
        <div
          v-for="pollSummary in pollSummaries"
          :key="`${pollSummary.pollId}:${pollSummary.userStatus.accessMode}:${pollSummary.userStatus.isOwner}`"
          class="mw-100">
          <PollCardComponent :poll="pollSummary" />
        </div>
      </div>
      <div
        v-if="isLoading"
        class="spinner-border text-secondary m-3"
        role="status" />
    </div>
  </div>
</template>

<script lang="ts">
import { useMainStore } from '/@/store'
import PollCardComponent from '/@/components/PollCard.vue'
import apis, { PollSummary, PollTag } from '/@/lib/apis'
import { onMounted, onUnmounted, computed } from 'vue'
import { watch } from 'vue'
import { ref } from 'vue'
import { watchEffect } from 'vue'

export default {
  name: 'HomePage',
  components: { PollCardComponent },
  setup() {
    const store = useMainStore()
    const isLogined = computed(() => !!store.userID)

    const tags = ref<PollTag[] | null>(null)
    ;(async () => {
      tags.value = (await apis.getTags()).data
    })()

    const searchTitle = ref<string>('')

    const limit = 10
    let searchState = {
      offset: 0,
      isEnd: false
    }
    const isLoading = ref(false)
    const pollSummaries = ref<PollSummary[] | null>(null)

    const getPolls = async (searchTitle: string, isLogined: boolean) => {
      if (searchState.isEnd || isLoading.value) return
      isLoading.value = true

      const publicPollPromise = apis.getPolls(
        limit,
        searchState.offset,
        searchTitle || undefined,
        true
      )

      const originalPollSummaries = pollSummaries.value
      let newPollSummaries: PollSummary[] = []
      if (isLogined) {
        const privatePollPromise = apis.getPolls(
          limit,
          searchState.offset,
          searchTitle || undefined,
          false
        )
        newPollSummaries = (
          await Promise.race([publicPollPromise, privatePollPromise])
        ).data
        privatePollPromise.then(res => {
          pollSummaries.value =
            originalPollSummaries?.concat(res.data) ?? res.data
        })
      } else {
        newPollSummaries = (await publicPollPromise).data
      }

      searchState = {
        offset: searchState.offset + newPollSummaries.length,
        isEnd: newPollSummaries.length < limit
      }
      pollSummaries.value =
        originalPollSummaries?.concat(newPollSummaries) ?? newPollSummaries
      isLoading.value = false
    }
    getPolls('', isLogined.value)

    watch([isLogined, searchTitle], async () => {
      pollSummaries.value = null
      searchState = {
        offset: 0,
        isEnd: false
      }
      await getPolls(searchTitle.value, isLogined.value)
    })

    const scrollHandler = async () => {
      const hasReached =
        window.innerHeight + Math.ceil(window.scrollY) >=
        document.body.offsetHeight
      if (hasReached) {
        await getPolls(searchTitle.value, isLogined.value)
      }
    }

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

    const inputTag = ref<string>('')
    const searchTag = ref<string>('')

    const tagFilterdPollSummaries = ref<PollSummary[] | null>(null)
    watchEffect(() => {
      tagFilterdPollSummaries.value =
        pollSummaries.value?.filter(v =>
          v.tags?.some(e => !searchTag.value || e.name === searchTag.value)
        ) ?? []
    })

    const autocompletes = ref<PollTag[]>([])
    watchEffect(() => {
      if (inputTag.value === searchTag.value) return

      searchTag.value = ''

      if (!inputTag.value) {
        autocompletes.value = []
        return
      }

      autocompletes.value =
        tags.value
          ?.filter(v => v.name.startsWith(inputTag.value))
          .slice(0, 5) ?? []
    })

    return {
      searchTitle,
      pollSummaries: tagFilterdPollSummaries,
      isLoading,
      inputTag,
      searchTag,
      autocompletes
    }
  }
}
</script>
