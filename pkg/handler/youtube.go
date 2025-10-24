package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/K0ng2/zeedzad/config"
	"github.com/K0ng2/zeedzad/model"
	repoModel "github.com/K0ng2/zeedzad/repository/model"
)

const (
	opztvChannelID     = "UCsGx1qSnAS2P1YCJPYnYVUg"
	defaultMaxResults  = 50
	youtubeMaxPageSize = 50
	syncLogPrefix      = "[YouTube Sync]"
)

type videoSyncStats struct {
	added        int
	skipped      int
	errors       int
	totalFetched int
}

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
	maxResults := fiber.Query(c, "max_results", defaultMaxResults)

	stats, err := h.syncVideosFromChannel(c.RequestCtx(), maxResults)
	if err != nil {
		return c.Status(err.statusCode).JSON(model.Error{Error: err.message})
	}

	result := model.SyncResult{
		Added:   stats.added,
		Skipped: stats.skipped,
		Errors:  stats.errors,
		Total:   stats.totalFetched,
	}

	return c.JSON(Response(result, nil))
}

func (h *Handler) SyncYouTubeVideosScheduled() {
	stats, err := h.syncVideosFromChannel(context.Background(), defaultMaxResults)
	if err != nil {
		fmt.Printf("%s failed: %s\n", syncLogPrefix, err.message)
		return
	}

	fmt.Printf("%s completed - Added: %d, Skipped: %d, Errors: %d, Total: %d\n",
		syncLogPrefix, stats.added, stats.skipped, stats.errors, stats.totalFetched)
}

type syncError struct {
	message    string
	statusCode int
}

func (e *syncError) Error() string {
	return e.message
}

func (h *Handler) syncVideosFromChannel(ctx context.Context, maxResults int) (*videoSyncStats, *syncError) {
	service, err := h.createYouTubeService()
	if err != nil {
		return nil, &syncError{
			message:    "failed to create youtube service: " + err.Error(),
			statusCode: http.StatusInternalServerError,
		}
	}

	uploadsPlaylistID, syncErr := h.getUploadsPlaylistID(service, opztvChannelID)
	if syncErr != nil {
		return nil, syncErr
	}

	stats := &videoSyncStats{}
	if syncErr := h.fetchAndStoreVideos(ctx, service, uploadsPlaylistID, maxResults, stats); syncErr != nil {
		return nil, syncErr
	}

	return stats, nil
}

func (h *Handler) createYouTubeService() (*youtube.Service, error) {
	return youtube.NewService(context.Background(), option.WithAPIKey(config.YOUTUBE_API_KEY))
}

func (h *Handler) getUploadsPlaylistID(service *youtube.Service, channelID string) (string, *syncError) {
	channelCall := service.Channels.List([]string{"contentDetails"}).Id(channelID)
	channelResponse, err := channelCall.Do()
	if err != nil {
		return "", &syncError{
			message:    "failed to fetch channel: " + err.Error(),
			statusCode: http.StatusInternalServerError,
		}
	}

	if len(channelResponse.Items) == 0 {
		return "", &syncError{
			message:    "channel not found",
			statusCode: http.StatusNotFound,
		}
	}

	return channelResponse.Items[0].ContentDetails.RelatedPlaylists.Uploads, nil
}

func (h *Handler) fetchAndStoreVideos(ctx context.Context, service *youtube.Service, playlistID string, maxResults int, stats *videoSyncStats) *syncError {
	pageToken := ""

	for stats.totalFetched < maxResults {
		items, nextToken, syncErr := h.fetchPlaylistPage(service, playlistID, pageToken)
		if syncErr != nil {
			return syncErr
		}

		if len(items) == 0 {
			break
		}

		if syncErr := h.processVideoItems(ctx, items, maxResults, stats); syncErr != nil {
			return syncErr
		}

		if nextToken == "" {
			break
		}
		pageToken = nextToken
	}

	return nil
}

func (h *Handler) fetchPlaylistPage(service *youtube.Service, playlistID, pageToken string) ([]*youtube.PlaylistItem, string, *syncError) {
	call := service.PlaylistItems.List([]string{"snippet"}).
		PlaylistId(playlistID).
		MaxResults(youtubeMaxPageSize)

	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	response, err := call.Do()
	if err != nil {
		return nil, "", &syncError{
			message:    "failed to fetch playlist items: " + err.Error(),
			statusCode: http.StatusInternalServerError,
		}
	}

	return response.Items, response.NextPageToken, nil
}

func (h *Handler) processVideoItems(ctx context.Context, items []*youtube.PlaylistItem, maxResults int, stats *videoSyncStats) *syncError {
	for _, item := range items {
		if stats.totalFetched >= maxResults {
			break
		}
		stats.totalFetched++

		h.processVideoItem(ctx, item, stats)
	}
	return nil
}

func (h *Handler) processVideoItem(ctx context.Context, item *youtube.PlaylistItem, stats *videoSyncStats) {
	videoID := item.Snippet.ResourceId.VideoId

	if h.videoExists(ctx, videoID) {
		stats.skipped++
		return
	}

	video := h.buildVideoModel(item, videoID)
	if err := h.repo.CreateVideo(ctx, video); err != nil {
		fmt.Printf("failed to insert video %s: %v\n", videoID, err)
		stats.errors++
		return
	}

	stats.added++
}

func (h *Handler) videoExists(ctx context.Context, videoID string) bool {
	existingVideo, err := h.repo.GetVideoByYouTubeID(ctx, videoID)
	return err == nil && existingVideo != nil
}

func (h *Handler) buildVideoModel(item *youtube.PlaylistItem, videoID string) repoModel.Videos {
	publishedAt := h.parsePublishedDate(item.Snippet.PublishedAt)
	thumbnail := h.extractThumbnailURL(item.Snippet.Thumbnails)
	now := time.Now()

	return repoModel.Videos{
		ID:          toNullableString(videoID),
		Title:       item.Snippet.Title,
		Thumbnail:   toNullableString(thumbnail),
		PublishedAt: publishedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (h *Handler) parsePublishedDate(dateStr string) time.Time {
	publishedAt, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Now()
	}
	return publishedAt
}

func (h *Handler) extractThumbnailURL(thumbnails *youtube.ThumbnailDetails) string {
	if thumbnails != nil && thumbnails.High != nil {
		return thumbnails.High.Url
	}
	return ""
}

func toNullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
