<template>
  <div class="card">
    <div class="card-header text-start">
      <h2 class="card-title">{{ title }}</h2>
    </div>
    <div class="card-body">
      <div v-if="state.only_browsable">
        <button
          v-for="q in question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1"
          disabled>
          {{ q.choice }}
        </button>
      </div>
      <div v-if="state.can_answer">
        <button
          v-for="(q, index) in question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1"
          @click="submitPollID(index)">
          {{ q.choice }}
        </button>
        <textarea
          :v-model="state.comment"
          placeholder="コメント"
          rows="3"
          cols="50"
          maxlength="2000"
          class="m-2"></textarea>
      </div>
      <div v-if="state.can_access_details">
        <PollResultComponent
          :poll-id="state.PollResult.pollId"
          :type="state.PollResult.type"
          :count="state.PollResult.count"
          :result="state.PollResult.result">
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
import { defineComponent, onMounted, PropType, reactive } from 'vue'
import PollResultComponent from '/@/components/PollResult.vue'
import apis, {
  Choice,
  User,
  UserStatus,
  UserStatusAccessModeEnum,
  PostPollId,
  PollResults,
  PollType
} from '../lib/apis'

interface State {
  only_browsable: boolean
  can_answer: boolean
  can_access_details: boolean
  comment: string
  PollResult: PollResults
}

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
    const state = reactive<State>({
      only_browsable: false,
      can_answer: false,
      can_access_details: false,
      comment: '',
      PollResult: {
        pollId: '',
        type: PollType.Radio,
        count: 0,
        result: []
      }
    })
    onMounted(async () => {
      if (state.can_access_details)
        state.PollResult = (await apis.getPollsPollIDResults(props.pollId)).data
    })
    state.only_browsable =
      props.userStatus.accessMode == UserStatusAccessModeEnum.OnlyBrowsable
    state.can_answer =
      props.userStatus.accessMode == UserStatusAccessModeEnum.CanAnswer
    state.can_access_details =
      props.userStatus.accessMode == UserStatusAccessModeEnum.CanAccessDetails
    const submitPollID = async (index: number) => {
      const pollID: PostPollId = {
        answer: [props.question[index].id],
        comment: state.comment
      }
      try {
        await apis.postPollsPollID(props.pollId, pollID)
      } catch {
        alert('投票できませんでした。時間を空けてもう一度お試しください。')
        return
      }
      try {
        state.PollResult = (await apis.getPollsPollIDResults(props.pollId)).data
      } catch {
        alert('投票結果を取得できませんでした。')
        return
      }
      state.can_answer = false
      state.can_access_details = true
    }
    return {
      PollResultComponent,
      state,
      submitPollID
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
