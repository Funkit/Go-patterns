# API

## REST server with Chi router

When writing a REST API server in go using the [chi router](https://github.com/go-chi/chi), the following pattern can be used:

## Server struct

The Server struct is used to encapsulate the router and all necessary information for running the server.
The Server struct will have the following elements:

```
type Server struct {
	port   int
	router chi.Router
	c *someClient
}
```

- `port` is the port at which the REST API is exposed.
- `router` is the Chi router.
- `c` is some client used to get the data to expose (a DB connection, http client to another API...)

## Methods

### NewServer

```
func NewServer(port int) *Server {
	return &Server{
		router: chi.NewRouter(),
		port:   port,
	}
}
```

Initializes the router and various variables that allow the server to get data later (API address, DB details...). No connection is done at this stage.

### AddMiddlewares

```
func (s *Server) AddMiddlewares(middlewares ...func(handler http.Handler) http.Handler) {
	s.router.Use(middlewares...)
}
```

Middleware addition is separated from the `NewServer` or `Run` methods so that if the server is imported in a different project, the middleware can be customized.

### SubRoutes

```
func (s *Server) SubRoutes(baseURL string, r chi.Router) {
	s.router.Mount(baseURL, r)
}
```

Subroute addition is separated so that it can be extended by a separate project.

### Run

```
func (s *Server) Run() error {
	log.Printf("Listening on port %v\n", s.port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", s.port), s.router); err != nil {
		return err
	}
	return nil
}
```

The main method for starting up the server. 

### InitializeRoutes

```
func (s *Server) InitializeRoutes() {
	s.router.Get("/health", s.getSystemHealth())
}
```

This method allows all server routes to be shown in the same place. to keep it uncluttered, create methods that return handlers (see below).

### Routes methods

As we are using the `InitializeRoutes` method, the best way is to use a function that return an HTTP handler :

```
func (s *Server) getSystemHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	    //Processing, setting response body...
	}
}
```

## Main function

```
var (
	confFile = flag.String("c", "", "Path to the configuration file")
)

func main() {
	flag.Parse()

	conf, err := utils.GenericYAMLParsing[config](*confFile)
	if err != nil {
		panic(err)
	}

	server := api.NewServer(conf.ServerPort, conf.dbDetails)
	server.AddMiddlewares(middleware.Logger, render.SetContentType(render.ContentTypeJSON), middleware.Recoverer)
	server.InitializeRoutes()

	if err := server.Run(); err != nil {
		panic(err)
	}
}
```

