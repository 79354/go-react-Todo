package router

import(
	"go-react-Todo/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router{
	router:= mux.NewRouter()
	router.HandleFunc("/api/task", middleware.GetAllTasks).Methods("GET", "OPTIONS")	// fetchall

	router.HandleFunc("/api/tasks", middleware.CreateTask).Methods("POST", "OPTIONS")	// create
	router.HandleFunc("/api/tasks/{id}", middleware.TaskComplete).Methods("PUT", "OPTIONS")		// complete
	router.HandleFunc("/api/undoTask/{id}", middleware.UndoTask).Methods("PUT", "OPTIONS")		// edit

	//delete
	router.HandleFunc("/api/deleteTask/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/deleteAllTasks", middleware.DeleteAllTasks).Methods("DELETE", "OPTIONS")


	return router
}