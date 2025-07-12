package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeSubscriptionResponse struct {
	IsSubscribed bool   `json:"isSubscribed"`
	Message      string `json:"message,omitempty"`
}

func CheckYouTubeSubscription(c *gin.Context) {
	// Get channel ID from query parameter
	channelID := c.Query("channelId")
	if channelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel ID is required"})
		return
	}

	// Get Firebase ID from context (set by middleware)
	firebaseID, exists := c.Get("firebaseId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	_, ok := firebaseID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Firebase ID"})
		return
	}

	// Get Google access token from request header
	// The frontend should send this in the X-Google-Access-Token header
	googleAccessToken := c.GetHeader("X-Google-Access-Token")
	if googleAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "Google access token required. Please sign in with Google and pass the access token in X-Google-Access-Token header",
			"instructions": "To get the Google access token, use the GoogleAuthProvider.getCredential() method in your frontend after Google sign-in",
		})
		return
	}

	// Check YouTube subscription using the access token
	isSubscribed, err := checkYouTubeSubscriptionStatus(c.Request.Context(), googleAccessToken, channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check subscription: " + err.Error()})
		return
	}

	response := YouTubeSubscriptionResponse{
		IsSubscribed: isSubscribed,
		Message:      "Subscription status checked successfully",
	}

	c.JSON(http.StatusOK, response)
}

func checkYouTubeSubscriptionStatus(ctx context.Context, accessToken, channelID string) (bool, error) {
	// Create YouTube service with the access token
	youtubeService, err := youtube.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)))
	if err != nil {
		return false, fmt.Errorf("failed to create YouTube service: %v", err)
	}

	// Get the authenticated user's channel ID
	channelsCall := youtubeService.Channels.List([]string{"id"})
	channelsCall.Mine(true)
	channelsResponse, err := channelsCall.Do()
	if err != nil {
		return false, fmt.Errorf("failed to get user channels: %v", err)
	}

	if len(channelsResponse.Items) == 0 {
		return false, fmt.Errorf("no channels found for user")
	}

	userChannelID := channelsResponse.Items[0].Id

	// Check if the user is subscribed to the target channel
	subscriptionsCall := youtubeService.Subscriptions.List([]string{"snippet"})
	subscriptionsCall.ChannelId(userChannelID)
	subscriptionsCall.ForChannelId(channelID)
	subscriptionsResponse, err := subscriptionsCall.Do()
	if err != nil {
		if googleapi.IsNotModified(err) {
			// User is not subscribed
			return false, nil
		}
		return false, fmt.Errorf("failed to check subscriptions: %v", err)
	}

	// Check if there are any subscriptions (user is subscribed)
	return len(subscriptionsResponse.Items) > 0, nil
}
