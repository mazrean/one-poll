<template>
  <div>
    <router-link :to="{ name: 'home' }">ホーム</router-link> |
    <router-link :to="{ name: 'signup' }">サインアップ</router-link> |
    <router-link :to="{ name: 'signin' }">サインイン</router-link> |
    <router-link :to="{ name: 'profile' }">プロフィール</router-link>
    <router-link :to="{ name: 'details' }">ディテール</router-link>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'RoutesComponent',
  setup() {
    return {}
  }
})
</script>
