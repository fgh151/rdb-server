package bl

import (
	"db-server/modules/project"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Runner string

const (
	BLRunnerPHP  Runner = "php"
	BLRunnerDart Runner = "dart"
)

type BusinessLogic struct {
	// The logic UUID
	// example: 6204011c-30e6-408b-8aaa-dd8219860b4b
	Id          uuid.UUID `gorm:"primarykey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Folder      string    `json:"folder"`

	Runner Runner `json:"runner"`

	// Linked project  UUID
	// example: 6204011c-30e6-408b-8aaa-dd8214860b4b
	ProjectId uuid.UUID `json:"project_id"`

	// Linked project
	Project   project.Project
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (l BusinessLogic) Run() {
	switch l.Runner {
	case BLRunnerPHP:
		//cmd := "docker run -it --rm --name test -v \"$PWD\":/usr/src/app -w /usr/src/app php:7.4-cli php test.php "
	}
}
