package admin

import (
	"context"
	"net/http"
	"os"

	astradb "github.com/datastax/astra-db-go"
	utils "github.com/izzy-Ti/RaGO/internals/utils"

	DBS "github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/embeddings"
	"github.com/izzy-Ti/RaGO/internals/models"
)

type AdminReq struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Image   []string `json:"image"`
}
type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

var jwtSecret = []byte(os.Getenv("JWT_KEY"))
var ASTRA *astradb.Db

func SaveAdminPost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	userId, err := utils.UserId(token.Value, []byte(jwtSecret))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	var req AdminReq
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	adminPost := models.Posts{
		Title:    req.Title,
		Content:  req.Content,
		Uploadby: userId,
	}
	if err := DBS.DB.Create(&adminPost).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ctx := context.Background()
	col := ASTRA.Collection("GORag3")

	fullText := req.Title + "\n\n" + req.Content

	vec, err := embeddings.EmbedText(fullText)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	astrapost := embeddings.Posts{
		Title:   req.Title,
		Content: req.Content,
		Vector:  vec,
	}

	if _, err := col.InsertOne(ctx, astrapost); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	res := Response{
		Message: "post successfully",
		Success: true,
	}
	utils.WriteJson(w, http.StatusOK, res)

}
