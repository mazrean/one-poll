<template>
  <div class="container">
    <certification-form class="m-auto" sign="サインイン" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import CertificationForm from '../components/CertificationForm.vue'
import api from '/@/lib/apis'

export default defineComponent({
  name: 'SigninPage',
  components: { CertificationForm }
})
</script>
