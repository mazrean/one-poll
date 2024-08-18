<template>
  <div class="container position-relative top-50 start-50 translate-middle">
    <certification-form
      ref="form"
      class="m-auto"
      sign="サインイン"
      :form-error="formError"
      @on-submit-event="onSubmitForm" />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
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

    return { form, formError, onSubmitForm }
  }
})
</script>
