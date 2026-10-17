package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"synthori/ediary/m/database"
	"synthori/ediary/m/handlers"
	"synthori/ediary/m/routes"
	"synthori/ediary/m/services"
	"time"

	"github.com/joho/godotenv"
)

func run() error {
	dsn := os.Getenv("DB_NAME")

	if dsn == "" {
		dsn = "ediary.db"
	}

	db, err := database.OpenDB(dsn)
	if err != nil {
		return fmt.Errorf("open db (dsn=%q): %w", dsn, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Error open db")
			return
		}
	}()
	log.Println("db opened")

	userRepo := database.NewSQLiteUserDatabase(db)
	userSvc := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	diaryRepo := database.NewSQLiteDiaryDatabase(db)
	diarySvc := services.NewDiaryService(*diaryRepo)
	diaryHandler := handlers.NewDiaryHandler(*diarySvc)

	hs := handlers.Handlers{User: userHandler,
		Diary: diaryHandler}

	r := routes.InitRoutes(&hs)
	log.Println("routes registered")

	port, exists := os.LookupEnv("DEFAULT_PORT")
	if !exists {
		return errors.New("ports didn't found")
	}

	log.Printf("listening on :%s", port)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           r,
	}
	err = server.ListenAndServe()
	log.Printf("ListenAndServe returned: %+v", err)
	return err
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		slog.Error("failed to open log file", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("close file: %v", err)
		}
	}()

	logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	err = run()
	if err != nil {
		slog.Error("run execution failed",
			"error_type", fmt.Sprintf("%T", err),
			"error_value", fmt.Sprintf("%+v", err),
			"error_quoted", fmt.Sprintf("%q", err.Error()),
			"error_len", len(err.Error()),
			"err", err,
		)
	}
	slog.Info("exited normally")
}
