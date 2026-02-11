package chat

import (
	"net/http"

	utils "github.com/izzy-Ti/RaGO/internals/Utils"
	"github.com/izzy-Ti/RaGO/internals/rag"
)

type ASK struct {
	Query string `json:"query"`
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

	ans, err := rag.RAG(req.Query)
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, err)
	}
	res := Response{
		Ans:     ans,
		Message: "answer successfull",
		Success: true,
	}

	utils.WriteJson(w, http.StatusOK, res)
}
