package models

type MataPelajaran struct {
	ID            int    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	KelasID       int    `gorm:"column:kelasid" json:"kelasId"`
	NamaPelajaran string `gorm:"column:namapelajaran" json:"namapelajaran"`
	JSONQuestion  string `gorm:"column:jsonquestion" json:"jsonquestion"`
	SortBy        string `gorm:"column:sort_by" json:"sortBy"`
}

// TableName mengembalikan nama tabel yang benar untuk GORM.
func (MataPelajaran) TableName() string {
	return "matapelajaran"
}

type Question struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Question     string   `json:"question"`
	Answers      []string `json:"answers"`
	CorrectIndex int      `json:"correctIndex"`
	// Image        string   `json:"image"`
}

// QuestionsResponse represents the full JSON structure
type QuestionsResponse struct {
	Questions []Question `json:"questions"`
}

// MataPelajaranResponse represents the response structure for mata pelajaran
type MataPelajaranResponse struct {
	ID            int               `json:"id"`
	KelasID       int               `json:"kelasId"`
	NamaPelajaran string            `json:"namapelajaran"`
	Questions     QuestionsResponse `json:"questions"`
	SortBy        string            `json:"sortBy"`
}