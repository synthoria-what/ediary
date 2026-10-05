package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"synthori/ediary/m/database"
	"synthori/ediary/m/handlers"
	"synthori/ediary/m/routes"
	"synthori/ediary/m/services"

	"github.com/joho/godotenv"
)

func run() error {
	dsn := os.Getenv("DB_NAME")

	if dsn == "" {
		dsn = "ediary.db"
	}

	log.Printf("DB_NAME=%q", dsn)

	db, err := database.OpenDB(dsn)
	if err != nil {
		return fmt.Errorf("open db (dsn=%q): %w", dsn, err)
	}
	defer db.Close()
	log.Println("db opened")

	userRepo := database.NewSQLiteUserDatabase(db)

	userSvc := services.NewUserService(userRepo)
	log.Println("user svc built")

	userHandler := handlers.NewUserHandler(userSvc)
	log.Println("handler built")

	hs := handlers.Handlers{User: userHandler}

	r := routes.InitRoutes(&hs)
	log.Println("routes registered")

	log.Println("listening on :8000")
	err = http.ListenAndServe(":8000", r)
	log.Printf("ListenAndServe returned: %+v", err)
	return err
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("failed to open log file", "error", err)
		os.Exit(1)
	}
	defer file.Close()

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
			"err", err, // slog автоматически вызовет err.Error()
		)
	}
	slog.Info("exited normally")
}
