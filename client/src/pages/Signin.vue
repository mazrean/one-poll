<template>
  <div class="container position-relative top-50 start-50 translate-middle">
    <certification-form
      ref="form"
      class="m-auto"
      sign="サインイン"
      :formError="formError"
      @on-submit-event="onSubmitForm" />
  </div>
</template>

<script lang="ts">
import { Component, defineComponent, ref } from 'vue'
import CertificationForm from '/@/components/CertificationForm.vue'
import api, { PostUser } from '/@/lib/apis'
import { useMainStore } from '/@/store/index'
import { useRouter } from 'vue-router'

export default defineComponent({
  name: 'SigninPage',
  components: { CertificationForm },
  setup() {
    const store = useMainStore()
    const router = useRouter()
    const errorMessage = 'ユーザー名またはパスワードが間違っています'
    const formError = ref<Boolean>(false)
    const form = ref<InstanceType<typeof CertificationForm>>()
    const onSubmitForm = async (name: string, password: string) => {
      const user: PostUser = { name: name, password: password }
      try {
        await api.postUsersSignin(user)
      } catch (err) {
        //formの入力を消す
        form.value?.resetForm()
        console.log('postに失敗したよ')
        return
      }
      await store.setUserID()
      router.push('/')
    }
    return { form, errorMessage, formError, onSubmitForm }
  }
})
</script>
