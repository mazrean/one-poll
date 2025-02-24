<template>
  <div class="card m-0 p-2" style="border-radius: 1em">
    <div
      class="card-header text-start bg-white px-1 py-2 mw-100"
      style="width: 340px">
      <h4 class="card-title">{{ poll.title }}</h4>
      <div class="card-tags bi bi-tags-fill text-muted d-flex flex-wrap">
        <span
          v-for="tag in poll.tags ?? [{ id: -1, name: '' }]"
          :key="tag.id"
          class="ms-1">
          {{ tag.name }},
        </span>
      </div>
    </div>
    <div class="card-body px-1 py-2" style="width: 360px">
      <div v-if="state.onlyBrowsable">
        <button
          v-for="q in poll.question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1 w-100"
          disabled>
          {{ q.choice }}
        </button>
      </div>
      <div v-if="state.canAnswer">
        <button
          v-for="(q, i) in poll.question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1 w-100"
          @click="submitPollID(i)">
          {{ q.choice }}
        </button>
        <div class="m-2">
          <textarea
            v-model="state.comment"
            placeholder="コメント(任意)"
            maxlength="2000"
            class="form-control border-secondary">
          </textarea>
        </div>
      </div>
      <div v-if="state.canAccessDetails">
        <PollResultComponent
          :poll-id="state.pollResult.pollId"
          :type="state.pollResult.type"
          :count="state.pollResult.count"
          :result="state.pollResult.result">
        </PollResultComponent>
      </div>
    </div>
    <div class="footer d-flex justify-content-around mb-1">
      <div>{{ state.remain }}</div>
      <router-link
        v-if="state.canAccessDetails"
        class="link link-detail"
        :to="{ name: 'details', params: { pollId: poll.pollId } }">
        詳細
      </router-link>
      <div>@{{ poll.owner.name }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, reactive } from 'vue'
import PollResultComponent from '/@/components/PollResult.vue'
import apis, {
  UserStatusAccessModeEnum,
  PostPollId,
  PollResults,
  PollType,
  PollSummary
} from '/@/lib/apis'
import { watchEffect } from 'vue'

interface State {
  onlyBrowsable: boolean
  canAnswer: boolean
  canAccessDetails: boolean
  now: Date
  remain: string
  comment: string
  pollResult: PollResults
}

export default defineComponent({
  components: { PollResultComponent },
  props: {
    poll: {
      type: Object as PropType<PollSummary>,
      required: true
    }
  },
  setup(props) {
    const state = reactive<State>({
      onlyBrowsable: false,
      canAnswer: false,
      canAccessDetails: false,
      now: new Date(),
      remain: '',
      comment: '',
      pollResult: {
        pollId: '',
        type: PollType.Radio,
        count: 0,
        result: []
      }
    })

    state.onlyBrowsable =
      props.poll.userStatus.accessMode == UserStatusAccessModeEnum.OnlyBrowsable
    state.canAnswer =
      props.poll.userStatus.accessMode == UserStatusAccessModeEnum.CanAnswer
    state.canAccessDetails =
      props.poll.userStatus.accessMode ==
      UserStatusAccessModeEnum.CanAccessDetails

    const compRemain = (now: Date) => {
      if (!props.poll.deadline) {
        return '期限なし'
      }

      const deadline = new Date(props.poll.deadline)

      if (deadline.getTime() <= now.getTime()) {
        return '公開済み'
      }

      let dif = Math.floor((deadline.getTime() - now.getTime()) / (60 * 1000))
      const day = Math.floor(dif / 1440)
      dif %= 1440
      const hour = Math.floor(dif / 60)
      dif %= 60
      const minute = dif

      if (day > 0) {
        return `残り : ${day}日`
      } else if (hour > 0) {
        return `残り : ${hour}時間${minute}分`
      } else {
        return `残り : ${minute}分`
      }
    }

    const timer = setInterval(() => {
      state.now = new Date()
    }, 1000)
    watchEffect(() => {
      state.remain = compRemain(state.now)
      if (state.remain == '公開済み') {
        state.canAnswer = false
        state.canAccessDetails = true
        clearInterval(timer)
      }
    })

    const postPoll = async (pollID: PostPollId) => {
      try {
        await apis.postPollsPollID(props.poll.pollId, pollID)
      } catch {
        alert('投票できませんでした。時間を空けてもう一度お試しください。')
        return
      }
    }
    const getResult = async () => {
      try {
        state.pollResult = (
          await apis.getPollsPollIDResults(props.poll.pollId)
        ).data
      } catch {
        alert('投票結果を取得できませんでした。')
        return
      }
    }
    const submitPollID = async (i: number) => {
      const pollID: PostPollId = {
        answer: [props.poll.question[i].id],
        comment: state.comment
      }
      await postPoll(pollID)
      await getResult()
      state.canAnswer = false
      state.canAccessDetails = true
    }
    if (state.canAccessDetails) getResult()

    return {
      state,
      submitPollID,
      PollResultComponent
    }
  }
})
</script>

<style>
.vote-button {
  height: 30px;
}
</style>
