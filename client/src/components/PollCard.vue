<template>
  <div class="card">
    <div class="card-header text-start">
      <h2 class="card-title">{{ title }}</h2>
    </div>
    <div class="card-body">
      <div v-if="only_browsable">
        <button
          v-for="q in question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1"
          disabled>
          {{ q.choice }}
        </button>
      </div>
      <div v-if="can_answer">
        <button
          v-for="q in question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1"
          @click="vote()">
          {{ q.choice }}
        </button>
      </div>
      <div v-if="can_access_details">
        <PollResultComponent
          :poll-id="PollResults.pollId"
          :type="PollResults.type"
          :count="PollResults.count"
          :result="PollResults.result">
        </PollResultComponent>
      </div>
    </div>
    <div class="footer d-flex justify-content-around">
      <div>締切 : {{ deadline }}</div>
      <div>
        <a href="#">@{{ owner.name }}</a>
      </div>
      <div>作成日 : {{ createdAt }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import PollResultComponent from '/@/components/PollResult.vue'
import PollResults from '/@/assets/poll_results.json'
import { Choice, User, UserStatus, UserStatusAccessModeEnum } from '../lib/apis'

export default defineComponent({
  components: { PollResultComponent },
  props: {
    pollId: {
      type: String,
      required: true
    },
    title: {
      type: String,
      required: true
    },
    type: {
      type: String,
      required: true
    },
    deadline: {
      type: String,
      required: true
    },
    question: {
      type: Array as PropType<Choice[]>,
      required: true
    },
    createdAt: {
      type: String,
      required: true
    },
    qStatus: {
      type: String,
      required: true
    },
    owner: {
      type: Object as PropType<User>,
      required: true
    },
    userStatus: {
      type: Object as PropType<UserStatus>,
      required: true
    }
  },
  setup(props) {
    const only_browsable =
      props.userStatus.access_mode == UserStatusAccessModeEnum.OnlyBrowsable
    const can_answer =
      props.userStatus.access_mode == UserStatusAccessModeEnum.CanAnswer
    const can_access_details =
      props.userStatus.access_mode == UserStatusAccessModeEnum.CanAccessDetails
    const vote = () => {
      return
    }
    return {
      only_browsable,
      can_answer,
      can_access_details,
      vote,
      PollResults,
      PollResultComponent
    }
  }
})
</script>

<style>
.card {
  width: 500px;
}
.vote-button {
  width: 90%;
  height: 2.7rem;
}
</style>
