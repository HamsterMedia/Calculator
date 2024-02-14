package main
import (
	//"context"
	"fmt"
	"log"
	"net/http"
	//"strconv"
	"sync"
	"json"
	"github.com/gorilla/mux"
)
type Operator string
const (
	Add Operator = "+"
	Subtract Operator = "-"
	Multiply Operator = "*"
	Divide Operator = "/"
)
type Calculator struct {
	mux *mux.Router
	wg  sync.WaitGroup
}
func NewCalculator() *Calculator {
	c := &Calculator{
		mux: mux.NewRouter(),
	}
	c.mux.HandleFunc("/add", c.addHandler).Methods("POST")
	c.mux.HandleFunc("/subtract", c.subtractHandler).Methods("POST")
	c.mux.HandleFunc("/multiply", c.multiplyHandler).Methods("POST")
	c.mux.HandleFunc("/divide", c.divideHandler).Methods("POST")
	return c
}
func (c *Calculator) Run(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: c.mux,
	}
	log.Printf("Calculator microservice listening on %s", addr)
	c.wg.Add(4)
	go c.addWorker()
	go c.subtractWorker()
	go c.multiplyWorker()
	go c.divideWorker()
	defer c.wg.Wait()
	return srv.ListenAndServe()
}
func (c *Calculator) Stop() {
	if err := http.DefaultClient.Get("http://localhost:8080/shutdown"); err != nil {
		log.Printf("Error stopping Calculator microservice: %v", err)
	}
	c.wg.Wait()
}
func (c *Calculator) addHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	c.add <- req
	res := <-c.addRes
	if res.Err != nil {
		http.Error(w, res.Err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%d", res.Result)
}
func (c *Calculator) addWorker() {
	defer c.wg.Done()

	for req := range c.add {
		res := &result{
			Result: req.A + req.B,
		}
		c.addRes <- res
	}
}
func (c *Calculator) subtractHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	c.subtract <- req
	res := <-c.subtractRes
	if res.Err != nil {
		http.Error(w, res.Err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%d", res.Result)
}
func (c *Calculator) subtractWorker() {
	defer c.wg.Done()
	for req := range c.subtract {
		res := &result{
			Result: req.A - req.B,
		}
		c.subtractRes <- res
	}
}
func (c *Calculator) multiplyHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	c.multiply <- req
	res := <-c.multiplyRes
	if res.Err != nil {
		http.Error(w, res.Err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%d", res.Result)
}
func (c *Calculator) multiplyWorker() {
	defer c.wg.Done()

	for req := range c.multiply {
		res := &result{
			Result: req.A * req.B,
		}
		c.multiplyRes <- res
	}
}
func (c *Calculator) divideHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	c.divide <- req
	res := <-c.divideRes
	if res.Err != nil {
		http.Error(w, res.Err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%d", res.Result)
}
func (c *Calculator) divideWorker() {
	defer c.wg.Done()

	for req := range c.divide {
		if req.B == 0 {
			res := &result{
				Err: fmt.Errorf("division by zero"),
			}

			c.divideRes <- res
			continue
		}

		res := &result{
			Result: req.A / req.B,
		}

		c.divideRes <- res
	}
}
type result struct {
	Result int
	Err    error
}
var add = make(chan struct {
	A int
	B int
})
var addRes = make(chan *result)
var subtract = make(chan struct {
	A int
	B int
})
var subtractRes = make(chan *result)
var multiply = make(chan struct {
	A int
	B int
})
var multiplyRes = make(chan *result)
var divide = make(chan struct {
	A int
	B int
})
var divideRes = make(chan *result)