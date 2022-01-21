<template>
  <div class="container">
    <h1><em class="bi bi-person-fill" /> プロフィール</h1>
    <ul id="myTab" class="nav nav-tabs" role="tablist">
      <li class="nav-item" role="presentation">
        <button
          id="home-tab"
          class="nav-link active"
          data-bs-toggle="tab"
          data-bs-target="#home"
          type="button"
          role="tab"
          aria-controls="home"
          aria-selected="true">
          アカウント情報
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button
          id="profile-tab"
          class="nav-link"
          data-bs-toggle="tab"
          data-bs-target="#profile"
          type="button"
          role="tab"
          aria-controls="profile"
          aria-selected="false">
          作成した質問一覧
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button
          id="contact-tab"
          class="nav-link"
          data-bs-toggle="tab"
          data-bs-target="#contact"
          type="button"
          role="tab"
          aria-controls="contact"
          aria-selected="false">
          回答した質問一覧
        </button>
      </li>
    </ul>
    <div id="myTabContent" class="tab-content">
      <div
        id="home"
        class="tab-pane fade show active"
        role="tabpanel"
        aria-labelledby="home-tab">
        ログイン中のアカウントID<br />
        <strong>@{{ userID }}</strong>
      </div>
      <div
        id="profile"
        class="tab-pane fade"
        role="tabpanel"
        aria-labelledby="profile-tab">
        <div class="m-auto">
          <div
            v-if="state.isLoading[0]"
            class="spinner-border text-secondary"
            role="status"></div>
          <div v-else-if="state.PollOwners.length === 0">
            <p>表示可能な質問がありません。</p>
          </div>
          <div v-else class="d-flex flex-wrap justify-content-center">
            <div
              v-for="PollSummary in state.PollOwners"
              :key="PollSummary.pollId">
              <PollCardComponent
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
          </div>
        </div>
      </div>
      <div
        id="contact"
        class="tab-pane fade"
        role="tabpanel"
        aria-labelledby="contact-tab">
        <div class="m-auto">
          <div
            v-if="state.isLoading[1]"
            class="spinner-border text-secondary"
            role="status"></div>
          <div v-else-if="state.PollAnswers.length === 0">
            <p>表示可能な質問がありません。</p>
          </div>
          <div v-else class="d-flex flex-wrap justify-content-center">
            <div
              v-for="PollSummary in state.PollAnswers"
              :key="PollSummary.pollId">
              <PollCardComponent
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
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed, reactive, onMounted } from 'vue'
import { useMainStore } from '/@/store/index'
import PollCardComponent from '/@/components/PollCard.vue'
import apis, { PollSummary } from '../lib/apis'

interface State {
  PollOwners: PollSummary[]
  PollAnswers: PollSummary[]
  isLoading: boolean[]
}
export default defineComponent({
  name: 'ProfilePage',
  components: {},
  setup() {
    const state = reactive<State>({
      PollOwners: [],
      PollAnswers: [],
      isLoading: [true, true]
    })
    const store = useMainStore()
    const userID = computed(() => store.userID)
    onMounted(async () => {
      try {
        state.PollOwners = (await apis.getUsersMeOwners()).data
      } catch {
        state.PollOwners = []
      }
      state.isLoading[0] = false
    })
    onMounted(async () => {
      try {
        state.PollAnswers = (await apis.getUsersMeAnswers()).data
      } catch {
        state.PollAnswers = []
      }
      state.isLoading[1] = false
    })
    return {
      state,
      PollCardComponent,
      userID
    }
  }
})
</script>
