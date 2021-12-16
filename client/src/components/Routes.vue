<template>
  <div>
    <router-link :to="{ name: 'home' }">ホーム</router-link> |
    <router-link :to="{ name: 'signup' }">サインアップ</router-link> |
    <router-link :to="{ name: 'signin' }">サインイン</router-link> |
    <router-link :to="{ name: 'profile' }">プロフィール</router-link>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Routes',
  setup() {
    return {}
  }
})
</script>

<style module>
.re {
  font-weight: bold;
}
</style>
