# appointy-instagram-api

Name: Ambuj Gupta  
Reg No: 19BCE0366  
College: Vellore Institute Of Technology  

## Description

Round 1 Task  
Database: MongoDB (Atlas Cloud Platform).  
Backend: Go  

## How-To

Start local instance by using `go run .` command  
Local instance starts at `localhost:8080/`  
I have used Postman to test the API.

## Documentation

| Route             | Method | Description                                 |
|-------------------|--------|---------------------------------------------|
| /users            | GET    | Get all the users from DB                   |
|                   | POST   | Create a new user                           |
| /users/{id}       | GET    | Find the user with the given ID             |
| /posts            | GET    | Get all the posts from DB                   |
|                   | POST   | Create a new post                           |
| /posts/{id}       | GET    | Find the post with the given ID             |
| /posts/users/{id} | GET    | Get all the posts of user with given userID |

## Database MongoDB Atlas
![](/images/9_DatabaseCollection.png)

## Testing
### Execution
![](/images/1_vscode.png)

### Create new user
![](/images/2_new_user.png)
![](/images/3_new_user.png)

### Get all users
![](/images/4_get_all_users.png)

### Create new post
![](/images/5_new_post.png)

### Get all posts
![](/images/6_get_all_posts.png)

### Create new post
![](/images/7_new_post.png)

### Getting posts form a particular user
![](/images/8_posts_from_user.png)
