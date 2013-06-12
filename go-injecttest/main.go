package main

import (
	"code.google.com/p/go-inject"
	"code.google.com/p/go-inject/http"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Counter struct{}

func sayFoo(w http.ResponseWriter, request *http.Request) {
	w.Header().Add(
		"Content-Type",
		"text/plain",
	)
	w.Write([]byte(fmt.Sprintf("Foo\n")))
}

func ConfigureInjector(injector inject.Injector) {
	i := 0
	injector.BindInScope(Counter{}, func(_ inject.Context, _ inject.Container) interface{} {
		i += 1
		return i
	}, inject_http.RequestScoped{})

	inject_http.BindHandlerFunc(injector, "/",
		func(w http.ResponseWriter, request *http.Request) {
			w.Header().Add(
				"Content-Type",
				"text/plain",
			)
			w.Write([]byte(fmt.Sprintf("Hello! (%d)\n",
				injector.CreateContainer().GetInstance(request, Counter{}).(int))))
			w.Write([]byte(fmt.Sprintf("Hello! (%d)\n",
				injector.CreateContainer().GetInstance(request, Counter{}).(int))))
		})
}

func ConfigureFooInjector(injector inject.Injector) {
	inject_http.BindHandlerFunc(injector, "/foo/", sayFoo)
}

func main() {
	// Initialize the flag(s)
	flag.Parse()

	// Create the injector to begin configuring the bindings.
	injector := inject.CreateInjector()

	// Bind any flags used in the inject_http module. Flag bindings are in a separate function
	// so that the module can be used without flags, with explicit bindings (see multi_http.go).
	inject_http.ConfigureFlags(injector)

	// Bind any scope tags used in the inject_http module. Scope bindings are in a separate function
	// so the module can be used more than once. The scope binding can only ever happen once.
	inject_http.ConfigureScopes(injector)

	// Now configure the bindings in the inject_http module. This step can be performed more than
	// once in multiple child injectors.
	inject_http.ConfigureInjector(injector)

	// There are multiple ConfigureInjector functions to demonstrate that multiple modules
	// can configure HTTP mappings. If two modules try to bind the same path, the system will panic.
	ConfigureFooInjector(injector)
	ConfigureInjector(injector)

	// Look up an object by creating a Container. The Container ensures there are no dependency
	// loops configured in the Injector. Top-level code should create a new container for each
	// object lookup. Provider functions will have a container passed to them, and they use that
	// container for their lookups.
	container := injector.CreateContainer()
	httpServer := container.GetInstance(nil, inject_http.Server{}).(http.Server)

	// Start the HTTP server.
	log.Fatal(httpServer.ListenAndServe())
}
