<template>
	<div class="container mx-auto px-4 py-8">
		<div class="mb-8">
			<h1 class="text-4xl font-bold mb-4">OPZTV Videos</h1>
			<p class="text-lg opacity-70 mb-6">Match YouTube videos with Steam games</p>

			<!-- Search Bar -->
			<div class="form-control">
				<div class="input-group">
					<input
						v-model="searchQuery"
						type="text"
						placeholder="Search by video title or game name..."
						class="input input-bordered w-full"
						@keyup.enter="search"
					/>
					<button
						class="btn btn-primary"
						@click="search"
					>
						<font-awesome-icon icon="search" />
					</button>
				</div>
			</div>
		</div>

		<!-- Loading State -->
		<div
			v-if="loading"
			class="flex justify-center py-16"
		>
			<span class="loading loading-spinner loading-lg"></span>
		</div>

		<!-- Error State -->
		<div
			v-else-if="error"
			class="alert alert-error"
		>
			<font-awesome-icon icon="exclamation-circle" />
			<span>{{ error }}</span>
		</div>

		<!-- Video Grid -->
		<div v-else-if="videos.length > 0">
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mb-8">
				<div
					v-for="video in videos"
					:key="video.id"
					class="card bg-base-200 shadow-xl hover:shadow-2xl transition-shadow"
				>
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
							<div
								v-if="video.game"
								class="flex items-center gap-2"
							>
								<img
									v-if="video.game.icon"
									:src="video.game.icon"
									:alt="video.game.name"
									class="w-8 h-8 rounded"
								/>
								<span class="text-sm font-semibold">{{ video.game.name }}</span>
							</div>
							<button
								v-else
								class="btn btn-primary btn-sm w-full"
								@click="openGameModal(video)"
							>
								<font-awesome-icon icon="gamepad" />
								Match Game
							</button>
						</div>

						<!-- YouTube Link -->
						<div class="card-actions justify-end mt-2">
							<a
								:href="`https://www.youtube.com/watch?v=${video.youtube_id}`"
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
			</div>

			<!-- Pagination -->
			<div class="flex justify-center items-center gap-4">
				<button
					class="btn btn-outline"
					:disabled="currentPage === 1"
					@click="goToPage(currentPage - 1)"
				>
					<font-awesome-icon icon="chevron-left" />
					Previous
				</button>

				<div class="flex items-center gap-2">
					<span class="text-sm">
						Page {{ currentPage }} of {{ totalPages }}
					</span>
					<span class="text-sm opacity-70">
						({{ meta?.total || 0 }} total videos)
					</span>
				</div>

				<button
					class="btn btn-outline"
					:disabled="currentPage >= totalPages"
					@click="goToPage(currentPage + 1)"
				>
					Next
					<font-awesome-icon icon="chevron-right" />
				</button>
			</div>
		</div>

		<!-- Empty State -->
		<div
			v-else
			class="py-16 text-center opacity-70"
		>
			<font-awesome-icon
				icon="video"
				size="3x"
				class="mb-4"
			/>
			<p class="text-xl">No videos found</p>
		</div>

		<!-- Game Search Modal -->
		<game-search-modal
			v-if="selectedVideo"
			ref="gameModal"
			:video-id="selectedVideo.id"
			:video-title="selectedVideo.title"
			@game-matched="refreshVideos"
		/>
	</div>
</template>

<script setup lang="ts">
import type { Video, Meta } from '~/composables/useApi'

definePageMeta({
	layout: 'default',
})

const ITEMS_PER_PAGE = 24

const api = useApi()
const videos = ref<Video[]>([])
const meta = ref<Meta | null>(null)
const loading = ref(false)
const error = ref('')
const searchQuery = ref('')
const currentPage = ref(1)
const selectedVideo = ref<Video | null>(null)
const gameModal = ref<any>()

const totalPages = computed(() => {
	if (!meta.value?.total) return 0
	return Math.ceil(meta.value.total / ITEMS_PER_PAGE)
})

async function fetchVideos() {
	loading.value = true
	error.value = ''

	try {
		const offset = (currentPage.value - 1) * ITEMS_PER_PAGE
		const response = await api.getVideos({
			offset,
			limit: ITEMS_PER_PAGE,
			search: searchQuery.value || undefined,
		})

		videos.value = response.data
		meta.value = response.meta || null
	} catch (e: any) {
		error.value = e.message || 'Failed to fetch videos'
		videos.value = []
	} finally {
		loading.value = false
	}
}

function search() {
	currentPage.value = 1
	fetchVideos()
}

function goToPage(page: number) {
	currentPage.value = page
	fetchVideos()
	window.scrollTo({ top: 0, behavior: 'smooth' })
}

function formatDate(dateString: string) {
	const date = new Date(dateString)
	return date.toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric',
	})
}

function openGameModal(video: Video) {
	selectedVideo.value = video
	nextTick(() => {
		gameModal.value?.open()
	})
}

function refreshVideos() {
	fetchVideos()
}

// Load videos on mount
onMounted(() => {
	fetchVideos()
})
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
