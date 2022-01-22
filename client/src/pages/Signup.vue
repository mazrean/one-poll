<template>
  <div class="container position-relative top-50 start-50 translate-middle">
    <certification-form
      class="m-auto"
      sign="サインアップ"
      @on-submit-event="onSubmitForm" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import CertificationForm from '/@/components/CertificationForm.vue'
import api, { PostUser } from '/@/lib/apis'
import { useMainStore } from '/@/store/index'
import { useRouter } from 'vue-router'

export default defineComponent({
  name: 'SignupPage',
  components: { CertificationForm },
  setup() {
    const store = useMainStore()
    const router = useRouter()
    const onSubmitForm = async (name: string, password: string) => {
      const user: PostUser = { name: name, password: password }
      await api.postUsers(user)
      await store.setUserID()
      router.push({ path: '/', force: true })
    }
    return { onSubmitForm }
  }
})
</script>
