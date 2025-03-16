package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// Item 구조체
type Item struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	URL         string             `json:"url" bson:"url"`
	Image       string             `json:"image" bson:"image"`
	Price       string             `json:"price" bson:"price"`
	Source      string             `json:"source" bson:"source"`
	Category    string             `json:"category" bson:"category"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

var collection *mongo.Collection

func main() {
	// MongoDB 연결
	connectDB()

	// Gin 라우터 설정
	router := gin.Default()

	// CORS 미들웨어 설정
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API 라우트 정의
	api := router.Group("/api")
	{
		api.GET("/items", getItems)
		api.GET("/items/:category", getItemsByCategory)
		api.GET("/item/:id", getItemByID)
		api.GET("/search", searchItems)
	}

	// 정적 파일 제공
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// 서버 시작
	log.Println("서버 시작: http://localhost:8080")
	router.Run(":8080")
}

// MongoDB 연결 함수
func connectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// 연결 확인
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB에 연결되었습니다")
	collection = client.Database("aggregator").Collection("items")
}

// 전체 아이템 조회 API
func getItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 페이지네이션 처리
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var items []Item
	if err = cursor.All(ctx, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 전체 아이템 수 조회
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"limit":     limit,
		"totalPage": (total + int64(limit) - 1) / int64(limit),
	})
}

// 카테고리별 아이템 조회 API
func getItemsByCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	category := c.Param("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	skip := (page - 1) * limit

	filter := bson.M{"category": category}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var items []Item
	if err = cursor.All(ctx, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"limit":     limit,
		"totalPage": (total + int64(limit) - 1) / int64(limit),
	})
}

// 개별 아이템 조회 API
func getItemByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := c.Param("id")

	// 유효한 ObjectID인지 확인
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "유효하지 않은 ID 형식입니다"})
		return
	}

	var item Item
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "아이템을 찾을 수 없습니다"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// 아이템 검색 API
func searchItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "검색어가 필요합니다"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	skip := (page - 1) * limit

	// 텍스트 검색 필터
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var items []Item
	if err = cursor.All(ctx, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"limit":     limit,
		"totalPage": (total + int64(limit) - 1) / int64(limit),
	})
}
