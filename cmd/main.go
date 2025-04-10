package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/og11423074s/gocourse_enrollment/pkg/handler"
	"os"

	"github.com/og11423074s/gocourse_enrollment/internal/enrollment"
	"github.com/og11423074s/gocourse_enrollment/pkg/bootstrap"
	"net/http"
	"time"
)

func main() {

	// Load .env file
	_ = godotenv.Load()

	// Initialize logger
	logger := bootstrap.InitLogger()

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		logger.Fatal("PAGINATION_LIMIT_DEFAULT is not set")
	}

	// Connect to database
	db, err := bootstrap.DBConnection()

	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	// Enroll repository
	enrollRepo := enrollment.NewRepo(logger, db)

	// Enroll service
	enrollSrv := enrollment.NewService(logger, enrollRepo)

	// Enroll endpoints
	h := handler.NewEnrollmentHandler(ctx, enrollment.MakeEndpoints(enrollSrv, enrollment.Config{LimPageDef: pagLimDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 4 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		logger.Println("listen in ", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		logger.Fatal(err)
	}

}

func accessControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
