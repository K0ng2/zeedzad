<template>
	<dialog
		ref="modal"
		class="modal"
	>
		<div class="modal-box max-w-2xl">
			<h3 class="text-lg font-bold">Match Game</h3>
			<p class="py-2 text-sm opacity-70">Search and match a game for: {{ videoTitle }}</p>

			<!-- Search Input -->
			<div class="form-control mb-4">
				<div class="input-group">
					<input
						v-model="searchQuery"
						type="text"
						placeholder="Search IGDB games..."
						class="input input-bordered w-full"
						@keyup.enter="searchGames"
					/>
					<button
						class="btn btn-primary"
						@click="searchGames"
						:disabled="loading"
					>
						<font-awesome-icon
							v-if="loading"
							icon="spinner"
							spin
						/>
						<font-awesome-icon
							v-else
							icon="search"
						/>
					</button>
				</div>
			</div>

			<!-- Loading State -->
			<div
				v-if="loading"
				class="flex justify-center py-8"
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

			<!-- Search Results -->
			<div
				v-else-if="searchResults.length > 0"
				class="space-y-2 max-h-96 overflow-y-auto"
			>
				<div
					v-for="game in searchResults"
					:key="game.id"
					class="card bg-base-200 hover:bg-base-300 cursor-pointer transition-colors"
					@click="selectGame(game)"
				>
					<div class="card-body p-4 flex-row items-center gap-4">
						<div class="flex-1">
							<h4 class="card-title text-base">{{ game.name }}</h4>
							<p class="text-sm opacity-70">ID: {{ game.id }}</p>
							<a
								v-if="game.url"
								:href="game.url"
								target="_blank"
								class="text-xs text-primary hover:underline"
								@click.stop
							>
								View on IGDB →
							</a>
						</div>
						<button class="btn btn-sm btn-primary">
							Select
						</button>
					</div>
				</div>
			</div>

			<!-- Empty State -->
			<div
				v-else-if="searchQuery && !loading"
				class="py-8 text-center opacity-70"
			>
				<font-awesome-icon
					icon="search"
					size="2x"
					class="mb-2"
				/>
				<p>No games found. Try a different search term.</p>
			</div>

			<!-- Initial State -->
			<div
				v-else
				class="py-8 text-center opacity-70"
			>
				<font-awesome-icon
					icon="gamepad"
					size="2x"
					class="mb-2"
				/>
				<p>Search for a game to match with this video</p>
			</div>

			<!-- Modal Actions -->
			<div class="modal-action">
				<button
					class="btn"
					@click="close"
				>Cancel</button>
			</div>
		</div>
		<form
			method="dialog"
			class="modal-backdrop"
		>
			<button @click="close">close</button>
		</form>
	</dialog>
</template>

<script setup lang="ts">
import type { IGDBGameSearchResult } from '~/composables/useApi'

const props = defineProps<{
	videoId: string
	videoTitle: string
}>()

const emit = defineEmits<{
	gameMatched: []
}>()

const modal = ref<HTMLDialogElement>()
const searchQuery = ref('')
const searchResults = ref<IGDBGameSearchResult[]>([])
const loading = ref(false)
const error = ref('')

const api = useApi()
const toast = useToast()

function open() {
	modal.value?.showModal()
	searchQuery.value = props.videoTitle
	searchGames()
}

function close() {
	modal.value?.close()
	searchQuery.value = ''
	searchResults.value = []
	error.value = ''
}

async function searchGames() {
	if (!searchQuery.value.trim()) return

	loading.value = true
	error.value = ''

	try {
		const response = await api.searchIGDBGames(searchQuery.value)
		searchResults.value = response.data
	} catch (e: any) {
		error.value = e.message || 'Failed to search games'
		searchResults.value = []
	} finally {
		loading.value = false
	}
}

async function selectGame(igdbGame: IGDBGameSearchResult) {
	loading.value = true
	error.value = ''

	try {
		// First, create/get the game in our database
		const gameResponse = await api.createGame({
			id: igdbGame.id,
			name: igdbGame.name,
			url: igdbGame.url,
		})

		// Then update the video with this game
		await api.updateVideoGame(props.videoId, gameResponse.data.id.toString())

		toast.showSuccess('Game matched successfully!')
		emit('gameMatched')
		close()
	} catch (e: any) {
		error.value = e.message || 'Failed to match game'
		toast.showError('Failed to match game')
	} finally {
		loading.value = false
	}
}

defineExpose({ open })
</script>
