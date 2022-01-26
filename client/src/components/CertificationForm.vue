<template>
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">{{ sign }}</h5>
      <form
        class="needs-validation"
        :class="{ 'was-validated': validated }"
        novalidate
        @submit="validation">
        <div>
          <label for="user_name" class="card-text">ユーザー名</label><br />
          <input
            id="name"
            v-model="name"
            type="text"
            class="form-control"
            placeholder="ユーザー名を入力"
            required
            pattern="[0-9a-zA-Z_]{4,16}" />
          <div class="invalid-feedback">{{ userNameErrorMessage }}</div>
        </div>
        <div>
          <label for="password" class="card-text">パスワード</label><br />
          <input
            id="password"
            v-model="password"
            type="password"
            class="form-control"
            placeholder="パスワードを入力"
            required
            pattern="[0-9a-zA-Z_]{8,50}" />
          <div class="invalid-feedback">{{ passwordErrorMessage }}</div>
        </div>
        <button type="submit" class="btn btn-primary mt-3">{{ sign }}</button>
      </form>
      <router-link
        v-if="sign === 'サインイン'"
        class="mt-1"
        :to="{ name: 'signup' }"
        >新しくアカウントを作成する</router-link
      >
    </div>
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  ref,
  toRefs,
  watch,
  watchEffect,
  defineExpose
} from 'vue'

export default defineComponent({
  name: 'CertificationFormComponent',
  props: {
    sign: {
      type: String,
      default: ''
    },
    userNameErrorMessage: {
      type: String,
      default: 'ユーザー名は4~16文字の英数字・アンダーバーにしてください'
    },
    passwordErrorMessage: {
      type: String,
      default: 'パスワードは8~50文字の英数字・アンダーバーにしてください'
    },
    errorMessage: {
      type: String,
      default: 'ユーザー名またはパスワードが間違っています'
    }
  },
  emits: ['on-submit-event'],
  setup(props, context) {
    const name = ref<string>('')
    const password = ref<string>('')
    const userNameErrorMessage = ref<string>(props.userNameErrorMessage)
    const passwordErrorMessage = ref<string>(props.passwordErrorMessage)
    const validated = ref<boolean>(false)
    const onSubmitForm = () => {
      context.emit('on-submit-event', name.value, password.value)
    }
    let nameError = true
    let passError = true
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
    // watch(
    //   () => props.formError,
    //   () => {
    //     if (formError) {
    //       name.value = ''
    //       password.value = ''
    //       userNameErrorMessage.value = props.errorMessage
    //       passwordErrorMessage.value = props.errorMessage
    //       formError.value = false
    //       console.log(props.formError)
    //       console.log(formError.value)
    //     }
    //   }
    // )
    //method
    const validation = (event: Event): void => {
      validated.value = true
      if (nameError || passError) {
        event.preventDefault()
      }
      //validationが通り、実際にAPIを投げる
      else {
        event.preventDefault()
        onSubmitForm()
      }
    }
    const resetForm = () => {
      name.value = ''
      password.value = ''
      userNameErrorMessage.value = props.errorMessage
      passwordErrorMessage.value = props.errorMessage
    }
    return {
      name,
      password,
      userNameErrorMessage,
      passwordErrorMessage,
      validated,
      onSubmitForm,
      validation,
      resetForm
    }
  }
})
</script>
