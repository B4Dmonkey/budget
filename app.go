package main

// func (a *App) InitializeServer() {
// 	if err := CreateDatabase(); err != nil {
// 		log.Fatal(err.Error())
// 		return
// 	}
// fSys, err := fs.Sub(publicDir, "public")

// if err != nil {
// 	log.Fatal("Failed to load public dir", err)
// }

// a.mux.Handle(GET+" /assets/", http.FileServer(http.FS(fSys)))

// a.mux.HandleFunc(GET+" /styles.css", func(w http.ResponseWriter, r *http.Request) {
// 	styles, err := publicDir.ReadFile("public/styles.css")
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "text/css")
// 	w.Write(styles)
// })

// 	a.SetRouteHandlers(routes)
// }
