package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx"
)

// User модель
type User struct {
	ID         int    `json:"id" example:"1" format:"int64"`
	Name       string `json:"name" example:"Ivan"`
	Surname    string `json:"surname" example:"Ivanov"`
	Patronimic string `json:"patronimic" example:"Ivanovic"`
	Age        int16  `json:"age" example:"22" format:"int16"`
	Gender     string `json:"gender" example:"male"`
	Nation     string `json:"nation" example:"arab"`
}

type CountryResponse struct {
	Countries []CountryInfo `json:"country"`
}

type CountryInfo struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

// example
var (
	ErrNameInvalid    = errors.New("name is empty")
	ErrSurnameInvalid = errors.New("surname is empty")
)

// AddUser example
type AddUser struct {
	Name       string `json:"name" example:"Ivan"`
	Surname    string `json:"surname" example:"Ivanov"`
	Patronimic string `json:"patronimic" example:"Ivanovic"`
}

// Validation example
func (a AddUser) Validation() error {
	switch {
	case len(a.Name) == 0:
		return ErrNameInvalid
	case len(a.Surname) == 0:
		return ErrSurnameInvalid
	default:
		return nil
	}
}

// UpdateUser example
type UpdateUser struct {
	Name       string `json:"name" example:"Ivan"`
	Surname    string `json:"surname" example:"Ivanov"`
	Patronimic string `json:"patronimic" example:"Ivanovic"`
}

// Validation example
func (a UpdateUser) Validation() error {
	switch {
	case len(a.Name) == 0:
		return ErrNameInvalid
	case len(a.Surname) == 0:
		return ErrSurnameInvalid
	default:
		return nil
	}
}

// UsersAll example
func GetUsersByName(conn *pgx.Conn, name string, limit, offset int) ([]User, error) {
	query := "SELECT * FROM users WHERE name ILIKE $1 LIMIT $2 OFFSET $3"
	rows, err := conn.Query(query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronimic, &user.Age, &user.Gender, &user.Nation)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UserOne example
func UserOne(conn *pgx.Conn, id int) (User, error) {
	row := conn.QueryRow("SELECT * FROM users WHERE id=$1", id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronimic, &user.Age, &user.Gender, &user.Nation)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Если записи не найдено, возвращаем ошибку
			return User{}, ErrNoRow
		}
		return User{}, err
	}

	return user, nil
}

// Insert example
func (a User) Insert(conn *pgx.Conn) (int, error) {
	var newUserId int
	err := conn.QueryRow("INSERT INTO users (name, surname, patronimic, age, gender, nation) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", a.Name, a.Surname, a.Patronimic, a.Age, a.Gender, a.Nation).Scan(&newUserId)
	if err != nil {
		return -1, err
	}
	return newUserId, nil
}

// Delete example
func Delete(conn *pgx.Conn, id int) error {
	result, err := conn.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	// Проверяем, была ли удалена хотя бы одна запись
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user id=%d is not found", id)
	}

	return nil
}

func (a User) Update(conn *pgx.Conn) (User, error) {
	// SQL-запрос для обновления пользователя по его ID
	query := `
		UPDATE users 
		SET name = $1, surname = $2, patronimic = $3
		WHERE id = $4
	`

	// Выполняем запрос на обновление
	_, err := conn.Exec(query, a.Name, a.Surname, a.Patronimic, a.ID)
	if err != nil {
		return User{}, err
	}

	// Запрос для получения обновленного пользователя
	var updatedUser User
	getQuery := `
		SELECT *
		FROM users 
		WHERE id = $1
	`
	err = conn.QueryRow(getQuery, a.ID).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Surname, &updatedUser.Patronimic, &updatedUser.Age, &updatedUser.Gender, &updatedUser.Nation)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (a *User) AddData() error {
	age, err := getAge(a.Name)
	if err != nil {
		return err
	}
	gender, err := getGender(a.Name)
	if err != nil {
		return err
	}
	nation, err := getNation(a.Name)
	if err != nil {
		return err
	}

	a.Age = age
	a.Gender = gender
	a.Nation = nation
	return nil
}

func getAge(name string) (int16, error) {
	b := new(User)
	r_age, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return -1, err
	}
	defer r_age.Body.Close()
	body, err := io.ReadAll(r_age.Body)
	if err != nil {
		return -1, err
	}
	err = json.Unmarshal(body, &b)
	if err != nil {
		return -1, err
	}
	return b.Age, nil
}

func getGender(name string) (string, error) {
	b := new(User)
	r_age, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "", err
	}
	defer r_age.Body.Close()
	body, err := io.ReadAll(r_age.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &b)
	if err != nil {
		return "", err
	}
	return b.Gender, nil
}

func getNation(name string) (string, error) {
	c := new(CountryResponse)
	r_age, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", err
	}
	defer r_age.Body.Close()
	body, err := io.ReadAll(r_age.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &c)
	if err != nil {
		return "", err
	}
	return c.Countries[0].CountryId, nil
}
