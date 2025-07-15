package main

import "net/http"

func (app *application) routes() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/login", app.loginHandler)

    mux.HandleFunc("/teacher/", app.home_teacher)
    mux.HandleFunc("/teacher/help", app.help_teacher)
    mux.HandleFunc("/teacher/your_requests", app.your_requests_teacher)
    mux.HandleFunc("/teacher/meeting", app.meeting_teacher)
    mux.HandleFunc("/teacher/requests", app.requests_teacher)
    mux.HandleFunc("/teacher/choose", app.choose_teacher)
    mux.HandleFunc("/teacher/request", app.request_teacher)
    mux.HandleFunc("/teacher/finish", app.finish_teacher)

    mux.HandleFunc("/student/", app.home)
    mux.HandleFunc("/student/requests", app.requests)
    mux.HandleFunc("/student/help", app.help)
    mux.HandleFunc("/student/teachers", app.teachers)
    mux.HandleFunc("/student/delete", app.delete)
    mux.HandleFunc("/student/meeting", app.meeting)
    mux.HandleFunc("/student/request", app.request)
    mux.HandleFunc("/student/new-request", app.new_request)

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

    return mux
}
