package rest

import (
	"log"
	"net/http"
	"os"
	"plutus/config"
	"plutus/constants"
	gr "plutus/groups-reader"
	"plutus/redis"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/rs/cors"
)

// APIServer is the REST API application
type APIServer struct {
	router        *mux.Router
	redisClient   *redis.Client
	logger        *logrus.Logger
	UIAddress     string
	lastRefreshed time.Time
}

// NewAPIServer returns a new Server struct object
func NewAPIServer(redisClient *redis.Client, groupsReader gr.GroupsReader, logger *logrus.Logger) (*APIServer, error) {
	router := mux.NewRouter()

	uiAdress, err := config.GetUIAddress()
	if err != nil {
		return nil, err
	}
	server := &APIServer{
		router:        router,
		redisClient:   redisClient,
		logger:        logger,
		UIAddress:     uiAdress,
		lastRefreshed: time.Now(),
	}
	server.initRoutes()
	return server, nil
}

// initRoutes initialises the routes in the application
func (a *APIServer) initRoutes() {
	a.router.HandleFunc("/ui", a.uiHandler)
	a.router.HandleFunc("/api/v1/path", a.pathHandler).Methods("POST")
	a.router.HandleFunc("/api/v1/user", a.userHandler).Methods("POST")
	a.router.HandleFunc("/api/v1/namespaces", a.namespacesHandler).Methods("GET")
	a.router.HandleFunc("/api/v1/lastRefresh", a.lastRefreshHandler).Methods("GET")
}

// Run starts the rest api server
func (a *APIServer) Run(refresh func(redisClient *redis.Client, logger *logrus.Logger) error) {
	restAddr := constants.ENV_REST_ADDR
	addr, ok := os.LookupEnv(restAddr)
	if !ok {
		log.Fatalf("%s not specified", restAddr)
	}

	quit := make(chan struct{})

	// refresh every 30 minutes
	ticker := time.NewTicker(30 * 60 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				refresh(a.redisClient, a.logger)
				a.lastRefreshed = time.Now()
			case <-quit:
				ticker.Stop()
			}
		}
	}()

	a.logger.Info("Starting REST API server at ", addr)
	c := cors.Default()
	routerWithCORS := c.Handler(a.router)

	log.Fatal(http.ListenAndServe(addr, routerWithCORS))
}
