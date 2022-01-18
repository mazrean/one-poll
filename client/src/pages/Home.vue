<template>
  <div class="d-flex flex-wrap">
    <PollCardComponent
      v-for="PollSummary in PollSummaries"
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
  </div>
</template>

<script lang="ts">
import PollCardComponent from '/@/components/PollCard.vue'
import PollSummaries from '/@/assets/poll_summaries.json'
//import apis from '../lib/apis'

export default {
  name: 'HomePage',
  components: { PollCardComponent },
  setup() {
    //const PollSummaries = apis.getPolls()
    return {
      PollSummaries,
      PollCardComponent
    }
  }
}
</script>
