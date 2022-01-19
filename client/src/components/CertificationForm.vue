<template>
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">PollQ</h5>
      <form
        class="needs-validation"
        :class="{ 'was-validated': validated }"
        novalidate
        @submit="validation">
        <div>
          <label for="user_name" class="card-text">User</label><br />
          <input
            id="name"
            v-model="name"
            type="text"
            class="form-control"
            placeholder="ユーザー名を入力"
            required
            pattern="[0-9a-zA-Z_]{4,16}" />
          <div class="invalid-feedback">
            ユーザー名は英数＋アンダーバー込で4~16文字にしてください
          </div>
        </div>
        <div>
          <label for="password" class="card-text">password</label><br />
          <input
            id="password"
            v-model="password"
            type="password"
            class="form-control"
            placeholder="パスワードを入力"
            required
            pattern="[0-9a-zA-Z_]{8,50}" />
          <div class="invalid-feedback">
            パスワードは英数＋アンダーバー込で8~50文字にしてください
          </div>
        </div>
        <button type="submit" class="btn btn-primary">{{ sign }}</button>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from 'vue'

export default defineComponent({
  name: 'CertificationFormComponent',
  props: {
    sign: {
      type: String,
      default: 'Sign in'
    }
  },
  setup() {
    const name = ref<string>('')
    const password = ref<string>('')
    const validated = ref<boolean>(false)
    let nameError = false
    let passError = false
    watch(name, () => {
      let re = new RegExp('[0-9a-zA-Z_]{4,16}')
      if (re.test(name.value)) {
        nameError = false
      } else {
        nameError = true
      }
    })
    watch(password, () => {
      let re = new RegExp('[0-9a-zA-Z_]{8,50}')
      if (re.test(password.value)) {
        passError = false
      } else {
        passError = true
      }
    })
    //method
    const validation = (event: Event): void => {
      validated.value = true
      if (nameError || passError) {
        event.preventDefault()
      }
    }
    return { name, password, validated, validation }
  }
})
</script>