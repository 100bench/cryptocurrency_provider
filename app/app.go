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
	log.Printf("connection string %s", cfg.PostgresDSN())
	if err != nil {
		return errors.Wrap(err, "postgres.NewPgxClient")
	}
	defer storage.Close()

	// инициализация внешнего клиента
	client, err := coindesk.NewClientCoinDesk(cfg.CoinDeskAPIURL)
	log.Printf("coindesk client %s", cfg.CoinDeskAPIURL)

	if err != nil {
		return errors.Wrap(err, "coindesk.NewClientCoinDesk")
	}

	// инициализация брокера
	broker, err := kafka.NewBroker(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID)
	log.Printf("kafka broker %s", cfg.KafkaBrokers)
	if err != nil {
		return errors.Wrap(err, "kafka.NewBroker")
	}
	defer func() {
		_ = broker.Close()
	}()

	// конструкторы
	prodService, err := cases.NewProducer(broker)
	log.Printf("producer")
	if err != nil {
		return errors.Wrap(err, "cases.NewProducer")
	}
	storageService, err := cases.NewStorageService(storage)
	log.Printf("storage service")
	if err != nil {
		return errors.Wrap(err, "cases.NewStorageService")
	}
	consumerService, err := cases.NewConsumer(broker, storage)
	log.Printf("consumer service")
	if err != nil {
		return errors.Wrap(err, "cases.NewConsumer")
	}
	apiService, err := cases.NewServiceAPI(client, storage)
	log.Printf("api service")
	if err != nil {
		return errors.Wrap(err, "cases.NewServiceAPI")
	}

	// запуск периодического получения курсов валют с настраиваемым интервалом
	go startCurrencyFetcher(ctx, apiService, prodService, cfg.FetchInterval)

	// запуск consumer для обработки сообщений из Kafka
	go func() {
		if err := consumerService.Consume(ctx); err != nil {
			log.Printf("consumerservice appRun: %v", err)
		}
	}()

	// запуск HTTP сервера для API
	server, err := public.NewServer(storageService)
	log.Printf("server statr %v", server)
	if err != nil {
		return errors.Wrap(err, "public.NewServer")
	}

	httpServer := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: server,
	}
	log.Printf("application started")

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		log.Println("Shutting down server...")
		cancel()
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancelShutdown()
		_ = httpServer.Shutdown(ctxShutdown)
	}()

	return httpServer.ListenAndServe()
}

func startCurrencyFetcher(ctx context.Context, apiService *cases.ServiceAPI, prodService *cases.Producer, fetchInterval time.Duration) {
	ticker := time.NewTicker(fetchInterval)
	defer ticker.Stop()

	log.Printf("Starting currency fetcher - will fetch rates every %v", fetchInterval)

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
