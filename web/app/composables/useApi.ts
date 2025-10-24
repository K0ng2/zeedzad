export interface Video {
	id: string
	youtube_id: string
	title: string
	description?: string
	thumbnail?: string
	published_at: string
	channel_id: string
	channel_title?: string
	game?: Game
	created_at: string
	updated_at: string
}

export interface Game {
	id: number
	name: string
	url?: string
}

export interface GameResponse extends Game {
	created_at: string
	updated_at: string
}

export interface IGDBGameSearchResult {
	id: number
	name: string
	url: string
}

export interface Meta {
	total: number
	limit: number
	offset: number
}

export interface APIResponse<T> {
	data: T
	meta?: Meta
}

export function useApi() {
	const config = useRuntimeConfig()
	const baseURL = config.public.apiBase || '/api'

	async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
		const response = await fetch(`${baseURL}${endpoint}`, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				...options?.headers,
			},
		})

		if (!response.ok) {
			const error = await response.json().catch(() => ({ error: response.statusText }))
			throw new Error(error.error || 'API request failed')
		}

		return response.json()
	}

	return {
		// Video endpoints
		async getVideos(params: { offset?: number; limit?: number; search?: string } = {}) {
			const query = new URLSearchParams()
			if (params.offset) query.append('offset', params.offset.toString())
			if (params.limit) query.append('limit', params.limit.toString())
			if (params.search) query.append('search', params.search)

			return fetchAPI<APIResponse<Video[]>>(`/videos?${query}`)
		},

		async getVideoById(id: string) {
			return fetchAPI<APIResponse<Video>>(`/videos/${id}`)
		},

		async updateVideoGame(videoId: string, gameId: string) {
			return fetchAPI<void>(`/videos/${videoId}/game`, {
				method: 'PUT',
				body: JSON.stringify({ game_id: gameId }),
			})
		},

		// Game endpoints
		async getGames(params: { offset?: number; limit?: number; search?: string } = {}) {
			const query = new URLSearchParams()
			if (params.offset) query.append('offset', params.offset.toString())
			if (params.limit) query.append('limit', params.limit.toString())
			if (params.search) query.append('search', params.search)

			return fetchAPI<APIResponse<GameResponse[]>>(`/games?${query}`)
		},

		async getGameById(id: string) {
			return fetchAPI<APIResponse<GameResponse>>(`/games/${id}`)
		},

		async searchIGDBGames(query: string) {
			const params = new URLSearchParams({ q: query })
			return fetchAPI<APIResponse<IGDBGameSearchResult[]>>(`/games/igdb/search?${params}`)
		},

		async createGame(game: { id: number; name: string; url?: string }) {
			return fetchAPI<APIResponse<GameResponse>>('/games', {
				method: 'POST',
				body: JSON.stringify(game),
			})
		},
	}
}
