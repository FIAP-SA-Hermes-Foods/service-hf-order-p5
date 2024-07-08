package main

import (
	"context"
	"log"
	"net/http"
	"os"
	l "service-hf-order-p5/external/logger"
	orderrpc "service-hf-order-p5/internal/adapters/rpc"
	"service-hf-order-p5/internal/core/application"
	httpH "service-hf-order-p5/internal/handler/http"

	"github.com/marcos-dev88/genv"
)

func init() {
	if err := genv.New(); err != nil {
		l.Errorf("", "error set envs %v", " | ", err)
	}
}

func main() {

	router := http.NewServeMux()

	ctx := context.Background()

	orderRPC := orderrpc.NewOrderRPC(ctx, os.Getenv("HOST_ORDER"), os.Getenv("PORT_ORDER"))

	orderWorkerRPC := orderrpc.NewOrderWorkerRPC(ctx, os.Getenv("HOST_ORDER"), os.Getenv("PORT_ORDER"))

	app := application.NewApplication(ctx, orderRPC, orderWorkerRPC)

	h := httpH.NewHandler(app)

	router.Handle("/hermes_foods/order/", http.StripPrefix("/", httpH.Middleware(h.HandlerOrder)))
	router.Handle("/hermes_foods/order", http.StripPrefix("/", httpH.Middleware(h.HandlerOrder)))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("API_HTTP_PORT"), router))
}
