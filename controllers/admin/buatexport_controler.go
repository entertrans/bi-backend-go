package admincontrollers

import (
	"encoding/json"

	"github.com/entertrans/bi-backend-go/config"
	"github.com/entertrans/bi-backend-go/models"
)

// Response structure untuk questions
type Question struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Question    string   `json:"question"`
	Answers     []string `json:"answers"`
	CorrectIndex int     `json:"correctIndex"`
	// Image       string   `json:"image"`
}

type QuestionsResponse struct {
	Questions []Question `json:"questions"`
}

type MataPelajaranResponse struct {
	ID            int               `json:"id"`
	KelasID       int               `json:"kelasId"`
	NamaPelajaran string            `json:"namapelajaran"`
	Questions     QuestionsResponse `json:"questions"`
	SortBy        string            `json:"sortBy"`
}

// GetAllJSONQuestions - Get all mata pelajaran dengan JSON questions
func GetAllJSONQuestions() ([]MataPelajaranResponse, error) {
	var mataPelajaranList []models.MataPelajaran

	// Ambil semua data dari database
	err := config.DB.Select("id, kelasid, namapelajaran, jsonquestion, sort_by").Find(&mataPelajaranList).Error
	if err != nil {
		return nil, err
	}

	// Parse response
	return parseMataPelajaranResponses(mataPelajaranList)
}

// GetJSONQuestionsByKelasID - Get mata pelajaran by kelas ID
func GetJSONQuestionsByKelasID(kelasID int) ([]MataPelajaranResponse, error) {
	var mataPelajaranList []models.MataPelajaran

	// Ambil data berdasarkan kelas ID
	err := config.DB.Select("id, kelasid, namapelajaran, jsonquestion, sort_by").
		Where("kelasid = ?", kelasID).
		Find(&mataPelajaranList).Error
	if err != nil {
		return nil, err
	}

	// Parse response
	return parseMataPelajaranResponses(mataPelajaranList)
}

// GetJSONQuestionsByID - Get mata pelajaran by ID
func GetJSONQuestionsByID(id int) (*MataPelajaranResponse, error) {
	var mataPelajaran models.MataPelajaran

	// Ambil data berdasarkan ID
	err := config.DB.Select("id, kelasid, namapelajaran, jsonquestion, sort_by").
		Where("id = ?", id).
		First(&mataPelajaran).Error
	if err != nil {
		return nil, err
	}

	// Parse response
	response, err := parseSingleMataPelajaranResponse(mataPelajaran)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Helper function untuk parse multiple responses
func parseMataPelajaranResponses(mataPelajaranList []models.MataPelajaran) ([]MataPelajaranResponse, error) {
	var responses []MataPelajaranResponse

	for _, mp := range mataPelajaranList {
		response, err := parseSingleMataPelajaranResponse(mp)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// Helper function untuk parse single response
func parseSingleMataPelajaranResponse(mp models.MataPelajaran) (MataPelajaranResponse, error) {
	var questionsResp QuestionsResponse
	
	// Parse the JSONQuestion string into QuestionsResponse
	if err := json.Unmarshal([]byte(mp.JSONQuestion), &questionsResp); err != nil {
		return MataPelajaranResponse{}, err
	}

	return MataPelajaranResponse{
		ID:            mp.ID,
		KelasID:       mp.KelasID,
		NamaPelajaran: mp.NamaPelajaran,
		Questions:     questionsResp,
		SortBy:        mp.SortBy,
	}, nil
}