<template>
  <div class="container">
    <h1><em class="bi bi-clipboard-data" /> ディテール</h1>
    <div class="card">
      <div class="card-header">
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
        <div class="created-at">作成日 : {{ state.PollSummary.createdAt }}</div>
        <div class="d-flex flex-wrap justify-content-around">
          <div class="owner_name">@{{ state.PollSummary.owner.name }}</div>
          <button
            v-if="state.PollSummary.userStatus.isOwner && !state.outdated"
            type="button"
            class="close-poll btn btn-warning"
            @click="postPollsClose()">
            締め切る
          </button>
          <router-link :to="{ name: 'home' }">
            <button
              v-if="state.PollSummary.userStatus.isOwner"
              type="button"
              class="delete-poll btn btn-danger"
              @click="deletePolls()">
              削除
            </button>
          </router-link>
        </div>
      </div>
    </div>
    <table class="table table-sm table-striped table-bordered caption-top mt-4">
      <caption>
        コメント一覧
      </caption>
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
import { defineComponent, reactive } from 'vue'
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
  outdated: boolean
  PollSummary: PollSummary
  PollResult: PollResults
  PollComments: PollComment[]
  createdAt: Date
}

export default defineComponent({
  name: 'DetailsPage',
  components: { PollResultComponent },
  setup() {
    const state = reactive<State>({
      pollId: '',
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
      PollComments: [],
      createdAt: new Date()
    })
    state.createdAt = new Date(state.PollSummary.createdAt)

    const { pollId } = useRoute().params
    state.pollId = pollId.toString()
    const getPolls = async () => {
      try {
        state.PollSummary = (await apis.getPollsPollID(state.pollId)).data
      } catch {
        alert('投票を取得できませんでした。')
      }
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
        alert('投票を締め切りました。')
      } catch {
        alert('投票を締め切ることができませんでした。')
      }
    }
    const deletePolls = async () => {
      try {
        apis.deletePollsPollID(state.pollId)
        alert('投票を削除しました。')
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
