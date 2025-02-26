package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Genre    string `json:"genre"`
	Director string `json:"director"`
}

var movies = []Movie{
	{ID: "1", Title: "Inception", Year: 2010, Genre: "Sci-Fi", Director: "Christopher Nolan"},
	{ID: "2", Title: "The Dark Knight", Year: 2008, Genre: "Action", Director: "Christopher Nolan"},
	{ID: "3", Title: "Interstellar", Year: 2014, Genre: "Sci-Fi", Director: "Christopher Nolan"},
	{ID: "4", Title: "The Matrix", Year: 1999, Genre: "Sci-Fi", Director: "Lana Wachowski, Lilly Wachowski"},
	{ID: "5", Title: "Fight Club", Year: 1999, Genre: "Drama", Director: "David Fincher"},
	{ID: "6", Title: "Forrest Gump", Year: 1994, Genre: "Drama", Director: "Robert Zemeckis"},
	{ID: "7", Title: "The Godfather", Year: 1972, Genre: "Crime", Director: "Francis Ford Coppola"},
	{ID: "8", Title: "The Godfather: Part II", Year: 1974, Genre: "Crime", Director: "Francis Ford Coppola"},
}

func main() {
	router := gin.Default()
	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovieByID)
	router.POST("/movies", addMovie)
	router.PUT("/movies/:id", updateMovie)
	router.DELETE("/movies/:id", deleteMovie)
	router.Run(":8080")
}

func getMovies(c *gin.Context) {
	c.JSON(http.StatusOK, movies)
}

func getMovieByID(c *gin.Context) {
	id := c.Param("id")
	for _, movie := range movies {
		if movie.ID == id {
			c.JSON(http.StatusOK, movie)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
}

func addMovie(c *gin.Context) {
	var newMovie Movie
	if err := c.ShouldBindJSON(&newMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movies = append(movies, newMovie)
	c.JSON(http.StatusCreated, newMovie)
}

func updateMovie(c *gin.Context) {
	id := c.Param("id")
	var updatedMovie Movie
	if err := c.ShouldBindJSON(&updatedMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, movie := range movies {
		if movie.ID == id {
			movies[i] = updatedMovie
			c.JSON(http.StatusOK, updatedMovie)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
}

func deleteMovie(c *gin.Context) {
	id := c.Param("id")
	for i, movie := range movies {
		if movie.ID == id {
			movies = append(movies[:i], movies[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Movie deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
}
