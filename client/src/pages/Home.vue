<template>
  <div class="container">
    <h1><em class="bi bi-house-fill" /> ホーム</h1>
    <div class="m-auto">
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
            :question="PollSummary.question"
            :created-at="PollSummary.createdAt"
            :q-status="PollSummary.qStatus"
            :owner="PollSummary.owner"
            :user-status="PollSummary.userStatus"
            class="m-3">
          </PollCardComponent>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import PollCardComponent from '/@/components/PollCard.vue'
import apis, { PollSummary } from '../lib/apis'
import { reactive, onMounted } from 'vue'
interface State {
  PollSummaries: PollSummary[]
  isLoading: boolean
}
export default {
  name: 'HomePage',
  components: { PollCardComponent },
  setup() {
    const state = reactive<State>({
      PollSummaries: [],
      isLoading: true
    })
    onMounted(async () => {
      try {
        state.PollSummaries = (await apis.getPolls()).data
      } catch {
        state.PollSummaries = []
      }
      state.isLoading = false
    })
    return {
      state,
      PollCardComponent
    }
  }
}
</script>
