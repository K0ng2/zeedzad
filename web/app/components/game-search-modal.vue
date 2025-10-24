<template>
	<dialog
		ref="modal"
		class="modal modal-bottom sm:modal-middle"
	>
		<div class="modal-box max-w-3xl">
			<!-- Header -->
			<div class="flex items-start justify-between mb-6">
				<div>
					<h3 class="text-2xl font-bold mb-2">Match Game</h3>
					<p class="text-sm text-base-content/70">
						Search and match a game for:
					</p>
					<p class="text-sm font-semibold text-primary mt-1 line-clamp-1">
						{{ videoTitle }}
					</p>
				</div>
				<button
					@click="close"
					class="btn btn-sm btn-circle btn-ghost"
				>
					<font-awesome-icon
						icon="times"
						class="w-4 h-4"
					/>
				</button>
			</div>

			<!-- Modern Search Input -->
			<div class="mb-6">
				<div class="relative">
					<font-awesome-icon
						icon="search"
						class="absolute left-4 top-1/2 -translate-y-1/2 text-base-content/40 w-5 h-5"
					/>
					<input
						v-model="searchQuery"
						type="text"
						placeholder="Search IGDB games..."
						class="input input-bordered w-full pl-12 pr-28 input-lg focus:outline-primary"
						@keyup.enter="searchGames"
					/>
					<button
						class="btn btn-primary absolute right-1 top-1/2 -translate-y-1/2"
						@click="searchGames"
						:disabled="loading"
					>
						<font-awesome-icon
							v-if="loading"
							icon="spinner"
							spin
							class="w-4 h-4"
						/>
						<font-awesome-icon
							v-else
							icon="search"
							class="w-4 h-4"
						/>
						<span class="hidden sm:inline ml-2">Search</span>
					</button>
				</div>
			</div>

			<!-- Loading State -->
			<div
				v-if="loading"
				class="flex flex-col items-center justify-center py-12"
			>
				<span class="loading loading-spinner loading-lg text-primary mb-3"></span>
				<p class="text-sm text-base-content/60">Searching games...</p>
			</div>

			<!-- Error State -->
			<div
				v-else-if="error"
				class="alert alert-error shadow-lg"
			>
				<font-awesome-icon
					icon="exclamation-circle"
					class="w-5 h-5"
				/>
				<span>{{ error }}</span>
			</div>

			<!-- Search Results -->
			<div
				v-else-if="searchResults.length > 0"
				class="space-y-3 max-h-[60vh] overflow-y-auto pr-2 scrollbar-thin"
			>
				<div
					v-for="game in searchResults"
					:key="game.id"
					class="group relative bg-base-200 hover:bg-base-300 rounded-xl border border-base-300 hover:border-primary/50 cursor-pointer transition-all duration-200 hover:shadow-lg overflow-hidden"
					@click="selectGame(game)"
				>
					<div class="flex items-center gap-4 p-4">
						<div
							class="shrink-0 w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors"
						>
							<font-awesome-icon
								icon="gamepad"
								class="w-6 h-6 text-primary"
							/>
						</div>

						<div class="flex-1 min-w-0">
							<h4 class="font-semibold text-base mb-1 line-clamp-1 group-hover:text-primary transition-colors">
								{{ game.name }}
							</h4>
							<div class="flex items-center gap-2 text-xs text-base-content/60">
								<span class="badge badge-ghost badge-sm">ID: {{ game.id }}</span>
								<a
									v-if="game.url"
									:href="game.url"
									target="_blank"
									class="inline-flex items-center gap-1 text-primary hover:underline"
									@click.stop
								>
									<span>View on IGDB</span>
									<font-awesome-icon
										icon="external-link-alt"
										class="w-3 h-3"
									/>
								</a>
							</div>
						</div>

						<button class="btn btn-primary btn-sm gap-2 shrink-0">
							<font-awesome-icon
								icon="check"
								class="w-4 h-4"
							/>
							<span class="hidden sm:inline">Select</span>
						</button>
					</div>
				</div>
			</div>

			<!-- Empty State -->
			<div
				v-else-if="searchQuery && !loading"
				class="flex flex-col items-center justify-center py-12 text-center"
			>
				<div class="bg-base-200 rounded-full p-6 mb-4">
					<font-awesome-icon
						icon="search"
						class="w-12 h-12 text-base-content/30"
					/>
				</div>
				<h4 class="font-semibold text-lg mb-2">No games found</h4>
				<p class="text-sm text-base-content/60">Try a different search term</p>
			</div>

			<!-- Initial State -->
			<div
				v-else
				class="flex flex-col items-center justify-center py-12 text-center"
			>
				<div class="bg-primary/10 rounded-full p-6 mb-4">
					<font-awesome-icon
						icon="gamepad"
						class="w-12 h-12 text-primary"
					/>
				</div>
				<h4 class="font-semibold text-lg mb-2">Search for a game</h4>
				<p class="text-sm text-base-content/60">Enter a game name to find matches from IGDB</p>
			</div>
		</div>

		<!-- Modal Backdrop -->
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
		await api.updateVideoGame(props.videoId, gameResponse.data.id)

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
