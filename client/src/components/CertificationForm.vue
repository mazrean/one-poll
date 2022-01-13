<template>
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">PollQ</h5>
      <form
        class="needs-validation"
        :class="{ 'was-validated': Validated }"
        novalidate
        @submit="Validation">
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
            英数＋アンダーバー込で4~16文字を入力
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
            英数＋アンダーバー込で8~50文字を入力
          </div>
        </div>
        <button type="submit" class="btn btn-primary">{{ sign }}</button>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'CertificationFormComponent',
  props: {
    sign: {
      type: String,
      default: 'Sign in'
    }
  },
  setup() {
    return {}
  },
  data() {
    return {
      name: '',
      password: '',
      Validated: false,
      error: [false, false]
    }
  },
  watch: {
    name() {
      let re = new RegExp('[0-9a-zA-Z_]{4,16}')
      if (re.test(this.name)) {
        this.error[0] = false
      } else {
        this.error[0] = true
      }
    },
    password() {
      let re = new RegExp('[0-9a-zA-Z_]{8,50}')
      if (re.test(this.password)) {
        this.error[1] = false
      } else {
        this.error[1] = true
      }
    }
  },
  methods: {
    Validation: function (event: Event) {
      this.Validated = true
      this.error.forEach(error => {
        if (error) {
          event.preventDefault()
        }
      })
    }
  }
})
</script>
