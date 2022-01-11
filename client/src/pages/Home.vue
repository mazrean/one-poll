<template>
  <div class="d-flex flex-wrap">
    <PollCardComponent
      v-for="card in cards"
      :key="card.pollId"
      :title="card.title"
      :tags="card.tags"
      :question="card.question"
      :choices="card.choices"
      :created-at="card.createdAt"
      :owner="card.owner">
    </PollCardComponent>
  </div>
</template>

<script lang="ts">
import PollCardComponent from '/@/components/PollCard.vue'

export default {
  name: 'HomePage',
  components: { PollCardComponent },
  setup() {
    const cards = [
      {
        pollId: 'poll_id 1',
        title: 'Title 1',
        type: 'type 1',
        deadline: 'deadline 1',
        tags: [
          {
            id: 1,
            name: 'tag 1 1'
          },
          {
            id: 2,
            name: 'tag 1 2'
          },
          {
            id: 3,
            name: 'tag 1 3'
          }
        ],
        question: [
          {
            id: 1,
            choice: 'choice 1 1'
          },
          {
            id: 2,
            choice: 'choice 1 2'
          },
          {
            id: 3,
            choice: 'choice 1 3'
          },
          {
            id: 4,
            choice: 'choice 1 4'
          }
        ],
        createdAt: 'created_at 1',
        owner: 'owner 1'
      },
      {
        pollId: 'poll_id 2',
        title: 'Title 2',
        type: 'type 2',
        deadline: 'deadline 2',
        tags: [
          {
            id: 1,
            name: 'tag 2 1'
          },
          {
            id: 2,
            name: 'tag 2 2'
          },
          {
            id: 3,
            name: 'tag 2 3'
          }
        ],
        question: [
          {
            id: 1,
            choice: 'choice 2 1'
          },
          {
            id: 2,
            choice: 'choice 2 2'
          },
          {
            id: 3,
            choice: 'choice 2 3'
          },
          {
            id: 4,
            choice: 'choice 2 4'
          }
        ],
        createdAt: 'created_at 2',
        owner: 'owner 2'
      },
      {
        pollId: 'poll_id 3',
        title: 'Title 3',
        type: 'type 3',
        deadline: 'deadline 3',
        tags: [
          {
            id: 1,
            name: 'tag 3 1'
          },
          {
            id: 2,
            name: 'tag 3 2'
          },
          {
            id: 3,
            name: 'tag 3 3'
          }
        ],
        question: [
          {
            id: 1,
            choice: 'choice 3 1'
          },
          {
            id: 2,
            choice: 'choice 3 2'
          },
          {
            id: 3,
            choice: 'choice 3 3'
          }
        ],
        createdAt: 'created_at 3',
        owner: 'owner 3'
      }
    ]
    return {
      cards,
      PollCardComponent
    }
  }
}
</script>
