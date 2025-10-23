export interface Offset {
	limit: number
	offset: number
}

export interface Meta extends Offset {
	total: number
}

export interface APIResponse<T> {
	data: T
	meta?: Meta
}
