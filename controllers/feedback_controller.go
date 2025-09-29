package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"upload-service/configs"
	"upload-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getFeedbackCollection() *mongo.Collection {
	return configs.GetCollection(configs.DB, "feedback")
}

// SubmitFeedback handles POST requests to submit feedback
func SubmitFeedback() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		logger := configs.LogWithContext("feedback", "submit")
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var feedbackRequest models.FeedbackRequest
		if err := json.NewDecoder(r.Body).Decode(&feedbackRequest); err != nil {
			logger.Error("Failed to decode feedback request", "error", err)
			errorResponse(rw, err, http.StatusBadRequest)
			return
		}

		logger.Debug("Feedback request received", "has_email", feedbackRequest.Email != "", "has_phone", feedbackRequest.Phone != "")

		// Validate that at least email or phone is provided
		if feedbackRequest.Email == "" && feedbackRequest.Phone == "" {
			logger.Warn("Feedback validation failed: no contact information provided")
			errorResponse(rw, &json.SyntaxError{}, http.StatusBadRequest)
			return
		}

		// Validate that comment is not empty
		if feedbackRequest.Comment == "" {
			logger.Warn("Feedback validation failed: comment is empty")
			errorResponse(rw, &json.SyntaxError{}, http.StatusBadRequest)
			return
		}

		feedback := models.Feedback{
			Email:       feedbackRequest.Email,
			Phone:       feedbackRequest.Phone,
			Comment:     feedbackRequest.Comment,
			DateCreated: time.Now(),
		}

		result, err := getFeedbackCollection().InsertOne(ctx, feedback)
		if err != nil {
			logger.Error("Failed to insert feedback", "error", err)
			errorResponse(rw, err, http.StatusInternalServerError)
			return
		}

		logger.Info("Feedback submitted successfully", "id", result.InsertedID)
		configs.LogDatabaseOperation("insert", "feedback", start, nil)
		successResponse(rw, result.InsertedID)
	}
}

// GetFeedback handles GET requests to retrieve all feedback
func GetFeedback() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		logger := configs.LogWithContext("feedback", "get-all")
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Set up options for sorting by date created (newest first)
		findOptions := options.Find()
		findOptions.SetSort(bson.M{"date_created": -1})

		cursor, err := getFeedbackCollection().Find(ctx, bson.M{}, findOptions)
		if err != nil {
			logger.Error("Failed to query feedback collection", "error", err)
			errorResponse(rw, err, http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var feedback []models.Feedback
		if err = cursor.All(ctx, &feedback); err != nil {
			logger.Error("Failed to decode feedback documents", "error", err)
			errorResponse(rw, err, http.StatusInternalServerError)
			return
		}

		// Ensure we always return an empty array instead of null
		if feedback == nil {
			feedback = []models.Feedback{}
		}

		logger.Info("Feedback retrieved successfully", "count", len(feedback))
		configs.LogDatabaseOperation("find", "feedback", start, nil)

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(feedback)
	}
}
