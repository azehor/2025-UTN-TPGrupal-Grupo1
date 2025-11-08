package server

import (
	//stdlib imports
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	//internal project imports
	"quepc/api/internal/carreras"
	cStore "quepc/api/internal/carreras/store"
	//"quepc/api/internal/middleware"
	"quepc/api/internal/recomendaciones"

	"quepc/api/internal/softwares"
	sStore "quepc/api/internal/softwares/store"

	carrerasoftware "quepc/api/internal/carreraSoftware"
	csStore "quepc/api/internal/carreraSoftware/store"

	transport "quepc/api/internal/transport/http"

	//imports de componentes
	"quepc/api/internal/componentes"

	"quepc/api/internal/auth"
	"quepc/api/internal/componentes/discos"
	"quepc/api/internal/componentes/gabinetes"
	"quepc/api/internal/componentes/gpus"
	"quepc/api/internal/componentes/motherboards"
	"quepc/api/internal/componentes/procesadores"
	"quepc/api/internal/componentes/psus"
	"quepc/api/internal/componentes/rams"

	dsStore "quepc/api/internal/componentes/discos/store"
	gsStore "quepc/api/internal/componentes/gabinetes/store"
	gpStore "quepc/api/internal/componentes/gpus/store"
	mbStore "quepc/api/internal/componentes/motherboards/store"
	prsStore "quepc/api/internal/componentes/procesadores/store"
	psStore "quepc/api/internal/componentes/psus/store"
	rmsStore "quepc/api/internal/componentes/rams/store"

	dbPostgresql "quepc/api/internal/db"

	//external library imports
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/docgen"
)

func Start(addr string) error {
	r := chi.NewRouter()

	//r.Use(middleware.ContentTypeMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	dbPostgresql.InitDB()

	carrerasStore := cStore.New(dbPostgresql.DB)
	carreras := carreras.New(carrerasStore)

	softwaresStore := sStore.New(dbPostgresql.DB)
	softwares := softwares.New(softwaresStore)

	csStoreInst := csStore.New(dbPostgresql.DB)
	carreraSoftwares := carrerasoftware.New(csStoreInst)

	// inicializacion de componentes
	discosStore := dsStore.New(dbPostgresql.DB)
	discos := discos.New(discosStore)
	gabinetesStore := gsStore.New(dbPostgresql.DB)
	gabinetes := gabinetes.New(gabinetesStore)
	gpusStore := gpStore.New(dbPostgresql.DB)
	gpus := gpus.New(gpusStore)
	motherboardsStore := mbStore.New(dbPostgresql.DB)
	motherboards := motherboards.New(motherboardsStore)
	procesadoresStore := prsStore.New(dbPostgresql.DB)
	procesadores := procesadores.New(procesadoresStore)
	psusStore := psStore.New(dbPostgresql.DB)
	psus := psus.New(psusStore)
	ramsStore := rmsStore.New(dbPostgresql.DB)
	rams := rams.New(ramsStore)
	componentes := componentes.New(discos,
		gabinetes,
		gpus,
		motherboards,
		procesadores,
		psus,
		rams)

	recomendaciones := recomendaciones.New(softwaresStore, csStoreInst, gpusStore, procesadoresStore, ramsStore, discosStore, gabinetesStore, motherboardsStore, psusStore)

	authHandler := auth.NewHandler(dbPostgresql.DB)

	httpServer := transport.New(carreras,
		softwares,
		recomendaciones,
		carreraSoftwares,
		componentes,
		authHandler)
	httpServer.AddRoutes(r)

	// Se busca la variable de entorno y en base a eso se genera o no la documentacion
	if os.Getenv("GENERATE_DOCS") == "1" {
		_ = os.MkdirAll("docs", 0o755)
		md := docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{Intro: "## Endpoints\n\n"})
		_ = os.WriteFile("docs/routes.md", []byte(md), 0o644)

		fmt.Println("Documentacion generada en docs/routes.md")
		return nil
	}

	if err := auth.CrearSuperAdmin(dbPostgresql.DB); err != nil {
		return fmt.Errorf("error creando superadmin: %w", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		shutdown()
		os.Exit(1)
	}()

	fmt.Printf("Running server on %s\n", addr)
	err := http.ListenAndServe(addr, r)
	shutdown()
	return err
}

func shutdown() {
	fmt.Printf("Shutting down...\n")
}
