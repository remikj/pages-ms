package server

import (
	"fmt"
	"github.com/remikj/pages-ms/src/contoller"
	"net/http"
)
import "github.com/go-chi/chi/v5"
import "github.com/kelseyhightower/envconfig"

type Server struct {
	Config         *Configuration
	PageController contoller.PageController
}

type Configuration struct {
	Port int `envconfig:"SERVICE_PORT" default:"8080"`
}

func NewServerFromEnv(pageController contoller.PageController) (*Server, error) {
	configFromEnv, err := ConfigurationFromEnv()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return NewServer(configFromEnv, pageController), nil
}

func ConfigurationFromEnv() (*Configuration, error) {
	config := &Configuration{}
	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func NewServer(config *Configuration, pageController contoller.PageController) *Server {
	return &Server{
		Config:         config,
		PageController: pageController,
	}
}

func (s *Server) Run() error {
	router := chi.NewRouter()
	router.Get("/pages/{id}", s.PageController.HandlePageGet)
	server := &http.Server{Addr: fmt.Sprintf(":%v", s.Config.Port), Handler: router}
	fmt.Printf("Starting server on port: %v\n", s.Config.Port)
	return server.ListenAndServe()
}
