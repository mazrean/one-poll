<template>
  <!-- Modal -->
  <div
    id="newPollModal"
    class="modal fade"
    tabindex="-1"
    aria-labelledby="newPollModalLabel"
    aria-hidden="true">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h4 class="modal-title">新しい質問</h4>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form action="/hoge" method="POST" @submit.prevent>
            <label for="title" class="col-sm-2 col-form-label">タイトル</label>
            <input
              id="title"
              v-model="state.title"
              type="text"
              class="form-control mb-3"
              name="title"
              maxlength="32"
              placeholder="タイトルを入力" />
            <label for="detail" class="col-sm-2 col-form-label">質問文</label>
            <textarea
              id="detail"
              v-model="state.detail"
              class="form-control mb-3"
              name="detail"
              maxlength="140"
              placeholder="質問する内容を入力" />
            <label for="options" class="col-sm-2 col-form-label">選択肢</label>
            <ul
              v-for="(v, i) in state.options"
              :key="i"
              class="nav nav-pills flex-column">
              <li class="nav-item">
                <div class="d-inline-block d-flex">
                  <input
                    :id="'option-' + i"
                    v-model="state.options[i]"
                    type="text"
                    class="form-control d-flex my-1"
                    :name="'option-' + i"
                    maxlength="25"
                    placeholder="オプションの内容を入力" />
                  <button
                    :disabled="i < 2"
                    type="button"
                    class="btn btn-close d-flex my-auto ml-3"
                    @Click="deleteOption(i)"></button>
                </div>
              </li>
            </ul>
            <button
              v-show="state.options.length < 5"
              class="btn link"
              type="button"
              @Click="insertOption()">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="32"
                height="32"
                fill="currentColor"
                class="bi bi-plus"
                viewBox="0 0 16 16">
                <path
                  d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z" />
              </svg>
              <span class="align-middle">選択肢を追加する</span></button
            ><br />

            <label for="tag" class="col-sm-2 col-form-label">タグ</label>
            <div>
              <ul
                v-for="v in state.newTags"
                :key="v"
                class="nav nav-pills d-inline-block">
                <li class="nav-item d-flex">
                  <span
                    class="rounded-3 px-3 py-1 bg-light text-dark d-flex mx-1 my-1 fw-normal"
                    ><em class="bi bi-tags-fill" /> {{ v
                    }}<button
                      type="button"
                      class="btn btn-sm btn-close"
                      @Click="deleteTag(v)"></button
                  ></span>
                </li>
              </ul>
            </div>
            <input
              id="tag"
              v-model="state.newTag"
              :disabled="state.newTags.size >= 10"
              type="text"
              class="form-control mt-2"
              name="tag"
              maxlength="16"
              placeholder="新しく追加するタグ名を入力"
              @Input="calculateFilter()"
              @keydown.enter="insertTag()" />
            <ul v-for="v in state.autocompletes" :key="v" class="list-group">
              <button
                class="list-group-item list-group-item-action p-1"
                type="button"
                @click="onAutocomplete(v)">
                <em class="bi bi-tags-fill" /> {{ v }}
              </button>
            </ul>
            <label for="deadline" class="col-sm-2 col-form-label">締切</label>
            <select
              id="deadline"
              v-model="state.deadline"
              name="detail"
              class="form-select mb-3">
              <option value="0">1時間後</option>
              <option value="1">6時間後</option>
              <option value="2">12時間後</option>
              <option value="3">1日後</option>
              <option value="4">1週間後</option>
              <option value="5">設定しない</option>
            </select>
            <button
              type="button"
              class="btn btn-primary"
              data-bs-dismiss="modal"
              @click="submitPoll()">
              投稿する <em class="bi bi-arrow-up-right-square-fill" /></button
            ><br />
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from 'vue'
import api, { NewPoll, PollType } from '/@/lib/apis'

interface State {
  title: string
  detail: string
  deadline: string
  options: string[]
  newTags: Set<string>
  newTag: string
  tags: string[]
  autocompletes: string[]
}
export default defineComponent({
  name: 'NewPollComponent',
  setup() {
    const state = reactive<State>({
      title: '',
      detail: '',
      deadline: '0',
      options: ['', ''],
      newTags: new Set(),
      newTag: '',
      tags: [
        'ウマ娘',
        'ウマ娘プリティダービー',
        'スペシャルウィーク',
        'メジロマックイーン',
        'メジロアルダン',
        'メジロライアン',
        'メジロドーベル',
        'ファインモーション',
        'カレンチャン',
        'スマートファルコン',
        'ウマ娘mad'
      ],
      autocompletes: []
    })

    const insertTag = function () {
      if (state.newTag.length === 0) return
      state.newTags.add(state.newTag)
      state.newTag = ''
    }
    const deleteTag = function (str: string) {
      state.newTags.delete(str)
    }
    const insertOption = function () {
      state.options.push('')
    }
    const deleteOption = function (idx: number) {
      state.options.splice(idx, 1)
    }
    const calculateFilter = function () {
      if (state.newTag.length === 0) {
        state.autocompletes = []
      } else {
        state.autocompletes = state.tags
          .filter((v: string) => {
            return v.indexOf(state.newTag) === 0
          })
          .slice(0, 5)
      }
    }
    const onAutocomplete = function (str: string) {
      if (str.length === 0) return
      state.newTags.add(str)
      state.newTag = ''
      calculateFilter()
    }
    const submitPoll = async () => {
      const poll: NewPoll = {
        title: state.title,
        type: PollType.Radio,
        deadline: '2022-01-25T14:01:28.205Z',
        tags: Array.from(state.newTags),
        question: state.options
      }
      await api.postPolls(poll)
      state.title = ''
      state.detail = ''
      state.options = ['', '']
      state.newTags = new Set()
      state.newTag = ''
      calculateFilter()
      state.deadline = '0'
      location.href = '/'
    }

    return {
      state,
      insertTag,
      deleteTag,
      insertOption,
      deleteOption,
      calculateFilter,
      onAutocomplete,
      submitPoll
    }
  }
})
</script>
