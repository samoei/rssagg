package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/samoei/rssagg/internal/database"
)

func (apiConfig *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	feed_follow, err := apiConfig.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed follow: %v", err))
		return
	}

	repondWithJSON(w, 201, feed_follow)
}

func (apiConfig *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feed_follows, err := apiConfig.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error in fetching feed follows: %v", err))
		return
	}

	repondWithJSON(w, 200, feed_follows)
}

func (apiConfig *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feed-follow-id")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error in parsing feed follow id: %v", err))
		return
	}

	err = apiConfig.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     feedFollowId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error in deleting feed follow: %v", err))
		return
	}

	repondWithJSON(w, 200, struct{}{})
}
