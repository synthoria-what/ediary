package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"synthori/ediary/m/database"
	"synthori/ediary/m/handlers"
	"synthori/ediary/m/routes"
	"synthori/ediary/m/services"
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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("[ediary] ")

	err := run()
	if err != nil {
		// Печатаем всё, что можно: тип, значение, длину, байты
		log.Printf("run returned error")
		log.Printf("  type:   %T", err)
		log.Printf("  value:  %+v", err)
		log.Printf("  quoted: %q", err.Error())
		log.Printf("  len:    %d", len(err.Error()))
		log.Fatalf("fatal: %v", err)
	}
	log.Println("exited normally")
}
