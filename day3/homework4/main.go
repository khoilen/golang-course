package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Profile  string `json:"profile"`
}

const usersFilePath = "./data/users.json"

func readUsers() []User {
	var users []User
	file, err := ioutil.ReadFile(usersFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			return users
		}
		panic(err)
	}

	err = json.Unmarshal(file, &users)
	if err != nil {
		panic(err)
	}

	return users
}

func saveUsers(users []User) {
	file, err := json.MarshalIndent(users, "", " ")

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(usersFilePath, file, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", os.ModePerm)
	}

	router.POST("/register", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		users := readUsers()
		for _, user := range users {
			if user.Username == newUser.Username {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
				return
			}
		}

		users = append(users, newUser)
		saveUsers(users)
		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	})

	router.POST("/login", func(c *gin.Context) {
		var loginUser User

		if err := c.ShouldBindJSON(&loginUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		users := readUsers()
		for _, user := range users {
			if user.Username == loginUser.Username && user.Password == loginUser.Password {
				c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
				return
			}
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
	})

	router.GET("/user/:username", func(c *gin.Context) {
		username := c.Param("username")
		users := readUsers()

		for _, user := range users {
			if user.Username == username {
				c.JSON(http.StatusOK, user)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	router.POST("/user/:username", func(c *gin.Context) {
		var updateUser User
		if err := c.ShouldBindJSON(&updateUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		users := readUsers()
		for i, user := range users {
			if user.Username == updateUser.Username {
				users[i].Username = updateUser.Username
				users[i].Password = updateUser.Password
				saveUsers((users))
				c.JSON(http.StatusOK, gin.H{"message": "User information updated"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	router.POST("/user/:username/profile", func(c *gin.Context) {
		username := c.Param("username")

		file, err := c.FormFile("profile")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		users := readUsers()

		for i, user := range users {
			if user.Username == username {
				filename := filepath.Join("uploads", filepath.Base(file.Filename))
				if err := c.SaveUploadedFile(file, filename); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				users[i].Profile = filename
				saveUsers(users)
				c.JSON(http.StatusOK, gin.H{"message": "Profile picture uploaded"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	})

	router.Run(":8080")
}
