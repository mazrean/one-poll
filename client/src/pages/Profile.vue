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
          :class="`nav-link ${tab.active ? 'active' : ''}`"
          data-bs-toggle="tab"
          :data-bs-target="`#${tab.id}`"
          type="button"
          role="tab"
          :aria-controls="tab.id"
          :aria-selected="tab.active">
          {{ tab.label }}
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
        id="passkey"
        class="tab-pane fade"
        role="tabpanel"
        aria-labelledby="passkey-tab">
        <div class="m-auto">
          <div class="m-3">
            <button
              v-if="isPasskeyEnabled"
              type="button"
              class="btn btn-lg btn-primary"
              @click="onPasskeyRegister">
              <strong>パスキーを追加する</strong>
            </button>
          </div>
          <div
            v-if="!state.passkeys"
            class="spinner-border text-secondary m-3"
            role="status"></div>
          <div v-else-if="state.passkeys.length === 0" class="m-3">
            <p>パスキーがありません。</p>
          </div>
          <div v-else class="d-flex flex-wrap justify-content-center">
            <div v-for="passkey in state.passkeys" :key="passkey.id">
              <PasskeyCardComponent
                :passkey="passkey"
                @delete="onPasskeyDelete" />
            </div>
          </div>
        </div>
      </div>
      <div
        id="profile"
        class="tab-pane fade"
        role="tabpanel"
        aria-labelledby="profile-tab">
        <div class="m-auto">
          <div
            v-if="state.pollOwners === null"
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
            v-if="state.pollAnswers === null"
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
import apis, {
  PollSummary,
  WebAuthnCredential,
  WebAuthnCredentialType
} from '/@/lib/apis'
import { b64urlDecode, b64urlEncode, uuid2bytes } from '/@/lib/encoding'
import PasskeyCardComponent from '/@/components/PasskeyCard.vue'

interface State {
  passkeys: WebAuthnCredential[] | null
  pollOwners: PollSummary[] | null
  pollAnswers: PollSummary[] | null
}

export default defineComponent({
  name: 'ProfilePage',
  components: { PollCardComponent, PasskeyCardComponent },
  setup() {
    const store = useMainStore()
    const userID = computed(() => store.userID)

    const state = reactive<State>({
      passkeys: null,
      pollOwners: null,
      pollAnswers: null
    })
    onMounted(async () => {
      try {
        state.passkeys = (await apis.getWebauthnCredentials()).data
      } catch {
        state.passkeys = []
      }
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

    const onPasskeyRegister = async () => {
      if (!window.PublicKeyCredential) {
        alert('このブラウザは対応していません。')
        return
      }

      const res = (await apis.postWebauthnResisterStart()).data
      const option = {
        ...res,
        user: {
          ...res.user,
          id: uuid2bytes(res.user.id)
        },
        excludeCredentials: res.excludeCredentials?.map(cred => ({
          ...cred,
          id: b64urlDecode(cred.id)
        })),
        challenge: b64urlDecode(res.challenge),
        authenticatorSelection: {
          ...res.authenticatorSelection,
          userVerification: 'required' as const
        }
      }

      try {
        const credential = await navigator.credentials.create({
          publicKey: option
        })
        const { PublicKeyCredential, AuthenticatorAttestationResponse } = window
        if (
          !credential ||
          credential.type !== 'public-key' ||
          !(credential instanceof PublicKeyCredential) ||
          !(credential.response instanceof AuthenticatorAttestationResponse)
        ) {
          alert('登録に失敗しました。')
          return
        }

        const resisterRes = await apis.postWebauthnResisterFinish({
          id: credential.id,
          type: WebAuthnCredentialType.PublicKey,
          rawId: b64urlEncode(credential.rawId),
          response: {
            attestationObject: b64urlEncode(
              credential.response.attestationObject
            ),
            clientDataJSON: b64urlEncode(credential.response.clientDataJSON)
          }
        })
        if (resisterRes.status !== 200) {
          alert('登録に失敗しました。')
          return
        }
      } catch (e) {
        if (e instanceof Error) {
          if (e.name === 'InvalidStateError') {
            alert('既に登録されているパスキーです。')
            return
          }
        }
      }

      state.passkeys = (await apis.getWebauthnCredentials()).data
    }

    const onPasskeyDelete = async (id: string) => {
      const res = await apis.deleteWebauthnCredentials(id)
      if (res.status !== 200) {
        alert('削除に失敗しました。')
        return
      }

      state.passkeys =
        state.passkeys?.filter(passkey => passkey.id !== id) ?? null
    }

    let tabs = [
      { id: 'home', label: 'アカウント情報', active: true },
      { id: 'passkey', label: 'パスキー一覧', active: false },
      { id: 'profile', label: '作成した質問一覧', active: false },
      { id: 'contact', label: '回答した質問一覧', active: false }
    ]

    return {
      state,
      tabs,
      userID,
      onPasskeyRegister,
      onPasskeyDelete,
      isPasskeyEnabled: !!window.PublicKeyCredential
    }
  }
})
</script>
