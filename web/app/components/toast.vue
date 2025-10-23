<template>
	<Teleport to="body">
		<div class="toast z-50">
			<!-- Clear All button (shown when multiple toasts) -->
			<div
				v-if="toasts && toasts.length > 1"
				class="flex justify-end mb-2"
			>
				<button
					@click="clearAllToasts"
					class="btn btn-xs btn-ghost opacity-70 hover:opacity-100 text-xs"
				>
					<FontAwesomeIcon
						icon="times"
						class="w-3 h-3 mr-1"
					/>
					Clear All
				</button>
			</div>

			<TransitionGroup
				name="toast"
				tag="div"
				class="space-y-2"
			>
				<div
					v-for="toast in (toasts || [])"
					:key="toast.id"
					class="alert shadow-lg max-w-sm"
					:class="getToastClasses(toast.type)"
				>
					<div class="flex items-center gap-3">
						<!-- Icon -->
						<FontAwesomeIcon
							:icon="getToastIcon(toast.type)"
							class="w-6 h-6 shrink-0"
						/>

						<!-- Content -->
						<div class="flex-1 min-w-0">
							<div
								v-if="toast.title"
								class="font-semibold text-sm"
							>
								{{ toast.title }}
							</div>
							<div
								class="text-sm"
								:class="{ 'mt-1': toast.title }"
							>
								{{ toast.message }}
							</div>
						</div>

						<!-- Close button -->
						<button
							v-if="toast.closable !== false"
							@click="removeToast(toast.id)"
							class="btn btn-sm btn-circle btn-ghost opacity-70 hover:opacity-100"
						>
							<FontAwesomeIcon
								icon="times"
								class="w-4 h-4"
							/>
						</button>
					</div>

					<!-- Progress bar for auto-dismiss -->
					<div
						v-if="toast.duration && toast.duration > 0"
						class="absolute bottom-0 left-0 h-1 bg-current opacity-30 transition-all ease-linear"
						:style="{
							width: getProgressWidth(toast),
							transitionDuration: '100ms'
						}"
					></div>
				</div>
			</TransitionGroup>
		</div>
	</Teleport>
</template>

<script setup lang="ts">
import type { Toast } from '~/composables/useToast'

const { toasts, removeToast, clearAllToasts } = useToastStore()

// Reactive time for progress calculation
const currentTime = ref(Date.now())
let updateInterval: NodeJS.Timeout | null = null

// Update time every 100ms for smooth progress animation (client-side only)
onMounted(() => {
	updateInterval = setInterval(() => {
		currentTime.value = Date.now()
	}, 100)
})

// Cleanup interval on component unmount
onUnmounted(() => {
	if (updateInterval) {
		clearInterval(updateInterval)
	}
})

const getToastIcon = (type: Toast['type']) => {
	switch (type) {
		case 'success': return 'check-circle'
		case 'error': return 'exclamation-circle'
		case 'warning': return 'exclamation-triangle'
		case 'info': return 'info-circle'
		default: return 'info-circle'
	}
}

const getToastClasses = (type: Toast['type']) => {
	switch (type) {
		case 'success': return 'alert-success'
		case 'error': return 'alert-error'
		case 'warning': return 'alert-warning'
		case 'info': return 'alert-info'
		default: return 'alert-info'
	}
}

const getProgressWidth = (toast: Toast) => {
	if (!toast.duration || toast.duration <= 0) return '0%'

	const elapsed = currentTime.value - toast.createdAt
	const remaining = Math.max(0, toast.duration - elapsed)
	const progress = (remaining / toast.duration) * 100

	return `${Math.max(0, progress)}%`
}
</script>
<style scoped>
/* Toast transition animations */
.toast-enter-active,
.toast-leave-active {
	transition: all 0.3s ease;
}

.toast-enter-from {
	opacity: 0;
	transform: translateX(100%);
}

.toast-leave-to {
	opacity: 0;
	transform: translateX(100%);
}

.toast-move {
	transition: transform 0.3s ease;
}

/* Progress bar animation */
.alert {
	position: relative;
	overflow: hidden;
}
</style>
