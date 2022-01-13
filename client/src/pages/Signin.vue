<template>
  <div class="container">
    <!-- <Routes /> -->
    <certification-form sign="Sign in" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import CertificationForm from '../components/CertificationForm.vue'

export default defineComponent({
  name: 'SigninPage',
  components: { CertificationForm }
})
</script>
