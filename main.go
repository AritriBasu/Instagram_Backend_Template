package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	UserID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Post struct {
	PostID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption    string             `json:"caption,omitempty" bson:"_id,omitempty"`
	ImageURL   url.URL            `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	PostedTime time.Time          `json:"postedtime,omitempty" bson:"postedtime,omitempty"`
	PostUser   []User             `json:"postuser" bson:"posteduser"`
}

func createUser(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-type","application/json")
	var newuser User
	json.NewDecoder(request,Body).Decode(&newuser)
	collection := client.Database("myFirstDatabase").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx,newuser)
	json.NewEncoder(response).Encode(result)
}


func createPost(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var post Post
	json.NewDecoder(request,Body).Decode(&user)
	collection := client.Database("myFirstDatabase").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx,post)
	json.NewEncoder(response).Encode(result)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	uid := param["uid"]
	//fmt.Println(uid)
	var user User
	collection := client.Database("myFirstDatabase").Collection("users")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{UserID: uid}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "User not found!")
		return
	}
	json.NewEncoder(w).Encode(newuser)
}

func getPost(w http.ResponseWriter, r *http.Request){
	response.Header().Add("content-type","application/json")
	params := mux.Vars(request)
	param := mux.Vars(r)
	uid := param["pid"]
	var post Post
    collection := client.Database("myFirstDatabase").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err := collection.FindOne(ctx,Post{ID: pid}).Decode(&post)
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Post not found!")
		return
	}
	json.NewEncoder(response).Encode(post)
}

func getUserPost(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	pid := param["pid"]
	var userpost [] allpost
	collection := client.Database("myFirstDatabase").Collection("userposts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx,Post{ID: pid}).Decode(&post)
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "User not found!")
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
        var post Post 
		cursor.Decode(&post)
		userpost = append(userpost,Post)
	}
	if err !=nil{
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(([]byte(`{"message":"`+err.Error()+`"}`))
		return
	}
	json.NewEncoder(response).Encode(userpost)
}

func main() {
	//connecting to mongo db
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://Aritri:mongodb33@cluster0.pfccy.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	client, err= mongo.Connect(ctx, clientOptions)
	//to disconnect once control is returned to main
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/user/{uid}", getUser).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts/{pid}", getPost).Methods("GET")
	router.HandleFunc("/posts/users/{uid}", getUserPost).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}
