package chat

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/models"
	"github.com/izzy-Ti/RaGO/internals/rag"
	utils "github.com/izzy-Ti/RaGO/internals/utils"
)

type ASK struct {
	Query  string `json:"query"`
	ChatID uint   `json:"chatId"`
}
type Response struct {
	Ans     string `json:"answer"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}
type ResponseChat struct {
	Ans     interface{} `json:"answer"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
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
	title := req.Query
	if len(title) > 20 {
		title = title[:20]
	}

	if req.ChatID == 0 {
		chat = models.Chat{
			UserID: user.ID,
			Title:  title,
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
		utils.WriteJson(w, http.StatusInternalServerError, Response{
			Message: "RAG processing failed",
			Success: false,
		})
		return
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
func GetChat(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(*models.User)
	var chats []models.Chat
	err := db.DB.Preload("Messages").Where("user_id = ?", user.ID).Find(&chats).Error
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, err)
	}
	res := ResponseChat{
		Ans:     chats,
		Message: "answer successfull",
		Success: true,
	}
	utils.WriteJson(w, http.StatusOK, res)
}
func GetMessagesByChatID(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	vars := mux.Vars(r)
	chatIDstr, ok := vars["chat_ID"]
	if !ok {
		utils.WriteJson(w, http.StatusUnauthorized, ok)
		return
	}
	chatID, err := strconv.Atoi(chatIDstr)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, err)
		return
	}
	var chat models.Chat
	err = db.DB.Where("ID = ? AND User_ID = ?", chatID, user.ID).First(&chat).Error
	if err != nil {
		res := ResponseChat{
			Ans:     nil,
			Message: err.Error(),
			Success: false,
		}
		utils.WriteJson(w, http.StatusNotFound, res)
		return
	}
	var messages []models.Message
	err = db.DB.Where("Chat_ID = ?", chat.ID).Find(&messages).Error
	if err != nil {
		res := ResponseChat{
			Ans:     nil,
			Message: err.Error(),
			Success: false,
		}
		utils.WriteJson(w, http.StatusNotFound, res)
		return
	}
	res := ResponseChat{
		Ans:     messages,
		Message: "successfully fetched his",
		Success: true,
	}
	utils.WriteJson(w, http.StatusOK, res)
}
