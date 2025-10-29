package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jchiu21/streamsense/server/database"
	"github.com/jchiu21/streamsense/server/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieCollection, err := database.OpenCollection("movies")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error."})
			return
		}

		// set up context timeout (10 seconds max for query)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// query: find all movies with no filter
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
			return
		}
		defer cursor.Close(ctx)

		var movies []models.Movie

		// read database results from cursor into movie slice
		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieCollection, err := database.OpenCollection("movies")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error."})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required."})
			return
		}

		var movie models.Movie

		// find movie in database based on imdb_id param and store in movie var
		if err = movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found."})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieCollection, err := database.OpenCollection("movies")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error."})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// bind incoming http request JSON body to Movie struct
		var movie models.Movie
		if err = c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while binding incoming HTTP request body."})
			return
		}

		// validate movie from struct tags
		if err = validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed.", "details": err.Error()})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie."})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}
