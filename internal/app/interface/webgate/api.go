package webgate

import (
	"github.com/globbie/knowdy-media/internal/app/usecases/upload"
	"github.com/globbie/knowdy-media/internal/app/interface/monitor"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// "mime/multipart"
	"encoding/json"
	"net/http"
	"log"
)

const (
	userIdContextKey = "user_id"
)

type WebGate struct {
	UploadUseCase upload.Interface
}

func New(u upload.Interface) *WebGate {
	return &WebGate{
		UploadUseCase: u,
	}
}

func (wg *WebGate) Router() http.Handler {
	router := mux.NewRouter()

	// router.HandleFunc("/media", wg.authenticate(wg.getMediaList)).Methods(http.MethodGet)
	router.HandleFunc("/media", wg.authenticate(wg.postMedia)).Methods(http.MethodPost)

	router.Handle("/metrics", promhttp.Handler())

	router.Use(monitor.Measurer())
	router.Use(wg.logger)

	return router
}

// Response JSON message
type MediaUploadResponse struct {
	ErrDescription   string   `json:"errorDesc"`
	MediaList    []string `json:"uploadMediaList"`
}

func (wg *WebGate) postMedia(w http.ResponseWriter, r *http.Request) {
	/*userId, ok := r.Context().Value(userIdContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}*/
        err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

        for k, _ := range r.MultipartForm.File {
	        log.Println(k)

                _, err := wg.UploadUseCase.CreateFile("text/plain", "test.txt", 8, "me", "/")
	        if err != nil {
	 	   w.WriteHeader(http.StatusInternalServerError)
		   return
	        }
	}

        w.Header().Add("Content-Type", "application/json")
        result := MediaUploadResponse{
		  ErrDescription:    "OK",
		  MediaList:    []string{},
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
