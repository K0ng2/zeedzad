<template>
	<Teleport to="body">
		<div class="toast toast-end z-50 p-4">
			<!-- Clear All button (shown when multiple toasts) -->
			<div
				v-if="toasts && toasts.length > 1"
				class="flex justify-end mb-3"
			>
				<button
					@click="clearAllToasts"
					class="btn btn-xs btn-ghost gap-1 opacity-80 hover:opacity-100 shadow-sm"
				>
					<FontAwesomeIcon
						icon="times"
						class="w-3 h-3"
					/>
					<span>Clear All</span>
				</button>
			</div>

			<TransitionGroup
				name="toast"
				tag="div"
				class="space-y-3"
			>
				<div
					v-for="toast in (toasts || [])"
					:key="toast.id"
					class="alert shadow-2xl max-w-sm border border-base-300 backdrop-blur-sm"
					:class="getToastClasses(toast.type)"
				>
					<div class="flex items-start gap-3 w-full">
						<!-- Icon with modern styling -->
						<div
							class="shrink-0 w-10 h-10 rounded-full flex items-center justify-center"
							:class="getIconBgClass(toast.type)"
						>
							<FontAwesomeIcon
								:icon="getToastIcon(toast.type)"
								class="w-5 h-5"
							/>
						</div>

						<!-- Content -->
						<div class="flex-1 min-w-0">
							<div
								v-if="toast.title"
								class="font-bold text-base mb-1"
							>
								{{ toast.title }}
							</div>
							<div class="text-sm opacity-90">
								{{ toast.message }}
							</div>
						</div>

						<!-- Close button -->
						<button
							v-if="toast.closable !== false"
							@click="removeToast(toast.id)"
							class="btn btn-xs btn-circle btn-ghost shrink-0 opacity-70 hover:opacity-100"
						>
							<FontAwesomeIcon
								icon="times"
								class="w-4 h-4"
							/>
						</button>
					</div>

					<!-- Modern progress bar for auto-dismiss -->
					<div
						v-if="toast.duration && toast.duration > 0"
						class="absolute bottom-0 left-0 right-0 h-1 bg-current opacity-25 overflow-hidden"
					>
						<div
							class="h-full bg-current transition-all ease-linear"
							:style="{
								width: getProgressWidth(toast),
								transitionDuration: '100ms'
							}"
						></div>
					</div>
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

const getIconBgClass = (type: Toast['type']) => {
	switch (type) {
		case 'success': return 'bg-success/20 text-success'
		case 'error': return 'bg-error/20 text-error'
		case 'warning': return 'bg-warning/20 text-warning'
		case 'info': return 'bg-info/20 text-info'
		default: return 'bg-info/20 text-info'
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
/* Enhanced Toast transition animations */
.toast-enter-active {
	transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.toast-leave-active {
	transition: all 0.3s cubic-bezier(0.4, 0, 1, 1);
}

.toast-enter-from {
	opacity: 0;
	transform: translateX(120%) scale(0.9);
}

.toast-leave-to {
	opacity: 0;
	transform: translateX(120%) scale(0.95);
}

.toast-move {
	transition: transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Progress bar styling */
.alert {
	position: relative;
	overflow: hidden;
}
</style>
