<template>
  <div class="container">
    <h1><em class="bi bi-clipboard-data" /> ディテール</h1>
    <div class="card mx-auto my-3 p-2" style="border-radius: 1em">
      <div class="card-header bg-white">
        <div class="poll-title">{{ state.PollSummary.title }}</div>
        <div class="poll-tags bi bi-tags-fill text-muted d-flex flex-wrap">
          <span
            v-for="tag in state.PollSummary.tags"
            :key="tag.id"
            class="ms-1">
            {{ tag.name }},
          </span>
        </div>
      </div>
      <div class="card-body">
        <div class="">
          <PollResultComponent
            :poll-id="state.PollResult.pollId"
            :type="state.PollResult.type"
            :count="state.PollResult.count"
            :result="state.PollResult.result">
          </PollResultComponent>
        </div>
      </div>
      <div class="footer">
        <div class="d-flex flex-wrap justify-content-around mx-auto my-1">
          <div class="remain">{{ state.remain }}</div>
          <div class="created-at">作成日 : {{ state.createdAt }}</div>
        </div>
        <div class="d-flex flex-wrap justify-content-around mx-auto my-1">
          <div class="owner_name my-auto">
            @{{ state.PollSummary.owner.name }}
          </div>
          <button
            v-if="state.PollSummary.userStatus.isOwner && !state.outdated"
            type="button"
            class="close-poll btn btn-warning"
            @click="postPollsClose()">
            締め切る
          </button>
          <router-link
            v-if="state.PollSummary.userStatus.isOwner"
            :to="{ name: 'home' }">
            <button
              type="button"
              class="delete-poll btn btn-danger"
              @click="deletePolls()">
              削除
            </button>
          </router-link>
        </div>
      </div>
    </div>
    <h3 class="mt-4"><em class="bi bi-chat-left-text-fill" /> コメント一覧</h3>
    <table class="table table-sm table-striped table-bordered caption-top">
      <caption />
      <thead>
        <tr>
          <th scope="col">#</th>
          <th scope="col" class="text-start">コメント</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(PollComment, i) in state.PollComments"
          :key="PollComment.createdAt">
          <th scope="row">{{ i + 1 }}</th>
          <td class="text-start">{{ PollComment.content }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import PollResultComponent from '/@/components/PollResult.vue'
import { defineComponent, reactive, watchEffect } from 'vue'
import apis, {
  PollComment,
  PollResults,
  PollStatus,
  PollSummary,
  PollType,
  UserStatusAccessModeEnum
} from '/@/lib/apis'
import { useRoute } from 'vue-router'

interface State {
  pollId: string
  now: Date
  deadline: string
  remain: string
  createdAt: string
  outdated: boolean
  PollSummary: PollSummary
  PollResult: PollResults
  PollComments: PollComment[]
}

export default defineComponent({
  name: 'DetailsPage',
  components: { PollResultComponent },
  setup() {
    const state = reactive<State>({
      pollId: '',
      now: new Date(),
      deadline: '',
      remain: '',
      createdAt: '',
      outdated: true,
      PollSummary: {
        pollId: '',
        title: '',
        type: PollType.Radio,
        deadline: '',
        tags: [],
        question: [],
        createdAt: '',
        qStatus: PollStatus.Outdated,
        owner: {
          uuid: '',
          name: ''
        },
        userStatus: {
          isOwner: false,
          accessMode: UserStatusAccessModeEnum.CanAccessDetails
        }
      },
      PollResult: {
        pollId: '',
        type: PollType.Radio,
        count: 0,
        result: []
      },
      PollComments: []
    })

    const comp_remain = () => {
      if (state.deadline === '-1') {
        state.remain = '期限なし'
        return
      }
      const deadline = new Date(state.deadline)
      if (deadline.getTime() - state.now.getTime() <= 0 || state.outdated) {
        state.remain = '公開済み'
        return
      }
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
          ? '残り : ' + day.toString() + '日'
          : hour > 0
          ? '残り : ' + hour.toString() + '時間' + minute.toString() + '分'
          : '残り : ' + minute.toString() + '分'
    }
    comp_remain()
    setInterval(() => {
      state.now = new Date()
    }, 5000)
    watchEffect(() => {
      state.now, state.outdated, comp_remain()
    })

    const { pollId } = useRoute().params
    state.pollId = pollId.toString()
    const getPolls = async () => {
      try {
        state.PollSummary = (await apis.getPollsPollID(state.pollId)).data
      } catch {
        alert('投票を取得できませんでした。')
      }
      const date = new Date(state.PollSummary.createdAt)
      state.deadline =
        typeof state.PollSummary.deadline !== 'undefined'
          ? state.PollSummary.deadline
          : '-1'
      state.createdAt =
        date.getFullYear().toString() +
        '/' +
        (date.getMonth() + 1).toString() +
        '/' +
        date.getDate().toString() +
        ' ' +
        date.getHours().toString() +
        ':' +
        date.getMinutes().toString()
      state.outdated = state.PollSummary.qStatus === PollStatus.Outdated
    }
    const getResult = async () => {
      try {
        state.PollResult = (await apis.getPollsPollIDResults(state.pollId)).data
      } catch {
        alert('投票結果を取得できませんでした。')
      }
    }
    const getComments = async () => {
      try {
        state.PollComments = (
          await apis.getPollsPollIDComments(state.pollId)
        ).data
      } catch {
        alert('コメントを取得できませんでした。')
      }
    }
    const postPollsClose = async () => {
      try {
        apis.postPollsClose(state.pollId)
        state.outdated = true
      } catch {
        alert('投票を締め切ることができませんでした。')
      }
    }
    const deletePolls = async () => {
      try {
        apis.deletePollsPollID(state.pollId)
      } catch {
        alert('投票を削除できませんでした。')
      }
    }
    getPolls()
    getResult()
    getComments()

    return {
      state,
      PollResultComponent,
      postPollsClose,
      deletePolls
    }
  }
})
</script>

<style>
.card {
  width: 490px;
}
</style>
