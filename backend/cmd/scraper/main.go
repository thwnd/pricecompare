package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// Item 구조체는 수집된 데이터를 저장합니다
type Item struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Image       string    `json:"image"`
	Price       string    `json:"price"`
	Source      string    `json:"source"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	// 데이터를 저장할 슬라이스
	var items []Item

	// 컬렉터 인스턴스 생성
	c := colly.NewCollector(
		colly.AllowedDomains("quotes.toscrape.com"), // 테스트용 도메인으로 변경
		colly.MaxDepth(2), // 크롤링 깊이 제한
	)

	// 타임아웃 설정
	c.SetRequestTimeout(120 * time.Second)

	// 병렬 처리 설정
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	// User-Agent 설정
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// 방문한 URL 추적
	visitedURLs := make(map[string]bool)

	// 테스트용 quotes.toscrape.com에서 인용구 수집
	c.OnHTML(".quote", func(e *colly.HTMLElement) {
		item := Item{
			Title:       e.ChildText(".text"),
			Description: e.ChildText(".author"),
			URL:         e.Request.URL.String(),
			Image:       "", // 테스트 사이트에는 이미지가 없음
			Price:       "",
			Source:      "Quotes To Scrape",
			Category:    "quotes",
			CreatedAt:   time.Now(),
		}

		// 태그를 Category에 저장
		tags := e.ChildTexts(".tag")
		if len(tags) > 0 {
			item.Category = tags[0]
		}

		// 데이터가 유효한지 확인
		if item.Title != "" {
			items = append(items, item)
			fmt.Printf("인용구 수집: %s\n", item.Title)
		}
	})

	// 다음 페이지 링크 찾기
	c.OnHTML(".next a", func(e *colly.HTMLElement) {
		nextPageURL := e.Request.AbsoluteURL(e.Attr("href"))
		if !visitedURLs[nextPageURL] {
			visitedURLs[nextPageURL] = true
			c.Visit(nextPageURL)
		}
	})

	// 에러 처리
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("요청 에러: %s\n", err)
	})

	// 페이지 방문 완료 후 처리
	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("페이지 수집 완료: %s\n", r.Request.URL)
	})

	// 시작 URL 방문
	fmt.Println("스크래핑 시작...")
	c.Visit("https://quotes.toscrape.com/")

	// 결과를 JSON 파일로 저장
	writeItemsToJSON(items, "quotes.json")
	fmt.Printf("스크래핑 완료! 총 %d개 인용구 수집\n", len(items))
}

// 수집된 데이터를 JSON 파일로 저장
func writeItemsToJSON(items []Item, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("파일 생성 에러: %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(items)
}
