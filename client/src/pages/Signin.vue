<template>
  <div class="container">
    <!-- <Routes /> -->
    <certification-form Sign="Sign in"/>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Routes from '../components/Routes.vue'
import CertificationForm from '../components/CertificationForm.vue'

export default defineComponent({
  name: 'SigninPage',
  components: { Routes, CertificationForm }
})
</script>
