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
        class="tab-pane fade show active m-3"
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
            v-if="!state.pollOwners"
            class="spinner-border text-secondary m-3"
            role="status"></div>
          <div v-else-if="state.pollOwners.length === 0" class="m-3">
            <p>表示可能な質問がありません。</p>
          </div>
          <div v-else class="d-flex flex-wrap justify-content-center">
            <div
              v-for="pollSummary in state.pollOwners"
              :key="pollSummary.pollId">
              <PollCardComponent :poll="pollSummary" />
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
            v-if="!state.pollAnswers"
            class="spinner-border text-secondary m-3"
            role="status"></div>
          <div v-else-if="state.pollAnswers.length === 0" class="m-3">
            <p>表示可能な質問がありません。</p>
          </div>
          <div v-else class="d-flex flex-wrap justify-content-center">
            <div
              v-for="pollSummary in state.pollAnswers"
              :key="pollSummary.pollId">
              <PollCardComponent :poll="pollSummary" />
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
import apis, { PollSummary } from '/@/lib/apis'

interface State {
  pollOwners: PollSummary[] | null
  pollAnswers: PollSummary[] | null
}

export default defineComponent({
  name: 'ProfilePage',
  components: { PollCardComponent },
  setup() {
    const store = useMainStore()
    const userID = computed(() => store.userID)

    const state = reactive<State>({
      pollOwners: null,
      pollAnswers: null
    })
    onMounted(async () => {
      try {
        state.pollOwners = (await apis.getUsersMeOwners()).data
      } catch {
        state.pollOwners = []
      }
    })
    onMounted(async () => {
      try {
        state.pollAnswers = (await apis.getUsersMeAnswers()).data
      } catch {
        state.pollAnswers = []
      }
    })

    return {
      state,
      userID
    }
  }
})
</script>
