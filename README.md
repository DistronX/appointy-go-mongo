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

## Testing
### Execution
![](/images/1_vscode.png)

### Create new user

