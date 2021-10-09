/*
An instagram-api using golang and mongodb

Author: Ambuj Gupta
Reg No: 19BCE0366
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Post struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption   string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageURL  string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	Timestamp time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	UserId    primitive.ObjectID `json:"_userid,omitempty" bson:"_userid,omitempty"`
}

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

var client *mongo.Client

func main() {
	// Configuration
	mongoDbURI := "mongodb+srv://dbUser:appointy@distronx.lz6vr.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(mongoDbURI)

	client, _ = mongo.NewClient(clientOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//setup connection to MongoDB Atlas Cluster
	err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Error: ", err)
	} else {
		log.Println("Connected.")
	}

	handleRequest()
}

//Handles all the https routes
func handleRequest() {
	//root
	http.HandleFunc("/", root)

	//Posts Requests
	http.HandleFunc("/posts", getPosts)
	http.HandleFunc("/posts/", getPostById)
	http.HandleFunc("/posts/users/", getPostsOfUser)

	//Users Requests
	http.HandleFunc("/users", getAllUsers)
	http.HandleFunc("/users/", getUserById)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

// Handles the / route
func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Appointy-Instagram-API\n\nAmbuj Gupta\n19BCE0366\nVellore Institute of Technology")
	fmt.Println("Displaying Root.")
}

// Get all posts if Method=="GET"
// Create new post of Method=="POST"
func getPosts(response http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		var posts []Post
		response.Header().Set("content-type", "application/json")

		collection := client.Database("test").Collection("Post")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var post Post
			cursor.Decode(&post)
			posts = append(posts, post)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		fmt.Println("Endpoint Hit: Get All Posts")
		json.NewEncoder(response).Encode(posts)
	} else {
		request.ParseForm()
		decoder := json.NewDecoder(request.Body)
		var newPost Post

		newPost.Timestamp = time.Now()

		err := decoder.Decode(&newPost)
		if err != nil {
			panic(err)
		}
		log.Println(newPost.Id)
		fmt.Println("Endpoint Hit: Create New Post")
		insertPost(newPost)
	}
}

// Helper function that inserts the new post into DB
func insertPost(post Post) {
	collection := client.Database("test").Collection("Post")
	insertResult, err := collection.InsertOne(context.TODO(), post)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}

// Get the Post with post id given in the route /posts/{id}
func getPostById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	id := strings.TrimPrefix(request.URL.Path, "/posts/")
	objID, _ := primitive.ObjectIDFromHex(id)

	var post Post
	collection := client.Database("test").Collection("Post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, Post{Id: objID}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	fmt.Println("Returned Post Id: ", post.Id)
	json.NewEncoder(response).Encode(post)
}

func getAllUsers(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("content-type", "application/json")
		var users []User
		collection := client.Database("test").Collection("User")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var user User
			cursor.Decode(&user)
			users = append(users, user)
		}
		if err = cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		fmt.Println("Endpoint Hit: Get All Users")
		json.NewEncoder(response).Encode(users)
	} else {
		request.ParseForm()
		decoder := json.NewDecoder(request.Body)
		var newUser User

		err := decoder.Decode(&newUser)
		if err != nil {
			panic(err)
		}
		newUser.Password = hashPass(newUser.Password) //Password Encryption
		log.Println(newUser.Id)
		fmt.Println("Endpoint Hit: User Created")
		insertUser(newUser)
	}
}

// Helper function that inserts the new user into DB
func insertUser(user User) {
	collection := client.Database("test").Collection("User")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created user with ID:", insertResult.InsertedID)
}

// Get the User with user id given in the route /users/{id}
func getUserById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	id := strings.TrimPrefix(request.URL.Path, "/users/")
	objID, _ := primitive.ObjectIDFromHex(id)

	var user User
	collection := client.Database("test").Collection("User")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, User{Id: objID}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	fmt.Println("Returned User Id: ", user.Id)
	json.NewEncoder(response).Encode(user)
}

// Get all the posts of user with user id mentioned in route /posts/users/{id}
func getPostsOfUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	id := strings.TrimPrefix(request.URL.Path, "/posts/users/")
	objID, _ := primitive.ObjectIDFromHex(id)

	var posts []Post
	collection := client.Database("test").Collection("Post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, Post{UserId: objID})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var post Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}
	if err = cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	fmt.Printf("Endpoint Hit: Get Posts of UserId %s\n", id)
	json.NewEncoder(response).Encode(posts)
}

// Password hashing function
func hashPass(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}
