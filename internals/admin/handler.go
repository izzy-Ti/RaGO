package admin

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	astradb "github.com/datastax/astra-db-go"
	"github.com/izzy-Ti/RaGO/internals/db"
	utils "github.com/izzy-Ti/RaGO/internals/utils"

	DBS "github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/embeddings"
	"github.com/izzy-Ti/RaGO/internals/models"
	"github.com/ledongthuc/pdf"
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
type PostResponse struct {
	Posts   interface{} `json:"posts"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
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
func AllPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Posts

	err := DBS.DB.Order("created_at DESC").Find(&posts).Error
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, Response{
			Message: "Database error",
			Success: false,
		})
	}

	utils.WriteJson(w, http.StatusOK, PostResponse{
		Posts:   posts,
		Message: "fetch successfull",
		Success: true,
	})

}
func UploadSitePDF(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()
	os.MkdirAll("./temp", os.ModePerm)
	tmpPath := "./temp/" + header.Filename
	out, err := os.Create(tmpPath)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	f, err := os.Open(tmpPath)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer f.Close()
	var pdfText strings.Builder

	reader, err := pdf.NewReader(f, header.Size)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	totalPage := reader.NumPage()

	for i := 1; i <= totalPage; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		content, _ := page.GetPlainText(nil)
		pdfText.WriteString(content)
	}
	vec, err := embeddings.EmbedText(pdfText.String())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	col := ASTRA.Collection("GORag3")
	ctx := context.Background()

	astrapost := embeddings.Posts{
		Title:   header.Filename,
		Content: pdfText.String(),
		Vector:  vec,
	}

	if _, err := col.InsertOne(ctx, astrapost); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	content := pdfText.String()
	if len(content) > 250 {
		content = content[:250]
	}
	if !utf8.ValidString(content) {
		content = string([]rune(content))
	}
	ragKNow := &models.RagKnowladge{
		Title:   header.Filename,
		Content: content,
	}
	if err := db.DB.Create(ragKNow).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	res := Response{
		Message: "post successfully",
		Success: true,
	}
	utils.WriteJson(w, http.StatusOK, res)

}
func AllKnow(w http.ResponseWriter, r *http.Request) {
	var know []models.RagKnowladge

	err := db.DB.Find(&know).Error
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, Response{
			Message: "Database error",
			Success: false,
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, PostResponse{
		Posts:   know,
		Message: "fetch successfull",
		Success: true,
	})
}
