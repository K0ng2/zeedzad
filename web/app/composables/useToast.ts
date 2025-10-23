export interface Toast {
	id: string
	type: 'success' | 'error' | 'warning' | 'info'
	title?: string
	message: string
	duration?: number // ms, 0 or undefined = no auto-dismiss
	closable?: boolean // default true
	createdAt: number
}

export interface ToastOptions {
	title?: string
	duration?: number // ms, 0 = no auto-dismiss
	closable?: boolean
}

// Global reactive toast store
const toasts = ref<Toast[]>([])
const MAX_TOASTS = 5
let toastIdCounter = 0

const addToast = (options: Omit<Toast, 'id' | 'createdAt'>) => {
	const toast: Toast = {
		id: `toast-${++toastIdCounter}`,
		createdAt: Date.now(),
		closable: true,
		duration: 5000, // 5 seconds default
		...options
	}

	// Ensure message is not empty
	if (!toast.message || toast.message.trim() === '') {
		console.warn('Toast message cannot be empty')
		return ''
	}

	toasts.value.push(toast)

	// Remove oldest toasts if we exceed the limit
	if (toasts.value.length > MAX_TOASTS) {
		toasts.value.splice(0, toasts.value.length - MAX_TOASTS)
	}

	// Auto-dismiss if duration is set
	if (toast.duration && toast.duration > 0) {
		setTimeout(() => {
			removeToast(toast.id)
		}, toast.duration)
	}

	return toast.id
}

const removeToast = (id: string) => {
	const index = toasts.value.findIndex(t => t.id === id)
	if (index > -1) {
		toasts.value.splice(index, 1)
	}
}

const clearAllToasts = () => {
	toasts.value = []
}

// Toast methods
const showSuccess = (message: string, options?: ToastOptions) => {
	return addToast({ type: 'success', message, ...options })
}

const showError = (message: string, options?: ToastOptions) => {
	return addToast({ type: 'error', message, duration: 0, ...options })
}

const showWarning = (message: string, options?: ToastOptions) => {
	return addToast({ type: 'warning', message, ...options })
}

const showInfo = (message: string, options?: ToastOptions) => {
	return addToast({ type: 'info', message, ...options })
}

export interface ToastMethods {
	showSuccess: typeof showSuccess
	showError: typeof showError
	showWarning: typeof showWarning
	showInfo: typeof showInfo
	removeToast: typeof removeToast
	clearAllToasts: typeof clearAllToasts
}

/**
 * Composable for showing toast notifications
 * Uses a global reactive store, so no injection required
 */
export const useToast = (): ToastMethods => {
	return {
		showSuccess,
		showError,
		showWarning,
		showInfo,
		removeToast,
		clearAllToasts
	}
}

/**
 * Get the reactive toasts array for the Toast component
 */
export const useToastStore = () => {
	return {
		toasts: readonly(toasts),
		removeToast,
		clearAllToasts
	}
}

/**
 * Quick success notification
 */
export const toastSuccess = (message: string, options?: ToastOptions) => {
	return showSuccess(message, options)
}

/**
 * Quick error notification
 */
export const toastError = (message: string, options?: ToastOptions) => {
	return showError(message, options)
}

/**
 * Quick warning notification
 */
export const toastWarning = (message: string, options?: ToastOptions) => {
	return showWarning(message, options)
}

/**
 * Quick info notification
 */
export const toastInfo = (message: string, options?: ToastOptions) => {
	return showInfo(message, options)
}