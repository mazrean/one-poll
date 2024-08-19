<template>
  <div class="container">
    <h1><em class="bi bi-person-fill" /> プロフィール</h1>
    <ul id="myTab" class="nav nav-tabs" role="tablist">
      <li
        v-for="tab in tabs"
        :key="tab.id"
        class="nav-item"
        role="presentation">
        <button
          :id="`${tab.id}-tab`"
          class="nav-link active"
          data-bs-toggle="tab"
          :data-bs-target="`#${tab.id}`"
          type="button"
          role="tab"
          :aria-controls="tab.id"
          :aria-selected="tab.selected">
          {{ tab.name }}
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

    const tabs = [
      { id: 'home', name: 'アカウント情報', selected: true },
      { id: 'profile', name: '作成した質問一覧', selected: false },
      { id: 'contact', name: '回答した質問一覧', selected: false }
    ]

    return {
      state,
      userID,
      tabs
    }
  }
})
</script>
