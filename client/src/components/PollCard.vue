<template>
  <div class="card">
    <div class="card-header text-start">
      <h4 class="card-title">{{ title }}</h4>
      <div class="card-tags bi bi-tags-fill text-muted d-flex flex-wrap">
        <span v-for="tag in tags" :key="tag.id" class="ms-1">
          {{ tag.name }},
        </span>
      </div>
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
          v-for="(q, i) in question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1"
          @click="submitPollID(i)">
          {{ q.choice }}
        </button>
        <textarea
          v-model="state.comment"
          placeholder="コメント(任意)"
          rows="3"
          cols="50"
          maxlength="2000"
          class="m-2">
        </textarea>
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
      <div>残り: {{ state.remain }}</div>
      <div v-if="state.can_access_details">
        <router-link :to="{ name: 'details', params: { pollId: pollId } }">
          詳細を見る
        </router-link>
      </div>
      <div v-else>詳細を見る</div>
      <div>
        <a href="#">@{{ owner.name }}</a>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, reactive, watch } from 'vue'
import PollResultComponent from '/@/components/PollResult.vue'
import apis, {
  Choice,
  User,
  UserStatus,
  UserStatusAccessModeEnum,
  PostPollId,
  PollResults,
  PollType,
  PollTag
} from '/@/lib/apis'

interface State {
  only_browsable: boolean
  can_answer: boolean
  can_access_details: boolean
  now: Date
  remain: string
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
    tags: {
      type: Array as PropType<PollTag[]>,
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
      now: new Date(),
      remain: '',
      comment: '',
      PollResult: {
        pollId: '',
        type: PollType.Radio,
        count: 0,
        result: []
      }
    })

    state.only_browsable =
      props.userStatus.accessMode == UserStatusAccessModeEnum.OnlyBrowsable
    state.can_answer =
      props.userStatus.accessMode == UserStatusAccessModeEnum.CanAnswer
    state.can_access_details =
      props.userStatus.accessMode == UserStatusAccessModeEnum.CanAccessDetails

    const deadline = new Date(props.deadline)
    const comp_remain = () => {
      let dif = Math.floor(
        (deadline.getTime() - state.now.getTime()) / (60 * 1000)
      )
      const day = Math.floor(dif / 1440)
      dif %= 1440
      const hour = Math.floor(dif / 60)
      dif %= 60
      const minute = dif
      state.remain =
        day > 0
          ? day.toString() + '日'
          : hour > 0
          ? hour.toString() + '時間' + minute.toString() + '分'
          : minute.toString() + '分'
    }
    comp_remain()
    setInterval(() => {
      state.now = new Date()
    }, 5000)
    watch(
      () => state.now,
      () => comp_remain()
    )

    const postPoll = async (pollID: PostPollId) => {
      try {
        await apis.postPollsPollID(props.pollId, pollID)
      } catch {
        alert('投票できませんでした。時間を空けてもう一度お試しください。')
        return
      }
    }
    const getResult = async () => {
      try {
        state.PollResult = (await apis.getPollsPollIDResults(props.pollId)).data
      } catch {
        alert('投票結果を取得できませんでした。')
        return
      }
    }
    const submitPollID = async (i: number) => {
      const pollID: PostPollId = {
        answer: [props.question[i].id],
        comment: state.comment
      }
      await postPoll(pollID)
      await getResult()
      state.can_answer = false
      state.can_access_details = true
    }

    if (state.can_access_details) getResult()

    return {
      state,
      submitPollID,
      PollResultComponent
    }
  }
})
</script>

<style>
.card {
  width: 490px;
}
.vote-button {
  width: 420px;
  height: 30px;
}
</style>
