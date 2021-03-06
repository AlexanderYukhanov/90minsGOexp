// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	storage "github.com/AlexanderYukhanov/90minsGOexp/memstorage"
	"github.com/AlexanderYukhanov/90minsGOexp/service"
	"log"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/AlexanderYukhanov/90minsGOexp/server/restapi/operations"
	"github.com/AlexanderYukhanov/90minsGOexp/server/restapi/operations/trainers"
	"github.com/AlexanderYukhanov/90minsGOexp/server/restapi/operations/users"
)

//go:generate swagger generate server --target ../../server --name Experimental --spec ../../swagger/swagger.yml --principal interface{}

func configureFlags(api *operations.ExperimentalAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ExperimentalAPI) http.Handler {
	service, err := service.New(storage.NewStorageWithDefaultData())
	if err != nil {
		log.Fatalf("Failed to create a service: %v", err)
	}
	api.UsersCreateAppointmentHandler = users.CreateAppointmentHandlerFunc(service.CreateAppointment)
	api.UsersListAvailableTimesForTrainerHandler =
		users.ListAvailableTimesForTrainerHandlerFunc(service.ListAvailableTimesForTrainer)
	api.TrainersListTrainerAppointmentsHandler =
		trainers.ListTrainerAppointmentsHandlerFunc(service.ListTrainerAppointments)

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.UsersCreateAppointmentHandler == nil {
		api.UsersCreateAppointmentHandler = users.CreateAppointmentHandlerFunc(func(params users.CreateAppointmentParams) middleware.Responder {
			return middleware.NotImplemented("operation users.CreateAppointment has not yet been implemented")
		})
	}
	if api.UsersListAvailableTimesForTrainerHandler == nil {
		api.UsersListAvailableTimesForTrainerHandler = users.ListAvailableTimesForTrainerHandlerFunc(func(params users.ListAvailableTimesForTrainerParams) middleware.Responder {
			return middleware.NotImplemented("operation users.ListAvailableTimesForTrainer has not yet been implemented")
		})
	}
	if api.TrainersListTrainerAppointmentsHandler == nil {
		api.TrainersListTrainerAppointmentsHandler = trainers.ListTrainerAppointmentsHandlerFunc(func(params trainers.ListTrainerAppointmentsParams) middleware.Responder {
			return middleware.NotImplemented("operation trainers.ListTrainerAppointments has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
