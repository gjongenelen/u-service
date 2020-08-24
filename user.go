package api

import (
	"errors"
	"github.com/dongri/phonenumber"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

const USER_LEVEL_ADMIN = 20
const USER_LEVEL_OPERATOR = 12
const USER_LEVEL_SYSTEM = 25
const USER_LEVEL_USER = 5

type User struct {
	Id       uuid.UUID   `json:"id"`
	Email    string      `json:"email"`
	Phone    string      `json:"phone"`
	Verified bool        `json:"verified"`
	Active   bool        `json:"active"`
	Region   string      `json:"region"`
	Name     string      `json:"name"`
	Level    int         `json:"level"`
	Accounts []uuid.UUID `json:"accounts"`
	Password string      `json:"password"`
}

func (u *User) Present() map[string]interface{} {
	return map[string]interface{}{
		"id":        u.Id,
		"email":     u.Email,
		"phone":     u.Phone,
		"activated": u.Verified,
		"name":      u.Name,
		"region":    u.Region,
		"accounts":  u.Accounts,
	}
}

func (u *User) HasAccessToAccount(id uuid.UUID) bool {
	if u.Level == USER_LEVEL_ADMIN {
		return true
	}
	for _, account := range u.Accounts {
		if id == account {
			return true
		}
	}
	return false
}

func Regions() []string {
	return []string{
		"Asia",
		"Africa",
		"North America",
		"South America",
		"Antarctica",
		"Europe",
		"Australia",
	}
}
func IsRegion(reg string) bool {
	for _, region := range Regions() {
		if reg == region {
			return true
		}
	}
	return false
}

func validatePassword(password string) error {
	if password == "" {
		return errors.New("surprise surprise, password may not be empty")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if len(password) > 100 {
		return errors.New("password exceeds 100 characters, you won't memorize it")
	}

	return nil
}

func (u *User) Validate(newUser bool) error {
	if u.Name == "" {
		return errors.New("name may not be empty")
	}
	if u.Email == "" {
		return errors.New("email may not be empty")
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(u.Email) < 3 || len(u.Email) > 254 || !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email")
	}

	if u.Phone != "" {
		country := phonenumber.GetISO3166ByNumber(u.Phone, true)
		if country.CountryCode == "" {
			return errors.New("invalid phone number")
		}
	}
	if !IsRegion(u.Region) {
		return errors.New("invalid region, please pick one of: " + strings.Join(Regions(), ","))
	}
	if u.Password == "" {
		return errors.New("surprise surprise, password may not be empty")
	}
	if newUser {
		err := validatePassword(u.Password)
		if err != nil {
			return err
		}
	}

	return nil
}
