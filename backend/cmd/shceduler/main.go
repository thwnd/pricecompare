package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	// 로그 설정
	logFile, err := os.OpenFile("scheduler.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("로그 파일 열기 오류: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// 크론 스케줄러 생성
	c := cron.New(cron.WithSeconds())

	// 매일 자정에 스크래퍼 실행 (0 0 0 * * *)
	_, err = c.AddFunc("0 0 0 * * *", func() {
		runScraper()
	})
	if err != nil {
		log.Fatalf("크론 작업 추가 중 오류: %v", err)
	}

	// 애플리케이션 시작 시 바로 한 번 실행
	log.Println("시작 시 스크래퍼 실행...")
	runScraper()

	// 스케줄러 시작
	c.Start()
	log.Println("스케줄러가 시작되었습니다. 매일 자정에 스크래퍼가 실행됩니다.")

	// 프로그램이 종료되지 않도록 대기
	select {}
}

// 스크래퍼 실행 함수
func runScraper() {
	startTime := time.Now()
	log.Println("스크래퍼 실행 시작...")

	// 스크래퍼 프로그램 경로 설정
	scraperCmd := exec.Command("go", "run", "../scraper/main.go")

	// 스크래퍼의 출력을 로그에 연결
	scraperCmd.Stdout = os.Stdout
	scraperCmd.Stderr = os.Stderr

	// 스크래퍼 실행
	err := scraperCmd.Run()
	if err != nil {
		log.Printf("스크래퍼 실행 중 오류: %v", err)
		return
	}

	elapsed := time.Since(startTime)
	log.Printf("스크래퍼 실행 완료 (소요 시간: %v)\n", elapsed)

	// 데이터베이스 연결 상태 확인 또는 다른 필요한 작업 수행 가능
	fmt.Println("데이터 수집 작업이 완료되었습니다.")
}
