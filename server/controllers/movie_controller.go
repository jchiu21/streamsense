package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jchiu21/streamsense/server/database"
	"github.com/jchiu21/streamsense/server/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieCollection, err := database.OpenCollection("movies")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error."})
			return
		}

		// set up timeout (10 seconds max for query)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// query: find all movies with no filter
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		defer cursor.Close(ctx)

		var movies []models.Movie

		// read database results from cursor into movie slice
		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
		}

		c.JSON(http.StatusOK, movies)
	}
}
