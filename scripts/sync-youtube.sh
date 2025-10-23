#!/bin/bash

# YouTube Video Sync Script
# This script syncs videos from OPZTV YouTube channel to the database

# Configuration
API_BASE="${API_BASE:-http://localhost:8088}"
MAX_RESULTS="${MAX_RESULTS:-50}"

# Check if API key is provided
if [ -z "$YOUTUBE_API_KEY" ]; then
    echo "Error: YOUTUBE_API_KEY environment variable is required"
    echo "Usage: YOUTUBE_API_KEY=your_key ./sync-youtube.sh"
    exit 1
fi

echo "Syncing videos from OPZTV YouTube channel..."
echo "API Base: $API_BASE"
echo "Max Results: $MAX_RESULTS"
echo ""

# Make the API call
response=$(curl -s -X POST "$API_BASE/api/videos/sync?api_key=$YOUTUBE_API_KEY&max_results=$MAX_RESULTS")

# Check if curl was successful
if [ $? -ne 0 ]; then
    echo "Error: Failed to connect to API"
    exit 1
fi

# Parse and display results
echo "$response" | jq '.'

# Extract counts
added=$(echo "$response" | jq -r '.data.added // 0')
skipped=$(echo "$response" | jq -r '.data.skipped // 0')
total=$(echo "$response" | jq -r '.data.total // 0')

echo ""
echo "Sync completed!"
echo "  Added: $added videos"
echo "  Skipped: $skipped videos (already exist)"
echo "  Total processed: $total videos"
