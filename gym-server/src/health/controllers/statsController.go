package controllers

import (
	"encoding/json"
	"health/src/health/dto"
	"health/src/health/models"
	u "health/src/health/utils"
	"net/http"
	"strconv"

	"github.com/morkid/paginate"
)

func AddStatHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VID  string
		Name string
		Type string
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, "Invalid request body"))
		return
	}
	types := map[string]bool{
		"half":    true,
		"decline": true,
		"full":    true,
	}
	if !types[req.Type] {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, "Invalid stat type [half, decline, full]"))
		return
	}
	id, err := strconv.Atoi(req.VID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(u.InvalidBody, "Invalid request body"))
		return
	}

	userId := models.GetOrCreateUsernameId(req.Name) // Как он может быть 0?
	if userId == 0 {
		models.AddStat(uint(id), req.Name, req.Type) // В Stat никогда не добавится
	} else {
		models.AddPublicStat(uint(id), userId, req.Type)
	}

	u.Respond(w, u.Message(u.Ok, "Stats successfully added"))
}

func ListStatsHandler(w http.ResponseWriter, r *http.Request) {
	model := models.Db.Model(&models.PublicStat{})
	pg := paginate.New()
	paginated := pg.With(model).Request(r).Response(&[]dto.StatDto{})
	j, err := json.Marshal(paginated)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-type", "application/json")
	
	_, err = w.Write(j)
	if err != nil {
		panic(err)
	}
}

func ListStatsDetailsHandler(w http.ResponseWriter, r *http.Request) {
	model := models.Db. //у таблицы users нет методов создания записей
		Select("users.user_id, users.organization, users.\"position\", users.department, videos.\"name\", count(*), public_stats.\"type\"").
		Joins("inner join users on users.user_id = public_stats.user_id").
		Joins("inner join videos on videos.id = public_stats.video_id").
		Group("users.user_id, users.organization, users.\"position\", users.department, videos.\"name\", public_stats.\"type\"").
		Model(&models.PublicStat{})
	pg := paginate.New()
	paginated := pg.With(model).Request(r).Response(&[]dto.StatDetailDto{})
	j, err := json.Marshal(paginated)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		panic(err)
	}
}
