<template>
  <div class="d-flex flex-wrap">
    <PollCardComponent
      v-for="card in cards"
      :key="card.pollId"
      :poll-id="card.pollId"
      :title="card.title"
      :tags="card.tags"
      :question="card.question"
      :count="card.count"
      :created-at="card.createdAt"
      :owner="card.owner"
      :user-status="card.userStatus"
      class="m-3">
    </PollCardComponent>
  </div>
</template>

<script lang="ts">
import PollCardComponent from '/@/components/PollCard.vue'
import cards from '/@/assets/poll_card_data.json'

export default {
  name: 'HomePage',
  components: { PollCardComponent },
  setup() {
    return {
      cards,
      PollCardComponent
    }
  }
}
</script>
