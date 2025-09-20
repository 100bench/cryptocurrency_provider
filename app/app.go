package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/100bench/cryptocurrency_provider.git/deployment/config"
	"github.com/100bench/cryptocurrency_provider.git/internal/adapters/broker/kafka"
	"github.com/100bench/cryptocurrency_provider.git/internal/adapters/external_client/coindesk"
	"github.com/100bench/cryptocurrency_provider.git/internal/adapters/storage/postgres"
	"github.com/100bench/cryptocurrency_provider.git/internal/cases"
	"github.com/100bench/cryptocurrency_provider.git/internal/ports/http/public"
)

func RunApp() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	// инициализация хранилища
	storage, err := postgres.NewPgxClient(ctx, cfg.PostgresDSN())
	if err != nil {
		return errors.Wrap(err, "postgres.NewPgxClient")
	}
	defer storage.Close()

	// инициализация внешнего клиента
	client, err := coindesk.NewClientCoinDesk(cfg.CoinDeskAPIURL)
	if err != nil {
		return errors.Wrap(err, "coindesk.NewClientCoinDesk")
	}
	// инициализация брокера
	broker := kafka.NewBroker(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)
	defer func() {
		_ = broker.Close()
	}()

	// конструкторы, продумать последовательность
	prodService, err := cases.NewProducer(broker)
	if err != nil {
		return errors.Wrap(err, "cases.NewProducer")
	}
	storageService, err := cases.NewStorageService(storage)
	if err != nil {
		return errors.Wrap(err, "cases.NewStorageService")
	}
	consumerService, err := cases.NewConsumer(broker, storage)
	if err != nil {
		return errors.Wrap(err, "cases.NewConsumer")
	}
	apiService, err := cases.NewServiceAPI(client, storage)
	if err != nil {
		return errors.Wrap(err, "cases.NewServiceAPI")
	}

	// запуск периодического получения курсов валют каждые 5 минут
	go startCurrencyFetcher(ctx, apiService, prodService)

	// запуск consumer для обработки сообщений из Kafka
	go func() {
		if err := consumerService.Consume(ctx); err != nil {
			log.Printf("consumerservice appRun: %v", err)
		}
	}()

	// запуск HTTP сервера для API
	server, err := public.NewServer(storageService)
	if err != nil {
		return errors.Wrap(err, "public.NewServer")
	}

	httpServer := &http.Server{
		Addr:    cfg.HTTPAddr, // например, ":8080"
		Handler: server,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Shutting down server...")
		cancel()
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()
		_ = httpServer.Shutdown(ctxShutdown)
	}()

	return httpServer.ListenAndServe()
}

func startCurrencyFetcher(ctx context.Context, apiService *cases.ServiceAPI, prodService *cases.Producer) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("Starting currency fetcher - will fetch rates every 5 minutes")

	for {
		select {
		case <-ctx.Done():
			log.Println("Currency fetcher stopped")
			return
		case <-ticker.C:
			log.Println("Fetching currency rates...")

			// Получаем курсы валют
			rates, err := apiService.GetRates(ctx)
			if err != nil {
				log.Printf("Failed to get rates: %v", err)
				continue
			}

			// Отправляем в Kafka
			if err := prodService.Produce(ctx, rates); err != nil {
				log.Printf("Failed to produce rates to Kafka: %v", err)
				continue
			}

			log.Printf("Successfully fetched and sent %d rates to Kafka", len(rates))
		}
	}
}
