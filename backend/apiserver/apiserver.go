package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/ITegs/crs.pics/cloudprovider"
	"github.com/ITegs/crs.pics/database"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"net/http"
)

type ApiServer interface {
	Serve()
}

type apiServer struct {
	db database.DB
	cp cloudprovider.CP
}

func NewApiServer(DB database.DB, CP cloudprovider.CP) ApiServer {
	apiServer := &apiServer{
		db: DB,
		cp: CP,
	}

	return apiServer
}

func (api *apiServer) Serve() {
	fmt.Println("Program started!")

	apiHandler := api.buildApi()

	server := http.Server{
		Addr:    ":3000",
		Handler: apiHandler,
	}

	fmt.Printf("API server listening on %s\n", server.Addr)

	handler := cors.Default().Handler(server.Handler)
	err := http.ListenAndServe(server.Addr, handler)
	if err != nil {
		fmt.Println("SERVER FAILED: ", err)
	}
}

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

func (api *apiServer) buildApi() *httprouter.Router {
	router := httprouter.New()

	var routes = []*Route{
		{
			Method: http.MethodGet,
			Path:   "/",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("API is up and running!"))
				w.WriteHeader(http.StatusOK)
			}),
		},
		{
			Method:  http.MethodGet,
			Path:    "/links",
			Handler: http.HandlerFunc(api.GetAllLinks),
		},
		{
			Method:  http.MethodGet,
			Path:    "/link/:slug",
			Handler: http.HandlerFunc(api.GetLink),
		},
		{
			Method:  http.MethodPost,
			Path:    "/newLink",
			Handler: http.HandlerFunc(api.AddLink),
		},
		{
			Method:  http.MethodGet,
			Path:    "/getFileNames",
			Handler: http.HandlerFunc(api.GetFileNames),
		},
	}

	for i := 0; i < len(routes); i++ {
		r := routes[i]
		router.Handler(r.Method, "/api"+r.Path, r.Handler)
	}

	return router
}

func (api *apiServer) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	links := api.db.GetAllLinks()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(links)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	return
}

func (api *apiServer) GetLink(w http.ResponseWriter, r *http.Request) {
	p := httprouter.ParamsFromContext(r.Context())
	slug := p.ByName("slug")

	data := api.db.GetLinkBySlug(slug)

	w.Write([]byte(data))
	w.WriteHeader(http.StatusOK)
}

func (api *apiServer) AddLink(w http.ResponseWriter, r *http.Request) {
	var link database.Link
	json.NewDecoder(r.Body).Decode(&link)

	createdLink := api.db.AddLink(link)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&createdLink)
	w.WriteHeader(http.StatusCreated)
}

func (api *apiServer) GetFileNames(w http.ResponseWriter, r *http.Request) {
	fileNames := api.cp.GetAllFileNames()

	json.NewEncoder(w).Encode(&fileNames)
	w.WriteHeader(http.StatusOK)
}
