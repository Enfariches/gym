package controllers

import (
	"encoding/json"
	"health/src/health/models"
	u "health/src/health/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func GetVideosListHandler(w http.ResponseWriter, r *http.Request) {
	data := models.GetVideosList()
	u.Respond(w, data)
}

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, "Something went wrong!"))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, "Invalid request body"))
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			u.Respond(w, u.Message(u.InternalError, "Something went wrong!"))
		}
	}(file)

	vid, _ := models.AddVideo(handler.Filename)

	if err := u.SaveUploadedFile(handler, "upload/"+strconv.Itoa(int(vid.ID))+filepath.Ext(handler.Filename)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Respond(w, u.Message(u.UploadError, "Something went wrong!"))
		return
	}
	u.Respond(w, vid)
}

func UpdateVideoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidRequest, "Invalid request"))
		return
	}
	var req struct {
		Name string
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, "Invalid request body"))
		return
	}
	vid, err := models.UpdVid(uint(id), req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Respond(w, u.Message(u.InternalError, "Something went wrong!"))
		return
	}
	u.Respond(w, vid)
}

func DeleteVideoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidRequest, "Invalid request"))
		return
	}
	models.DeleteVideo(uint(id))
	u.Respond(w, u.Message(u.Ok, "Video successfully removed"))
}

func ExportVideo(w http.ResponseWriter, r *http.Request) {
	videoID := mux.Vars(r)["id"]
	filePath := "upload/" + videoID + ".mp4" // Путь к видео

	// Чтение файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Disposition", "attachment; filename="+videoID+".mp4")
	w.Write(data)
}

func ToggleVideoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidRequest, "Invalid request"))
		return
	}
	err = models.ToggleVid(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Respond(w, u.Message(u.InternalError, "Something went wrong!"))
		return
	}
	u.Respond(w, u.Message(u.Ok, "Video successfully toggled"))
}
