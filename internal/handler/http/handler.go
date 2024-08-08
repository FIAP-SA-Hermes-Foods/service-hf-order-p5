package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	l "service-hf-order-p5/external/logger"
	ps "service-hf-order-p5/external/strings"
	"service-hf-order-p5/internal/core/application"
	"service-hf-order-p5/internal/core/domain/entity/dto"
	"strconv"
)

type Handler interface {
	Handler(rw http.ResponseWriter, req *http.Request)
	HealthCheck(rw http.ResponseWriter, req *http.Request)
}

type handler struct {
	app application.Application
}



func NewHandler(app application.Application) handler {
	return handler{app: app}
}

func (h handler) HandlerProduct(rw http.ResponseWriter, req *http.Request) {

	var routes = map[string]http.HandlerFunc{
		"get hermes_foods/order":        h.getOrders,
		"get hermes_foods/order/{id}":   h.getOrderByID,
		"post hermes_foods/order":       h.saveOrder,
		"patch hermes_foods/order/{id}": h.updateOrderByID,
	}

	handler, err := router(req.Method, req.URL.Path, routes)

	if err == nil {
		handler(rw, req)
		return
	}

	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(`{"error": "route ` + req.Method + " " + req.URL.Path + ` not found"} `))
}

func (h handler) HealthCheck(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte(`{"error": "method not allowed"}`))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status": "OK"}`))
}

func (h *handler) saveOrder(rw http.ResponseWriter, req *http.Request) {
	msgID := l.MessageID(req.Header.Get(l.MessageIDKey))

	if req.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte(`{"error": "method not allowed"} `))
		return
	}

	var buff bytes.Buffer
	var reqOrder dto.RequestOrder

	if _, err := buff.ReadFrom(req.Body); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to read data body: %v"} `, err)
		return
	}

	if err := json.Unmarshal(buff.Bytes(), &reqOrder); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to Unmarshal: %v"} `, err)
		return
	}

	o, err := h.app.SaveOrder(msgID, reqOrder)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to save order: %v"} `, err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(ps.MarshalString(o)))
}

func (h *handler) getOrderByID(rw http.ResponseWriter, req *http.Request) {
	msgID := l.MessageID(req.Header.Get(l.MessageIDKey))
	id := getID("order", req.URL.Path)
	idconv, err := strconv.ParseInt(id, 10, 64)

	// var buff bytes.Buffer
	// var reqOrder dto.RequestOrder

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to get order by ID: %v"} `, err)
		return
	}

	o, err := h.app.GetOrderByID(msgID, strconv.FormatInt(idconv, 10))

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to get order by ID: %v"} `, err)
		return
	}

	if o == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"error": "order not found"}`))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(ps.MarshalString(o)))
}

func (h *handler) getOrders(rw http.ResponseWriter, req *http.Request) {
	oList, err := h.app.GetOrders("", "")
	//msgID := l.MessageID(req.Header.Get(l.MessageIDKey))

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to get order by ID: %v"} `, err)
		return
	}

	if oList == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"error": "order not found"}`))
		return
	}

	b, err := json.Marshal(oList)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to get order by ID: %v"} `, err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}

func (h *handler) updateOrderByID(rw http.ResponseWriter, req *http.Request) {
	id := getID("order", req.URL.Path)
	msgID := l.MessageID(req.Header.Get(l.MessageIDKey))
	idconv, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to get order by ID: %v"} `, err)
		return
	}

	var buff bytes.Buffer

	var reqOrder dto.RequestOrder

	if _, err := buff.ReadFrom(req.Body); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to read data body: %v"} `, err)
		return
	}

	if err := json.Unmarshal(buff.Bytes(), &reqOrder); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to Unmarshal: %v"} `, err)
		return
	}

	o, err := h.app.UpdateOrderByID(msgID, strconv.FormatInt(idconv, 10), reqOrder)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, `{"error": "error to get order by ID: %v"} `, err)
		return
	}

	if o == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"error": "order not found"}`))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(ps.MarshalString(o)))
}
