package main

import (
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) getRequestsByUserId(UserID int) ([]Request, error) {
	rows, err := app.db.Query("SELECT * FROM requests WHERE userid = $1", UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var req Request
		if err := rows.Scan(&req.ID, &req.Title, &req.Content, &req.Created, &req.Status, &req.UserID, &req.TeacherID, &req.ZoomID); err != nil {
			return nil, err
		}
		req.FormattedDate = req.Created.Format("15:04 02-01-2006")
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (app *application) getYourRequests(UserID int) ([]Request, error) {
	rows, err := app.db.Query("SELECT * FROM requests WHERE status='В процессе'  AND teacherid=$1", UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var req Request
		if err := rows.Scan(&req.ID, &req.Title, &req.Content, &req.Created, &req.Status, &req.UserID, &req.TeacherID, &req.ZoomID); err != nil {
			return nil, err
		}
		req.FormattedDate = req.Created.Format("15:04 02-01-2006")
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (app *application) getRequests(UserID int) ([]Request, error) {
	rows, err := app.db.Query("SELECT * FROM requests WHERE status='Открыт' AND NOT userid=$1", UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var req Request
		if err := rows.Scan(&req.ID, &req.Title, &req.Content, &req.Created, &req.Status, &req.ZoomID, &req.UserID, &req.TeacherID); err != nil {
			return nil, err
		}
		req.FormattedDate = req.Created.Format("15:04 02-01-2006")
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (app *application) getRequestByID(id int) (Request, error) {
	var req Request
	query := "SELECT id, title, content, created, status, zoomid, userid, teacherid FROM requests WHERE id = $1"
	row := app.db.QueryRow(query, id)

	row.Scan(&req.ID, &req.Title, &req.Content, &req.Created, &req.Status, &req.ZoomID, &req.UserID, &req.TeacherID)

	req.FormattedDate = req.Created.Format("15:04 02-01-2006")
	return req, nil
}

type Request struct {
	ID            int
	Title         string
	Content       string
	Created       time.Time
	FormattedDate string
	Status        string
	ZoomID        string
	UserID        int
	TeacherID     int
}

func (app *application) getUsers() ([]User, error) {
	rows, err := app.db.Query("SELECT * FROM users ORDER BY rating DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var usr User
		if err := rows.Scan(&usr.ID, &usr.Surname, &usr.Created, &usr.Description, &usr.Scores, &usr.Name, &usr.Course, &usr.Group, &usr.Rating, &usr.StudentID, &usr.Login, &usr.Password); err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

type User struct {
	ID          int
	Name        string
	Surname     string
	Created     time.Time
	Description string
	Scores      int
	Course      int
	Group       string
	Rating      float32
	StudentID   int
	Login       string
	Password    string
}

func (app *application) deleteRequestByID(id int) error {
	query := "DELETE FROM requests WHERE id = $1"
	result, err := app.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete request with ID %d: %v", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving result: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no request found with ID %d", id)
	}
	return nil
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addRequestToDB(title, content string, userID int) error {
	query := `INSERT INTO requests (title, content, created, status, zoomid, userid, teacherid) VALUES ($1, $2, NOW(), 'Открыт', 0, $3, 0)`

	_, err := app.db.Exec(query, title, content, userID)
	app.db.Exec("UPDATE users SET scores=scores-1 WHERE id=$1", userID)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) authenticateUser(username, password string) (*User, error) {
	users, err := app.getUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Login == username {
			if user.Password == password {
				return &user, nil
			} else {
				return nil, fmt.Errorf("invalid password")
			}
		}
	}

	return nil, fmt.Errorf("user not found")
}

func (app *application) getUserByID(userID int) (*User, error) {
	var user User

	query := `SELECT * FROM users WHERE id = $1`
	row := app.db.QueryRow(query, userID)

	row.Scan(
		&user.ID, &user.Surname, &user.Created, &user.Description,
		&user.Scores, &user.Name, &user.Course, &user.Group, &user.Rating, &user.StudentID,
		&user.Login, &user.Password)
	return &user, nil
}

func (app *application) UpdateRequestByID(id int, UserID int) error {
	query := "UPDATE requests SET status='В процессе' WHERE id = $1"
	result, err := app.db.Exec(query, id)
	query1 := "UPDATE requests SET teacherid=$2 WHERE id = $1"
	app.db.Exec(query1, id, UserID)
	if err != nil {
		return fmt.Errorf("failed to update request with ID %d: %v", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving result: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no request found with ID %d", id)
	}
	return nil
}

func (app *application) FinishRequestByID(id int, UserID int) error {
	query := "UPDATE requests SET status='Закрыт' WHERE id = $1"
	app.db.Exec("UPDATE users SET scores=scores+1 WHERE id = $1", UserID)
	result, err := app.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update request with ID %d: %v", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving result: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no request found with ID %d", id)
	}
	return nil
}
