<template>
	<div class="card bg-base-100 shadow-lg card-hover overflow-hidden border border-base-300">
		<!-- Thumbnail with Overlay -->
		<figure class="relative group overflow-hidden">
			<img
				:src="video.thumbnail || '/placeholder.jpg'"
				:alt="video.title"
				class="w-full aspect-video object-cover transition-transform duration-500 group-hover:scale-110"
			/>
			<!-- Gradient Overlay on Hover -->
			<div
				class="absolute inset-0 bg-linear-to-t from-base-300/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
			></div>

			<!-- Date Badge -->
			<div class="absolute top-3 right-3 badge badge-neutral badge-sm shadow-lg">
				<font-awesome-icon
					icon="calendar"
					class="w-3 h-3 mr-1"
				/>
				{{ formatDate(video.published_at) }}
			</div>

			<!-- Watch Button Overlay -->
			<a
				:href="`https://www.youtube.com/watch?v=${video.id}`"
				target="_blank"
				rel="noopener noreferrer"
				class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300"
			>
				<div class="btn btn-circle btn-lg btn-primary shadow-2xl">
					<font-awesome-icon
						:icon="['fab', 'youtube']"
						class="w-6 h-6"
					/>
				</div>
			</a>
		</figure>

		<div class="card-body p-5 space-y-4">
			<!-- Title -->
			<h2
				class="card-title text-base font-semibold line-clamp-2 leading-snug hover:text-primary transition-colors cursor-default"
				:title="video.title"
			>
				{{ video.title }}
			</h2>

			<!-- Game Info or Match Button -->
			<div class="space-y-3">
				<div
					v-if="video.game"
					class="space-y-3"
				>
					<!-- Game Badge -->
					<div class="flex items-center gap-2 p-3 bg-primary/5 rounded-lg border border-primary/20">
						<div class="shrink-0 w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center">
							<font-awesome-icon
								icon="gamepad"
								class="w-4 h-4 text-primary"
							/>
						</div>
						<div class="flex-1 min-w-0">
							<p class="text-sm font-semibold text-base-content line-clamp-1">
								{{ video.game.name }}
							</p>
							<a
								v-if="video.game.url"
								:href="video.game.url"
								target="_blank"
								class="text-xs text-primary hover:underline inline-flex items-center gap-1"
								@click.stop
							>
								View on IGDB
								<font-awesome-icon
									icon="external-link-alt"
									class="w-3 h-3"
								/>
							</a>
						</div>
					</div>

					<!-- Action Buttons -->
					<div class="flex gap-2">
						<button
							class="btn btn-outline btn-sm flex-1 gap-2"
							@click="handleChangeGame"
							:disabled="isDeleting"
						>
							<font-awesome-icon
								icon="pen"
								class="w-3 h-3"
							/>
							<span>Change</span>
						</button>
						<button
							class="btn btn-error btn-outline btn-sm flex-1 gap-2"
							@click="handleRemoveGame"
							:disabled="isDeleting"
						>
							<font-awesome-icon
								v-if="isDeleting"
								icon="spinner"
								spin
								class="w-3 h-3"
							/>
							<font-awesome-icon
								v-else
								icon="trash"
								class="w-3 h-3"
							/>
							<span>Remove</span>
						</button>
					</div>
				</div>

				<!-- Match Game Button -->
				<button
					v-else
					class="btn btn-primary btn-md w-full gap-2 shadow-md hover:shadow-lg"
					@click="handleMatchGame"
				>
					<font-awesome-icon
						icon="gamepad"
						class="w-4 h-4"
					/>
					<span class="font-semibold">Match Game</span>
				</button>
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
	gameRemoved: []
}>()

const api = useApi()
const toast = useToast()
const isDeleting = ref(false)

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

async function handleRemoveGame() {
	if (!confirm('Are you sure you want to remove this game match?')) {
		return
	}

	isDeleting.value = true

	try {
		await api.deleteVideoGame(props.video.id)
		toast.showSuccess('Game removed successfully!')
		emit('gameRemoved')
	} catch (e: any) {
		toast.showError(e.message || 'Failed to remove game')
	} finally {
		isDeleting.value = false
	}
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
