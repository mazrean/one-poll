<template>
  <div class="d-flex flex-wrap">
    {{ state.PollSummaries }}
    <PollCardComponent
      v-for="PollSummary in state.PollSummaries"
      :key="PollSummary.pollId"
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
    <div v-for="v in state.PollSummaries" :key="v.pollId">
      <p>{{ v.title }}</p>
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
