package router

import (
	"fmt"
	"net/http"

	"github.com/ds3lab/easeml/engine/api"
	"github.com/ds3lab/easeml/engine/api/handlers"
	"github.com/ds3lab/easeml/engine/api/middleware"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

// Route is a single route descriptor.
type Route struct {
	Name     string
	Methods  []string
	Pattern  string
	IsPrefix bool
	Handler  http.Handler
}

// Routes is a list of routes.
type Routes []Route

// New initializes a new gorilla/mux router.
func New(context api.Context) http.Handler {

	middlewareContext := middleware.Context(context)
	handlerContext := handlers.Context(context)

	// Set up CORS.
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		Debug:            false,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "HEAD"},
	})

	var commonMiddleware = alice.New(
		middlewareContext.Logging,
		middlewareContext.PanicRecovery,
		middlewareContext.RequestID,
		middlewareContext.Inject,
		//storageContext.Inject,
		middlewareContext.Authenticate,
	)

	var routes = Routes{
		Route{
			Name:    "Index",
			Methods: []string{"GET"},
			Pattern: "/",
			Handler: commonMiddleware.ThenFunc(Index),
		},
		Route{
			Name:    "GetUsers",
			Methods: []string{"GET"},
			Pattern: "/users",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.UsersGet),
		},
		Route{
			Name:    "PostUser",
			Methods: []string{"POST"},
			Pattern: "/users",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.UsersPost),
		},
		Route{
			Name:    "LoginUser",
			Methods: []string{"GET"},
			Pattern: "/users/login",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.UsersLoginGet),
		},
		Route{
			Name:    "LogoutUser",
			Methods: []string{"GET"},
			Pattern: "/users/logout",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.UsersLogoutGet),
		},
		Route{
			Name:    "GetUser",
			Methods: []string{"GET"},
			Pattern: "/users/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.UsersByIDGet),
		},
		Route{
			Name:    "PatchUser",
			Methods: []string{"PATCH"},
			Pattern: "/users/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.UsersByIDPatch),
		},
		Route{
			Name:    "GetProcesses",
			Methods: []string{"GET"},
			Pattern: "/processes",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.ProcessesGet),
		},
		Route{
			Name:    "GetProcess",
			Methods: []string{"GET"},
			Pattern: "/processes/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.ProcesssesByIDGet),
		},
		Route{
			Name:    "GetDatasets",
			Methods: []string{"GET"},
			Pattern: "/datasets",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.DatasetsGet),
		},
		Route{
			Name:    "PostDatasets",
			Methods: []string{"POST"},
			Pattern: "/datasets",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.DatasetsPost),
		},
		Route{
			Name:    "GetDataset",
			Methods: []string{"GET"},
			Pattern: "/datasets/{user-id}/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.DatasetsByIDGet),
		},
		Route{
			Name:    "PatchDataset",
			Methods: []string{"PATCH"},
			Pattern: "/datasets/{user-id}/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.DatasetsByIDPatch),
		},
		Route{
			Name:     "UploadDataset",
			Methods:  []string{"POST"},
			Pattern:  "/datasets/{user-id}/{dataset-id}/upload",
			IsPrefix: false,
			Handler:  commonMiddleware.ThenFunc(handlerContext.DatasetsUploadHandler("/api/v1/datasets/{user-id}/{dataset-id}/upload")),
		},
		Route{
			Name:     "HeadPatchDataset",
			Methods:  []string{"HEAD", "PATCH"},
			Pattern:  "/datasets/{user-id}/{dataset-id}/upload/{upload-id}",
			IsPrefix: false,
			Handler:  commonMiddleware.ThenFunc(handlerContext.DatasetsUploadHandler("/api/v1/datasets/{user-id}/{dataset-id}/upload/")),
		},
		Route{
			Name:     "GetDatasets",
			Methods:  []string{"GET"},
			Pattern:  "/datasets/{user-id}/{dataset-id}/data",
			IsPrefix: true,
			Handler:  commonMiddleware.ThenFunc(handlerContext.DatasetsDownloadHandler("/api/v1/datasets/{user-id}/{dataset-id}/data")),
		},
		Route{
			Name:    "GetModules",
			Methods: []string{"GET"},
			Pattern: "/modules",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.ModulesGet),
		},
		Route{
			Name:    "PostModules",
			Methods: []string{"POST"},
			Pattern: "/modules",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.ModulesPost),
		},
		Route{
			Name:    "GetModule",
			Methods: []string{"GET"},
			Pattern: "/modules/{user-id}/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.ModulesByIDGet),
		},
		Route{
			Name:    "PatchModule",
			Methods: []string{"PATCH"},
			Pattern: "/modules/{user-id}/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.ModulesByIDPatch),
		},
		Route{
			Name:     "PostModule",
			Methods:  []string{"POST"},
			Pattern:  "/modules/{user-id}/{module-id}/upload",
			IsPrefix: false,
			Handler:  commonMiddleware.ThenFunc(handlerContext.ModulesUploadHandler("/api/v1/modules/{user-id}/{module-id}/upload")),
		},
		Route{
			Name:     "HeadPatchModule",
			Methods:  []string{"HEAD", "PATCH"},
			Pattern:  "/modules/{user-id}/{module-id}/upload/{upload-id}",
			IsPrefix: false,
			Handler:  commonMiddleware.ThenFunc(handlerContext.ModulesUploadHandler("/api/v1/modules/{user-id}/{module-id}/upload/")),
		},
		Route{
			Name:    "GetJobs",
			Methods: []string{"GET"},
			Pattern: "/jobs",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.JobsGet),
		},
		Route{
			Name:    "PostJobs",
			Methods: []string{"POST"},
			Pattern: "/jobs",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.JobsPost),
		},
		Route{
			Name:    "GetJobs",
			Methods: []string{"GET"},
			Pattern: "/jobs/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.JobsByIDGet),
		},
		Route{
			Name:    "PatchJob",
			Methods: []string{"PATCH"},
			Pattern: "/jobs/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.JobsByIDPatch),
		},
		Route{
			Name:    "GetTasks",
			Methods: []string{"GET"},
			Pattern: "/tasks",
			Handler: commonMiddleware.Append(middlewareContext.DisallowAnon).ThenFunc(handlerContext.TasksGet),
		},
		Route{
			Name:    "GetTask",
			Methods: []string{"GET"},
			Pattern: "/tasks/{job-id}/{id}",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.TasksByIDGet),
		},
		Route{
			Name:     "GetTaskPredictions",
			Methods:  []string{"GET"},
			Pattern:  "/tasks/{job-id}/{task-id}/predictions",
			IsPrefix: true,
			Handler:  commonMiddleware.ThenFunc(handlerContext.TaskDataDownloadHandler("/api/v1/tasks/{job-id}/{task-id}/predictions", "predictions")),
		},
		Route{
			Name:     "GetTaskLogs",
			Methods:  []string{"GET"},
			Pattern:  "/tasks/{job-id}/{task-id}/logs",
			IsPrefix: true,
			Handler:  commonMiddleware.ThenFunc(handlerContext.TaskDataDownloadHandler("/api/v1/tasks/{job-id}/{task-id}/logs", "logs")),
		},
		Route{
			Name:     "GetTaskMetadata",
			Methods:  []string{"GET"},
			Pattern:  "/tasks/{job-id}/{task-id}/metadata",
			IsPrefix: true,
			Handler:  commonMiddleware.ThenFunc(handlerContext.TaskDataDownloadHandler("/api/v1/tasks/{job-id}/{task-id}/metadata", "metadata")),
		},
		Route{
			Name:     "GetTaskParameters",
			Methods:  []string{"GET"},
			Pattern:  "/tasks/{job-id}/{task-id}/parameters",
			IsPrefix: true,
			Handler:  commonMiddleware.ThenFunc(handlerContext.TaskDataDownloadHandler("/api/v1/tasks/{job-id}/{task-id}/parameters", "parameters")),
		},
		Route{
			Name:    "GetTaskImage",
			Methods: []string{"GET"},
			Pattern: "/tasks/{job-id}/{id}/image/download",
			Handler: commonMiddleware.Append(middlewareContext.HideFromAnon).ThenFunc(handlerContext.TaskImageDownload),
		},
	}

	router := mux.NewRouter().StrictSlash(true).PathPrefix("/api/v1").Subrouter()
	for _, route := range routes {
		var handler http.Handler
		handler = route.Handler
		//handler = Logger(handler, route.Name)

		r := router.Methods(route.Methods...)
		if route.IsPrefix {
			r = r.PathPrefix(route.Pattern)
		} else {
			r = r.Path(route.Pattern)
		}
		r = r.Name(route.Name).Handler(handler)
	}
	//router.Use(c.Handler)

	return c.Handler(router)
}

// Index defines the response for the root GET request.
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Easeml API v1 root!")
}
