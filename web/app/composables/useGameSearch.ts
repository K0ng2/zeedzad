import type { IGDBGameSearchResult, GameResponse } from './useApi'

export function useGameSearch() {
  const api = useApi()

  const searchQuery = ref('')
  const searchResults = ref<IGDBGameSearchResult[]>([])
  const dbResults = ref<GameResponse[]>([])
  const showingDbResults = ref(false)
  const loading = ref(false)
  const error = ref('')

  // Combined results for display
  const displayResults = computed(() => {
    if (showingDbResults.value) {
      // Convert database results to match IGDB format
      return dbResults.value.map(game => ({
        id: game.id,
        name: game.name,
        url: game.url || '',
      }))
    }
    return searchResults.value
  })

  async function searchGames(query: string) {
    if (!query.trim()) return

    loading.value = true
    error.value = ''
    searchResults.value = []
    dbResults.value = []
    showingDbResults.value = false

    try {
      // First, search in local database
      const dbResponse = await api.getGames({
        search: query,
        limit: 50,
      })

      if (dbResponse.data && dbResponse.data.length > 0) {
        // Found results in database
        dbResults.value = dbResponse.data
        showingDbResults.value = true
      } else {
        // No results in database, search IGDB API
        const igdbResponse = await api.searchIGDBGames(query)
        searchResults.value = igdbResponse.data
        showingDbResults.value = false
      }
    } catch (e: any) {
      error.value = e.message || 'Failed to search games'
      searchResults.value = []
      dbResults.value = []
    } finally {
      loading.value = false
    }
  }

  async function matchGameToVideo(
    game: IGDBGameSearchResult,
    videoId: string,
    isFromDatabase: boolean
  ) {
    loading.value = true
    error.value = ''

    try {
      let gameId = game.id

      // If selecting from IGDB results (not from DB), create the game first
      if (!isFromDatabase) {
        const gameResponse = await api.createGame({
          id: game.id,
          name: game.name,
          url: game.url,
        })
        gameId = gameResponse.data.id
      }

      // Update the video with this game
      await api.updateVideoGame(videoId, gameId)

      return { success: true }
    } catch (e: any) {
      error.value = e.message || 'Failed to match game'
      throw e
    } finally {
      loading.value = false
    }
  }

  function reset() {
    searchQuery.value = ''
    searchResults.value = []
    dbResults.value = []
    showingDbResults.value = false
    error.value = ''
  }

  return {
    searchQuery,
    searchResults,
    dbResults,
    showingDbResults,
    loading,
    error,
    displayResults,
    searchGames,
    matchGameToVideo,
    reset,
  }
}
