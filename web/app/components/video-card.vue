<template>
  <div class="card bg-base-200 shadow-xl hover:shadow-2xl transition-shadow">
    <!-- Thumbnail -->
    <figure class="relative">
      <img
        :src="video.thumbnail || '/placeholder.jpg'"
        :alt="video.title"
        class="w-full aspect-video object-cover"
      />
      <div class="absolute top-2 right-2 badge badge-neutral">
        {{ formatDate(video.published_at) }}
      </div>
    </figure>

    <div class="card-body p-4">
      <!-- Title -->
      <h2
        class="card-title text-base line-clamp-2"
        :title="video.title"
      >
        {{ video.title }}
      </h2>

      <!-- Game Info or Match Button -->
      <div class="mt-2">
        <div v-if="video.game">
          <div class="flex items-center gap-2 mb-2">
            <font-awesome-icon
              icon="gamepad"
              class="w-4 h-4 opacity-70"
            />
            <span class="text-sm font-semibold">{{ video.game.name }}</span>
            <a
              v-if="video.game.url"
              :href="video.game.url"
              target="_blank"
              class="text-xs text-primary hover:underline"
            >
              IGDB →
            </a>
          </div>
          <button
            class="btn btn-outline btn-xs w-full"
            @click="handleChangeGame"
          >
            <font-awesome-icon icon="pen" />
            Change Game
          </button>
        </div>
        <button
          v-else
          class="btn btn-primary btn-sm w-full"
          @click="handleMatchGame"
        >
          <font-awesome-icon icon="gamepad" />
          Match Game
        </button>
      </div>

      <!-- YouTube Link -->
      <div class="card-actions justify-end mt-2">
        <a
          :href="`https://www.youtube.com/watch?v=${video.id}`"
          target="_blank"
          rel="noopener noreferrer"
          class="btn btn-sm btn-ghost"
        >
          <font-awesome-icon :icon="['fab', 'youtube']" />
          Watch
        </a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Video } from '~/composables/useApi'

const props = defineProps<{
  video: Video
}>()

const emit = defineEmits<{
  openGameModal: [video: Video]
}>()

function formatDate(dateString: string) {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

function handleMatchGame() {
  emit('openGameModal', props.video)
}

function handleChangeGame() {
  emit('openGameModal', props.video)
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
