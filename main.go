package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Movie struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Genre    string `json:"genre"`
	Director string `json:"director"`
}

var users = []User{
	{Username: "admin", Password: "admin123", Role: "admin"},
	{Username: "user", Password: "user123", Role: "user"},
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
	router.POST("/register", register)
	router.POST("/login", login)
	userRoutes := router.Group("/user")
	userRoutes.Use(authMiddleware("user"))
	{
		userRoutes.GET("/movies", getMovies)
		userRoutes.GET("/movies/:id", getMovieByID)
	}
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(authMiddleware("admin"))
	{
		adminRoutes.POST("/movies", addMovie)
		adminRoutes.PUT("/movies/:id", updateMovie)
		adminRoutes.DELETE("/movies/:id", deleteMovie)
	}
	router.Run(":8080")
}

func generateToken(username string, role string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func register(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	for _, user := range users {
		if user.Username == newUser.Username {
			c.JSON(http.StatusConflict, gin.H{"message": "username already exists"})
			return
		}
	}
	users = append(users, newUser)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	var user User
	for _, u := range users {
		if u.Username == creds.Username && u.Password == creds.Password {
			user = u
			break
		}
	}
	if user.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	token, err := generateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func authMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}
		if claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"message": "access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func getMovies(c *gin.Context) {
	c.JSON(http.StatusOK, movies)
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
