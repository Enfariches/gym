package controllers

import (
	"encoding/json"
	"health/src/health/models"
	u "health/src/health/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ListAllScheduleHandler(w http.ResponseWriter, r *http.Request) {
	data := models.ListAllSchedule()
	u.Respond(w, data)
}

func ListScheduleHandler(w http.ResponseWriter, r *http.Request) {
	dayStr := mux.Vars(r)["day"] //Берется из значения в URL
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidRequest, "Invalid request"))
		return
	}
	data := models.ListSchedule(uint8(day))
	u.Respond(w, data)
}

func AddScheduleHandler(w http.ResponseWriter, r *http.Request) {
	dayStr := mux.Vars(r)["day"]
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidRequest, "Invalid request"))
		return
	}
	data := models.AddSchedule(uint8(day))
	u.Respond(w, data)
}

func UpdateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	schedule := &models.Shedule{}
	err := json.NewDecoder(r.Body).Decode(schedule)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, err.Error()))
		return
	}
	schedule.Update()
	u.Respond(w, u.Message(u.Ok, "Schedule successfully update"))
}

func DeleteScheduleHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidRequest, "Invalid request"))
		return
	}
	models.DeleteSchedule(uint(id))
	u.Respond(w, u.Message(u.Ok, "Schedule successfully removed"))
}
