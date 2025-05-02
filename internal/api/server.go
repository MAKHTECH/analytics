package api

import (
	api "analytics/internal/api/gen"
	"analytics/internal/api/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

func InitServer(address string, port int, service handlers.AnalyticService) (server *http.Server) {
	// Инициализация HTTP сервера для API и метрик
	serverApi := handlers.NewServerAPI(service)
	httpServer, err := api.NewServer(serverApi)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", httpServer)
	mux.Handle("/metrics", promhttp.Handler())

	server = &http.Server{
		Addr:         address + ":" + strconv.Itoa(port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Запуск HTTP сервера в отдельной горутине
	go func() {
		log.Printf("Starting HTTP server on port %d", port)
		if err = server.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	return
}
