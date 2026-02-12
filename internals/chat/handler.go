package chat

import (
	"net/http"

	utils "github.com/izzy-Ti/RaGO/internals/Utils"
	"github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/models"
	"github.com/izzy-Ti/RaGO/internals/rag"
)

type ASK struct {
	Query  string `json:"query"`
	ChatID uint   `json:"chatId`
}
type Response struct {
	Ans     string `json:"answer"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func Ask(w http.ResponseWriter, r *http.Request) {
	var req ASK
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, Response{
			Message: "Invalid request",
			Success: false,
		})
		return
	}
	user, _ := r.Context().Value("user").(*models.User)
	var chat models.Chat

	if req.ChatID == 0 {
		chat = models.Chat{
			UserID: user.ID,
			Title:  req.Query[:20],
		}
		db.DB.Create(&chat)
	} else {
		db.DB.First(&chat, req.ChatID)
	}

	userMsg := models.Message{
		ChatID:  chat.ID,
		Role:    "USER",
		Content: req.Query,
	}
	db.DB.Create(&userMsg)

	ans, err := rag.RAG(req.Query)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, err)
	}
	assistantMsg := models.Message{
		ChatID:  chat.ID,
		Role:    "ASSISTANT",
		Content: ans,
	}
	db.DB.Create(&assistantMsg)

	res := Response{
		Ans:     ans,
		Message: "answer successfull",
		Success: true,
	}

	utils.WriteJson(w, http.StatusOK, res)
}
