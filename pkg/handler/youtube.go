package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/K0ng2/zeedzad/config"
	"github.com/K0ng2/zeedzad/model"
	repoModel "github.com/K0ng2/zeedzad/repository/model"
)

// SyncYouTubeVideos godoc
// @Summary Sync videos from YouTube channel
// @Description Fetch and sync videos from OPZTV YouTube channel
// @Tags videos
// @Accept  json
// @Produce  json
// @Param api_key query string true "YouTube API Key"
// @Param max_results query int false "Maximum results to fetch" default(50)
// @Success 200 {object} model.APIResponse[model.SyncResult]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /videos/sync [post]
func (h *Handler) SyncYouTubeVideos(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	maxResults := fiber.Query(c, "max_results", 50)

	// Create YouTube service - use context.Background() like test endpoint
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(config.YOUTUBE_API_KEY))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: "failed to create youtube service: " + err.Error()})
	}

	// Fetch videos from OPZTV channel
	channelID := "UCsGx1qSnAS2P1YCJPYnYVUg"

	// First, get the "uploads" playlist ID for the channel
	channelCall := service.Channels.List([]string{"contentDetails"})
	channelCall = channelCall.Id(channelID)
	channelResponse, err := channelCall.Do()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: "failed to fetch channel: " + err.Error()})
	}

	if len(channelResponse.Items) == 0 {
		return c.Status(http.StatusNotFound).JSON(model.Error{Error: "channel not found"})
	}

	uploadsPlaylistID := channelResponse.Items[0].ContentDetails.RelatedPlaylists.Uploads

	// Insert videos into database
	added := 0
	skipped := 0
	errors := 0
	pageToken := ""
	fetchedCount := 0

	// Paginate through playlist items
	for fetchedCount < maxResults {
		playlistCall := service.PlaylistItems.List([]string{"snippet", "contentDetails"})
		playlistCall = playlistCall.PlaylistId(uploadsPlaylistID)
		playlistCall = playlistCall.MaxResults(50) // Fetch 50 per page (max allowed)

		if pageToken != "" {
			playlistCall = playlistCall.PageToken(pageToken)
		}

		playlistResponse, err := playlistCall.Do()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: "failed to fetch playlist items: " + err.Error()})
		}

		if len(playlistResponse.Items) == 0 {
			break // No more items
		}

		// Process items
		for _, item := range playlistResponse.Items {
			if fetchedCount >= maxResults {
				break
			}
			fetchedCount++

			videoID := item.Id

			// Check if video already exists
			existingVideo, err := h.repo.GetVideoByYouTubeID(ctx, videoID)
			if err == nil && existingVideo != nil {
				skipped++
				continue
			}

			// Parse published date
			publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
			if err != nil {
				publishedAt = time.Now()
			}

			// Prepare video data
			thumbnail := ""
			if item.Snippet.Thumbnails != nil && item.Snippet.Thumbnails.High != nil {
				thumbnail = item.Snippet.Thumbnails.High.Url
			}

			video := repoModel.Videos{
				ID:           stringPtr(uuid.New().String()),
				YoutubeID:    videoID,
				Title:        item.Snippet.Title,
				Description:  stringPtr(item.Snippet.Description),
				Thumbnail:    stringPtr(thumbnail),
				PublishedAt:  publishedAt,
				ChannelID:    item.Snippet.ChannelId,
				ChannelTitle: stringPtr(item.Snippet.ChannelTitle),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			err = h.repo.CreateVideo(ctx, video)
			if err != nil {
				errors++
				continue
			}

			added++
		}

		// Check if there's a next page
		if playlistResponse.NextPageToken == "" {
			break
		}
		pageToken = playlistResponse.NextPageToken
	}

	result := model.SyncResult{
		Added:   added,
		Skipped: skipped,
		Errors:  errors,
		Total:   int(fetchedCount),
	}

	return c.JSON(Response(result, nil))
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
