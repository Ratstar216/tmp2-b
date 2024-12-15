package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	"tmp2-backend/gemini"
	"tmp2-backend/models" // Assuming your models are in the 'models' package

	"time"

	"github.com/gorilla/mux"
	"github.com/oklog/ulid/v2"
	"gorm.io/driver/mysql" // Or your preferred database driver
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Database setup (using SQLite for this example)

	// err := godotenv.Load()
	// err := ""
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	DB_USER := os.Getenv("MYSQL_USER")
	DB_PASS := os.Getenv("MYSQL_PASSWORD")
	DB_NAME := os.Getenv("MYSQL_DATABASE")
	DB_HOST := os.Getenv("MYSQL_HOST")
	dbPort := "3306"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASS, DB_HOST, dbPort, DB_NAME)
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db = _db

	// Migrate the schema
	// db.AutoMigrate(&models.Tweet{}, &models.User{})

	// Sample data (remove in production)
	// createSampleData()

	// Router setup
	r := mux.NewRouter()

	// API endpoints
	r.HandleFunc("/tweets", getTweets).Methods("GET")
	r.HandleFunc("/tweets/{id}", getTweet).Methods("GET")
	r.HandleFunc("/tweets", createTweet).Methods("POST")
	r.HandleFunc("/tweets/{id}/reply", createReply).Methods("POST")
	r.HandleFunc("/tweets/{id}/like", likeTweet).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow requests from your frontend origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Wrap your router with the CORS handler
	handler := c.Handler(r)

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// getTweets handles GET /tweets - retrieves a list of tweets (without replies)
func getTweets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tweets []models.Tweet
	// db.Where("parent_id IS NULL").Order("created_at desc").Find(&tweets) // Fetch only top-level tweets TODO ちゃんと null にする
	// db.Where("parent_id = \"\"").Order("created_at desc").Find(&tweets)
	db.Order("created_at desc").Find(&tweets)
	json.NewEncoder(w).Encode(tweets)
}

// getTweet handles GET /tweets/{id} - retrieves a single tweet with its replies
func getTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	// id, err := strconv.Atoi(vars["id"])
	id := vars["id"]
	fmt.Println("id", id)
	fmt.Println("vars", vars)
	// if err != nil {
	// 	http.Error(w, "Invalid tweet ID", http.StatusBadRequest)
	// 	return
	// }

	// var tweet models.Tweet
	// // Preload Replies to eagerly load associated replies
	// if err := db.Preload("Replies").First(&tweet, id).Error; err != nil {
	// 	http.Error(w, "Tweet not found", http.StatusNotFound)
	// 	fmt.Println(err)
	// 	return
	// }
	// json.NewEncoder(w).Encode(tweet)

	// Preload Replies and use Where clause to filter by id
	// if err := db.Preload("Replies").Where("id = ?", id).First(&tweet).Error; err != nil {
	// 	http.Error(w, "Tweet not found", http.StatusNotFound)
	// 	fmt.Println(err)
	// 	return
	// }
	// json.NewEncoder(w).Encode(tweet)

	var tweets []models.Tweet
	// db.Where("parent_id IS NULL").Order("created_at desc").Find(&tweets) // Fetch only top-level tweets TODO ちゃんと null にする
	db.Where("id = ?", id).Find(&tweets)
	json.NewEncoder(w).Encode(tweets)

}

// createTweet handles POST /tweets - creates a new tweet
func createTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tweet models.Tweet

	json.NewDecoder(r.Body).Decode(&tweet)
	// In a real application, you would get the user ID from the authentication context.
	// For this example, we are hardcoding a user ID.
	// tweet.Username = "testuser"     // Replace with actual username
	// tweet.UserID = "abc"            // Replace with actual user ID
	tweet.ID = ulid.Make().String() // Replace with actual tweet ID
	tweet.Likes = 0
	now := time.Now()
	fmt.Println("tweet.Content", tweet.Content)
	tweet.Content = tweet.Content + gemini.Translate(tweet.Content)
	fmt.Println("tweet.Content", tweet.Content)

	// 時刻を文字列に変換（フォーマット指定）
	tweet.Created_at = now.Format("2006-01-02 15:04:05")

	db.Create(&tweet)
	json.NewEncoder(w).Encode(tweet)
}

