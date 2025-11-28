package controllers

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"upload-service/configs"
	"upload-service/models"
)

var (
	cleanupInProgress = make(map[string]bool)
	cleanupMutex      sync.Mutex
)

func MonitorLiveStreams() {
	ticker := time.NewTicker(12 * time.Second)
	defer ticker.Stop()

	fmt.Println("Live stream monitor started...")

	for range ticker.C {
		ensureAllStream()
	}
}

func ensureAllStream() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find streams marked as live
	filter := bson.M{
		"type":      TYPE_STREAM,
		"isdeleted": false,
		"is_live":   true,
	}

	cursor, err := getContentCollection().Find(ctx, filter)
	if err != nil {
		fmt.Println("Error finding streams:", err)
		return
	}
	defer cursor.Close(ctx)

	var streams []models.Content
	if err = cursor.All(ctx, &streams); err != nil {
		fmt.Println("Error decoding streams:", err)
		return
	}

	if len(streams) > 0 {
		fmt.Printf("Checking %d live streams...\n", len(streams))
	}

	for _, stream := range streams {
		checkSingleStream(stream)
	}
}

func checkSingleStream(stream models.Content) {
	if !stream.IsLive {
		return
	}

	streamingServerIP := configs.EnvStreamingServer()
	hlsURL := fmt.Sprintf("http://%s/hls/%s/index.m3u8", streamingServerIP, stream.StreamKey)

	lastModified, exists := checkHLSLastModified(hlsURL)

	if exists && time.Since(lastModified) > 45*time.Second {
		// Orphaned stream detected
		fmt.Printf("🔴 ORPHANED STREAM: %s (last modified: %v ago)\n", 
			stream.StreamKey, time.Since(lastModified))
		go triggerRemoteCleanup(stream.StreamKey)
	}
}

func checkHLSLastModified(hlsURL string) (time.Time, bool) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Head(hlsURL)
	if err != nil {
		return time.Time{}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return time.Time{}, false
	}

	lastModStr := resp.Header.Get("Last-Modified")
	if lastModStr == "" {
		return time.Now(), true
	}

	lastMod, err := http.ParseTime(lastModStr)
	if err != nil {
		return time.Now(), true
	}

	return lastMod, true
}

func triggerRemoteCleanup(streamKey string) {
	cleanupMutex.Lock()
	if cleanupInProgress[streamKey] {
		cleanupMutex.Unlock()
		return
	}
	cleanupInProgress[streamKey] = true
	cleanupMutex.Unlock()

	fmt.Printf("🧹 Triggering cleanup for: %s\n", streamKey)

	cleanupURL := fmt.Sprintf("http://%s/api/cleanup/%s", configs.EnvStreamingServer(), streamKey)
	resp, err := http.Post(cleanupURL, "application/json", nil)
	if err != nil {
		fmt.Printf("❌ Failed to trigger cleanup: %v\n", err)
	} else {
		resp.Body.Close()
		fmt.Printf("✅ Cleanup triggered for: %s\n", streamKey)
	}

	time.AfterFunc(1*time.Minute, func() {
		cleanupMutex.Lock()
		delete(cleanupInProgress, streamKey)
		cleanupMutex.Unlock()
	})
}