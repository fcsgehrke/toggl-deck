package postgres

import (
	"time"

	"github.com/gofrs/uuid"
	pq "github.com/lib/pq"
	"gorm.io/gorm"
)

type Deck struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid ()"`

	Shuffled  bool
	Remaining int
	Cards     pq.StringArray `gorm:"type:text[]"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