// createReply handles POST /tweets/{id}/reply - creates a reply to a tweet
func createReply(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	// parentID, err := strconv.Atoi(vars["id"])
	// if err != nil {
	// 	http.Error(w, "Invalid parent tweet ID", http.StatusBadRequest)
	// 	return
	// }
	parentID := vars["id"]

	var reply models.Tweet
	json.NewDecoder(r.Body).Decode(&reply)

	// In a real application, you would get the user ID from the authentication context.
	// For this example, we are hardcoding a user ID.
	reply.UserID = "1" // Replace with actual user ID
	reply.ID = ulid.Make().String()
	reply.Likes = 0
	now := time.Now()
	reply.Created_at = now.Format("2006-01-02 15:04:05")

	// reply.ParentID = &[]uint{uint(parentID)}[0] // Set the ParentID
	reply.ParentID = parentID

	db.Create(&reply)
	json.NewEncoder(w).Encode(reply)
}

func likeTweet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	tweetID := vars["id"]
	// if err != nil {
	// 	http.Error(w, "Invalid tweet ID", http.StatusBadRequest)
	// 	return
	// }

	// 本来は認証されたユーザー情報を取得する
	// userID := uint(1) // 仮のユーザーID

	// var tweet models.Tweet
	// if err := db.First(&tweet, tweetID).Error; err != nil {
	// 	http.Error(w, "Tweet not found", http.StatusNotFound)
	// 	return
	// }

	// // 既にいいねしているか確認 (効率化の余地あり)
	// var existingLike models.UserLikes
	// result := db.Where("user_id = ? AND tweet_id = ?", userID, tweetID).First(&existingLike)
	// if result.Error == nil {
	// 	http.Error(w, "Already liked", http.StatusConflict)
	// 	return
	// }

	// // いいねを関連付ける
	// like := models.UserLikes{UserID: userID, TweetID: uint(tweetID)}
	// if err := db.Create(&like).Error; err != nil {
	// 	http.Error(w, "Failed to like tweet", http.StatusInternalServerError)
	// 	return
	// }

	// Tweet の Likes フィールドをインクリメント (データベースには保存されない)
	// db.Model(&tweet).UpdateColumn("likes", gorm.Expr("likes + 1"))
	// db.Set("likes", "likes + 1").Where("id = ?", tweetID).Update("tweets")

	result := db.Model(&models.Tweet{}).Where("id = ?", tweetID).Update("likes", gorm.Expr("likes + ?", 1))
	if result.Error != nil {
		fmt.Println("Error updating likes:", result.Error)
		return
	}

	// 更新されたツイートを返す
	// db.Preload("Replies").First(&tweet, tweetID)
	// json.NewEncoder(w).Encode(tweet)
}

// createSampleData inserts some sample tweets and replies (for testing purposes)
func createSampleData() {
	user := models.User{Username: "testuser"}
	db.FirstOrCreate(&user, user)

	tweet1 := models.Tweet{ID: "121232", UserID: user.ID, Content: "First tweet!"}
	db.FirstOrCreate(&tweet1, tweet1)

	reply1 := models.Tweet{ID: "142432", UserID: user.ID, Content: "Reply to first tweet", ParentID: tweet1.ID}
	db.FirstOrCreate(&reply1, reply1)

	reply2 := models.Tweet{ID: "34342", UserID: user.ID, Content: "Another reply", ParentID: tweet1.ID}
	db.FirstOrCreate(&reply2, reply2)

	tweet2 := models.Tweet{ID: "123535532", UserID: user.ID, Content: "Second tweet!"}
	db.FirstOrCreate(&tweet2, tweet2)

	fmt.Println("Sample data created.")
}
