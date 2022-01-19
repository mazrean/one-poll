<template>
  <div class="d-flex flex-wrap">
    <div v-for="PollSummary in state.PollSummaries" :key="PollSummary.poll_id">
      <PollCardComponent
        :poll-id="PollSummary.poll_id"
        :title="PollSummary.title"
        :type="PollSummary.type"
        :deadline="PollSummary.deadline"
        :question="PollSummary.question"
        :created-at="PollSummary.created_at"
        :q-status="PollSummary.q_status"
        :owner="PollSummary.owner"
        :user-status="PollSummary.user_status"
        class="m-3">
      </PollCardComponent>
    </div>
  </div>
</template>

<script lang="ts">
import PollCardComponent from '/@/components/PollCard.vue'
//import PollSummaries from '/@/assets/poll_summaries.json'
import apis, { PollSummary } from '../lib/apis'
import { reactive, onMounted } from 'vue'
interface State {
  PollSummaries: PollSummary[]
}
export default {
  name: 'HomePage',
  components: { PollCardComponent },
  setup() {
    const state = reactive<State>({
      PollSummaries: []
    })
    onMounted(async () => {
      state.PollSummaries = (await apis.getPolls()).data
    })
    return {
      state,
      PollCardComponent
    }
  }
}
</script>
