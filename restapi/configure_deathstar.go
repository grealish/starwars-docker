package restapi

import (
	"crypto/tls"
	"net/http"
	"time"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/tgraf/starwars/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../swagger.yaml

var backtrace = `Panic: deathstar exploded

goroutine 1 [running]:
main.HandleGarbage(0x2080c3f50, 0x2, 0x4, 0x425c0, 0x5, 0xa)
        /code/src/github.com/empire/deathstar/
        temp/main.go:9 +0x64
main.main()
        /code/src/github.com/empire/deathstar/
        temp/main.go:5 +0x85
`

var info = `{
	"name": "Death Star",
	"model": "DS-1 Orbital Battle Station",
	"manufacturer": "Imperial Department of Military Research, Sienar Fleet Systems",
	"cost_in_credits": "1000000000000",
	"length": "120000",
	"crew": "342953",
	"passengers": "843342",
	"cargo_capacity": "1000000000000",
	"hyperdrive_rating": "4.0",
	"starship_class": "Deep Space Mobile Battlestation",
	"api": [
		"GET   /v1",
		"POST  /v1/requestlanding",
		"GET   /v1/vader/secrets",
		"PUT   /v1/exhaustport"
	]
}
`

var trap = `
  ___ _   _               _
 |_ _| |_( )___    __ _  | |_ _ __ __ _ _ __
  | || __|// __|  / _' | | __| '__/ _' | '_ \
  | || |_  \__ \ | (_| | | |_| | | (_| | |_) |
 |___|\__| |___/  \__,_|  \__|_|  \__,_| .__/
                                       |_|
`

func configureFlags(api *operations.DeathstarAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.DeathstarAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.TxtProducer = runtime.TextProducer()

	api.GetHandler = operations.GetHandlerFunc(func(params operations.GetParams) middleware.Responder {
		return operations.NewGetOK().WithPayload(info)
	})
	api.PutExhaustportHandler = operations.PutExhaustportHandlerFunc(func(params operations.PutExhaustportParams) middleware.Responder {
		go func() {
			time.Sleep(2 * time.Second)
			panic("deathstar exploded")
		}()
		return operations.NewPutExhaustportServiceUnavailable().WithPayload(backtrace)
	})
	api.PostRequestlandingHandler = operations.PostRequestlandingHandlerFunc(func(params operations.PostRequestlandingParams) middleware.Responder {
		return operations.NewPostRequestlandingOK().WithPayload("Ship landed\n")
	})

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
func configureServer(s *graceful.Server, scheme string) {
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
