package main

import (
	"fmt"
	_ "github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/student/" {
		app.notFound(w)
		return
	}

	session, err := app.session.Get(r, "session-name")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := app.getUserByID(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html-student/home.page.tmpl",
		"./ui/html-student/base.layout.tmpl",
		"./ui/html-student/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		User *User
	}{
		User: user,
	}

	err = ts.ExecuteTemplate(w, "home.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) requests(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/student/requests" {
		app.notFound(w)
		return
	}

	session, err := app.session.Get(r, "session-name")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	requests, err := app.getRequestsByUserId(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html-student/requests.page.tmpl",
		"./ui/html-student/base.layout.tmpl",
		"./ui/html-student/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		Requests []Request
	}{
		Requests: requests,
	}

	err = ts.ExecuteTemplate(w, "requests.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) help(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/student/help" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html-student/help.page.tmpl",
		"./ui/html-student/base.layout.tmpl",
		"./ui/html-student/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) teachers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/student/teachers" {
		app.notFound(w)
		return
	}

	users, err := app.getUsers()
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html-student/teachers.page.tmpl",
		"./ui/html-student/base.layout.tmpl",
		"./ui/html-student/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		Users []User
	}{
		Users: users,
	}

	err = ts.ExecuteTemplate(w, "teachers.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}
func (app *application) delete(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	idParam := queryParams.Get("delete")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = app.deleteRequestByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete record: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/student/requests", http.StatusSeeOther)
}

func (app *application) meeting(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/student/meeting" {
		app.notFound(w)
		return
	}

	queryParams := r.URL.Query()
	idParam := queryParams.Get("meeting")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	request, err := app.getRequestByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html-student/meeting.page.tmpl",
		"./ui/html-student/base.layout.tmpl",
		"./ui/html-student/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		Request Request
	}{
		Request: request,
	}

	err = ts.ExecuteTemplate(w, "meeting.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) request(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/student/request" {
		app.notFound(w)
		return
	}

	queryParams := r.URL.Query()
	idParam := queryParams.Get("request")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	request, err := app.getRequestByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html-student/request.page.tmpl",
		"./ui/html-student/base.layout.tmpl",
		"./ui/html-student/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
	}

	data := struct {
		Request Request
	}{
		Request: request,
	}

	err = ts.ExecuteTemplate(w, "request.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) new_request(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/student/new-request" {
		app.notFound(w)
	}

	if r.Method == http.MethodPost {

		title := r.FormValue("title")
		content := r.FormValue("content")

		if title == "" || content == "" {
			http.Error(w, "Заполните все поля", http.StatusBadRequest)
			return
		}

		session, err := app.session.Get(r, "session-name")
		if err != nil {
			app.serverError(w, err)
			return
		}
		userID, ok := session.Values["userID"].(int)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err = app.addRequestToDB(title, content, userID)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, "/student/requests", http.StatusSeeOther)
		return
	}

	files := []string{
		"./ui/html-student/new-request.page.tmpl",
		"./ui/html-student/base.layout.tmpl",
		"./ui/html-student/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		files := []string{
			"./ui/html-login/login.page.tmpl",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			app.serverError(w, err)
		}

	case http.MethodPost:

		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := app.authenticateUser(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		session, _ := app.session.Get(r, "session-name")
		session.Values["userID"] = user.ID
		session.Values["username"] = user.Name
		session.Save(r, w)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, "/student", http.StatusSeeOther)

	default:

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (app *application) home_teacher(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/teacher/" {
		app.notFound(w)
		return
	}

	session, err := app.session.Get(r, "session-name")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := app.getUserByID(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html-teacher/home.page.tmpl",
		"./ui/html-teacher/base.layout.tmpl",
		"./ui/html-teacher/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		User *User
	}{
		User: user,
	}

	err = ts.ExecuteTemplate(w, "home.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) help_teacher(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/teacher/help" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html-teacher/help.page.tmpl",
		"./ui/html-teacher/base.layout.tmpl",
		"./ui/html-teacher/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) your_requests_teacher(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/teacher/your_requests" {
		app.notFound(w)
		return
	}

	session, err := app.session.Get(r, "session-name")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	requests, err := app.getYourRequests(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html-teacher/requests.page.tmpl",
		"./ui/html-teacher/base.layout.tmpl",
		"./ui/html-teacher/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		Requests []Request
	}{
		Requests: requests,
	}

	err = ts.ExecuteTemplate(w, "requests.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) meeting_teacher(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/teacher/meeting" {
		app.notFound(w)
		return
	}

	queryParams := r.URL.Query()
	idParam := queryParams.Get("meeting")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID1232", http.StatusBadRequest)
		return
	}

	request, err := app.getRequestByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if r.Method == http.MethodPost {

		zoomID := r.FormValue("zoomID")
		if zoomID == "" {
			http.Error(w, "Zoom ID cannot be empty", http.StatusBadRequest)
			return
		}

		query := `UPDATE requests SET zoomid = $1 WHERE id = $2`
		_, err = app.db.Exec(query, zoomID, request.ID)
		if err != nil {
			http.Error(w, "Unable to update Zoom ID", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/teacher/your_requests", http.StatusSeeOther)
		return
	}

	files := []string{
		"./ui/html-teacher/meeting.page.tmpl",
		"./ui/html-teacher/base.layout.tmpl",
		"./ui/html-teacher/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		Request Request
	}{
		Request: request,
	}

	err = ts.ExecuteTemplate(w, "meeting.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) requests_teacher(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/teacher/requests" {
		app.notFound(w)
		return
	}

	session, err := app.session.Get(r, "session-name")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	requests, err := app.getRequests(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html-teacher/requests.page.tmpl",
		"./ui/html-teacher/base.layout.tmpl",
		"./ui/html-teacher/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		Requests []Request
	}{
		Requests: requests,
	}

	err = ts.ExecuteTemplate(w, "requests.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) choose_teacher(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	idParam := queryParams.Get("choose")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	session, err := app.session.Get(r, "session-name")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = app.UpdateRequestByID(id, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update record: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/teacher/your_requests", http.StatusSeeOther)
}

func (app *application) request_teacher(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/teacher/request" {
		app.notFound(w)
		return
	}

	queryParams := r.URL.Query()
	idParam := queryParams.Get("request")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	request, err := app.getRequestByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html-teacher/request.page.tmpl",
		"./ui/html-teacher/base.layout.tmpl",
		"./ui/html-teacher/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := struct {
		Request Request
	}{
		Request: request,
	}

	err = ts.ExecuteTemplate(w, "request.page.tmpl", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) finish_teacher(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	idParam := queryParams.Get("finish")

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	session, err := app.session.Get(r, "session-name")
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	err = app.FinishRequestByID(id, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update record: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/teacher/your_requests", http.StatusSeeOther)
}
