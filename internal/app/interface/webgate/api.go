package webgate

import (
	"github.com/globbie/knowdy-media/internal/app/usecases/upload"
	"github.com/globbie/knowdy-media/internal/app/interface/monitor"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"encoding/json"
	"net/http"
)

const (
	userIdContextKey = "user_id"
)

type WebGate struct {
	UploadFileSaver upload.FileSaver
}

func New(fs upload.FileSaver) *WebGate {
	return &WebGate{
		UploadFileSaver: fs,
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
type FileUploadResponse struct {
	ErrCode          string   `json:"errorCode"`
	ErrDescription   string   `json:"errorDesc"`
	ContentDispositionList     []ContentDisposition `json:"contentDispositionList"`
}

type ContentDisposition struct {
	MimeType         string   `json:"mimeType"`
	Name             string   `json:"name"`
	FileId           string   `json:"fileId"`
	FileName         string   `json:"fileName"`
	FileSize         uint64   `json:"fileSize"`
	ErrCode          string   `json:"errorCode"`
	ErrDescription   string   `json:"errorDesc"`
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

	// accumulate upload results
	result := FileUploadResponse{
		ErrCode: "OK",
		ErrDescription:    "Upload successful",
		ContentDispositionList:  []ContentDisposition{},
	}
	w.Header().Add("Content-Type", "application/json")

	var cds = make([]ContentDisposition, 0, 1)
        for k, _ := range r.MultipartForm.File {
	        file, fh, err := r.FormFile(k)
		if err != nil {
			result.ErrCode = "Incorrect form data"
			result.ErrDescription = "Please check the input form"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		var mimetype = "application/octet-stream"
		types, ok := fh.Header["Content-Type"]
		if ok {
			for _, t := range types {
				mimetype = http.DetectContentType([]byte(t))
				break
			}
		}

		var cd = ContentDisposition{mimetype, k,
			"0", fh.Filename, uint64(fh.Size), "OK", "OK"}

                rec, uc_err := wg.UploadFileSaver.SaveFile(file,
			mimetype, fh.Filename, uint64(fh.Size), "me") // TODO auth
	        if uc_err != nil {
			// TODO: shall we continue or not?

			// return current status
			if err = json.NewEncoder(w).Encode(result); err != nil {
				
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(wg.ConvertToHTTPStatus(uc_err))
			return
	        }

		cd.FileId = rec.Id
		cds = append(cds, cd)
	}

        // positive response
        w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
