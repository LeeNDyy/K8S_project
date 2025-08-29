package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var (
	statusCache = make(map[string]ServiceStatus)
	lastChecked = time.Now().Format("02.01.2006 15:04:05")
)

func main() {
	engine := html.New("./web/page", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./web/static")

	go startBackgroundMonitoring(30 * time.Second)

	app.Get("/", func(c *fiber.Ctx) error {
		stats := GetStats(statusCache)
		return c.Render("index", fiber.Map{
			"StatusMap":   statusCache,
			"Stats":       stats,
			"LastChecked": lastChecked,
		})
	})

	app.Get("/status", func(c *fiber.Ctx) error {
		stats := GetStats(statusCache)
		return c.Render("status_partial", fiber.Map{
			"StatusMap":   statusCache,
			"Stats":       stats,
			"LastChecked": lastChecked,
		})
	})

	app.Get("/api/stats", func(c *fiber.Ctx) error {
		stats := GetStats(statusCache)
		return c.JSON(stats)
	})

	app.Get("/api/services", func(c *fiber.Ctx) error {
		return c.JSON(statusCache)
	})

	log.Println("🚀 Сервер запущен на http://localhost:8000")
	log.Fatal(app.Listen(":8000"))
}

func startBackgroundMonitoring(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("🔍 Проверяем сервисы...")
		statusCache = CheckAllServices()
		lastChecked = time.Now().Format("02.01.2006 15:04:05")

		stats := GetStats(statusCache)
		if stats.DownServices > 0 {
			log.Printf("🚨 Обнаружено проблем: %d сервисов не работают", stats.DownServices)
		}
	}
}
