<template>
  <div class="container position-relative top-50 start-50 translate-middle">
    <certification-form
      ref="form"
      class="m-auto"
      sign="サインイン"
      :form-error="formError"
      @on-submit-event="onSubmitForm"
      @on-passkey-event="onPasskey" />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import CertificationForm from '/@/components/CertificationForm.vue'
import api, { PostUser, WebAuthnCredentialType } from '/@/lib/apis'
import { useMainStore } from '/@/store/index'
import { useRouter } from 'vue-router'
import { b64urlDecode, b64urlEncode } from '../lib/encoding'

export default defineComponent({
  name: 'SigninPage',
  components: { CertificationForm },
  setup() {
    const store = useMainStore()
    const router = useRouter()
    const formError = ref<boolean>(false)
    const form = ref<InstanceType<typeof CertificationForm>>()

    const onSubmitForm = async (name: string, password: string) => {
      const user: PostUser = { name: name, password: password }
      try {
        await api.postUsersSignin(user)
      } catch {
        //formの入力を消す
        form.value?.resetForm()
        return
      }
      await store.setUserID()
      router.push('/')
    }

    const onPasskey = async () => {
      const res = (await api.postWebauthnAuthenticateStart()).data
      const option = {
        ...res,
        challenge: b64urlDecode(res.challenge)
      }

      try {
        const credential = await navigator.credentials.get({
          publicKey: option
        })
        const { PublicKeyCredential, AuthenticatorAssertionResponse } = window
        if (
          !credential ||
          credential.type !== 'public-key' ||
          !(credential instanceof PublicKeyCredential) ||
          !(credential.response instanceof AuthenticatorAssertionResponse)
        ) {
          alert('認証に失敗しました。')
          return
        }

        const authRes = await api.postWebauthnAuthenticateFinish({
          id: credential.id,
          type: WebAuthnCredentialType.PublicKey,
          rawId: b64urlEncode(credential.rawId),
          response: {
            clientDataJSON: b64urlEncode(credential.response.clientDataJSON),
            authenticatorData: b64urlEncode(
              credential.response.authenticatorData
            ),
            signature: b64urlEncode(credential.response.signature)
          }
        })
        if (authRes.status !== 200) {
          alert('認証に失敗しました。')
          return
        }
      } catch {
        alert('認証に失敗しました。')
        return
      }

      await store.setUserID()
      router.push('/')
    }

    return { form, formError, onSubmitForm, onPasskey }
  }
})
</script>
