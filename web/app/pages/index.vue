<template>
	<div class="min-h-screen">
		<!-- Hero Section with Gradient -->
		<div class="gradient-hero border-b border-base-300">
			<div class="container mx-auto px-4 py-12 md:py-16">
				<div class="max-w-4xl mx-auto text-center mb-8">
					<h1
						class="text-5xl md:text-6xl font-bold mb-4 bg-linear-to-r from-primary to-secondary bg-clip-text text-transparent"
					>
						OPZTV Videos
					</h1>
					<p class="text-lg md:text-xl text-base-content/70 mb-8">
						Discover and match YouTube videos with games from IGDB
					</p>

					<!-- Modern Search Bar -->
					<div class="relative max-w-2xl mx-auto">
						<div class="relative flex items-center gap-2">
							<div class="relative flex-1">
								<font-awesome-icon
									icon="search"
									class="absolute left-4 top-1/2 -translate-y-1/2 text-base-content/40 w-5 h-5"
								/>
								<input
									v-model="searchQuery"
									type="text"
									placeholder="Search by video title or game name..."
									class="input input-lg input-bordered w-full pl-12 pr-4 shadow-lg focus:shadow-xl transition-shadow"
									@keyup.enter="search"
								/>
							</div>
							<button
								class="btn btn-primary btn-lg px-8 shadow-lg hover:shadow-xl"
								@click="search"
							>
								<font-awesome-icon
									icon="search"
									class="md:mr-2"
								/>
								<span class="hidden md:inline">Search</span>
							</button>
						</div>
					</div>

					<!-- Stats Bar -->
					<div class="flex justify-center gap-6 mt-8 text-sm">
						<div class="flex items-center gap-2">
							<div class="badge badge-primary badge-lg">
								<font-awesome-icon
									icon="video"
									class="mr-2"
								/>
								{{ meta?.total || 0 }} Videos
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Main Content -->
		<div class="container mx-auto px-4 py-8">

			<!-- Loading State -->
			<div
				v-if="loading"
				class="flex flex-col items-center justify-center py-24"
			>
				<span class="loading loading-spinner loading-lg text-primary mb-4"></span>
				<p class="text-base-content/60">Loading videos...</p>
			</div>

			<!-- Error State -->
			<div
				v-else-if="error"
				class="alert alert-error shadow-lg max-w-2xl mx-auto"
			>
				<font-awesome-icon
					icon="exclamation-circle"
					class="w-6 h-6"
				/>
				<span>{{ error }}</span>
			</div>

			<!-- Video Grid -->
			<div v-else-if="videos.length > 0">
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-6 mb-12">
					<video-card
						v-for="video in videos"
						:key="video.id"
						:video="video"
						@open-game-modal="openGameModal"
						@game-removed="refreshVideos"
					/>
				</div>

				<!-- Modern Pagination -->
				<div class="flex flex-col sm:flex-row justify-center items-center gap-4 py-8 border-t border-base-300">
					<button
						class="btn btn-outline btn-sm sm:btn-md gap-2"
						:disabled="currentPage === 1"
						@click="goToPage(currentPage - 1)"
					>
						<font-awesome-icon
							icon="chevron-left"
							class="w-4 h-4"
						/>
						<span class="hidden sm:inline">Previous</span>
					</button>

					<div class="flex items-center gap-3">
						<div class="badge badge-lg badge-ghost">
							Page {{ currentPage }} of {{ totalPages }}
						</div>
						<div class="badge badge-primary badge-outline">
							{{ meta?.total || 0 }} videos
						</div>
					</div>

					<button
						class="btn btn-outline btn-sm sm:btn-md gap-2"
						:disabled="currentPage >= totalPages"
						@click="goToPage(currentPage + 1)"
					>
						<span class="hidden sm:inline">Next</span>
						<font-awesome-icon
							icon="chevron-right"
							class="w-4 h-4"
						/>
					</button>
				</div>
			</div>

			<!-- Empty State -->
			<div
				v-else
				class="flex flex-col items-center justify-center py-24"
			>
				<div class="bg-base-200 rounded-full p-8 mb-6">
					<font-awesome-icon
						icon="video"
						class="w-16 h-16 text-base-content/30"
					/>
				</div>
				<h3 class="text-2xl font-semibold mb-2">No videos found</h3>
				<p class="text-base-content/60">Try adjusting your search criteria</p>
			</div>
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

const ITEMS_PER_PAGE = 25

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
