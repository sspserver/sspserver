//
// @project GeniusRabbit AdNet
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015 – 2016, 2018
//

package models

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"

	"github.com/guregu/null"
)

type UserStatus int

// User status
const (
	UserStatusInactive       UserStatus = 0 // 0
	UserStatusActive         UserStatus = 1 // 1
	UserStatusClauseForFraud UserStatus = 2 // 2
)

const (
	userPasswordSalt = "ad6eithae6AehaefooghaishieCegahv"
)

// User model
type User struct {
	ID uint64 `json:"id"`

	Status      UserStatus `json:"status"`
	IsSuperuser bool       `json:"is_superuser"`
	Username    string     `json:"username"`
	Password    string     `json:"password"`

	Companies []*Company `gorm:"many2many:company_m2m_member;" json:"companies"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

// TableName in database
func (u *User) TableName() string {
	return "usr"
}

// IsAuth user
func (u *User) IsAuth() bool {
	return u != nil && u.ID > 0
}

// PasswordHash code
func PasswordHash(password string) string {
	h := sha1.New()
	_, _ = io.WriteString(h, userPasswordSalt)
	_, _ = io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}
