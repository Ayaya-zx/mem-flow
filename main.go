package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Ayaya-zx/mem-flow/themes"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type themeServer struct {
	themeStore themes.ThemeStore
}

func newThemeServer(themeStore themes.ThemeStore) *themeServer {
	return &themeServer{themeStore: themeStore}
}

func main() {
	var port int

	v := viper.New()
	v.SetDefault("Port", 8765)

	v.AutomaticEnv()
	v.SetEnvPrefix("MEMFLOW")
	v.BindEnv("Port", "port")

	pflag.IntVarP(&port, "port", "p", 8765, "Port")
	pflag.Parse()
	if err := v.BindPFlag("Port", pflag.Lookup("port")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server := newThemeServer(themes.NewInmemThemeStore())
	mux := http.NewServeMux()

	mux.HandleFunc("GET /themes", server.getAllThemesHandler)
	mux.HandleFunc("POST /themes", server.createThemeHandler)
	mux.HandleFunc("GET /themes/{id}", server.getThemeHandler)
	mux.HandleFunc("PATCH /themes/{id}", server.repeateThemeHandler)
	mux.HandleFunc("DELETE /themes/{id}", server.deleteThemeHandler)
	mux.HandleFunc("GET /example", exampleHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", v.GetInt("Port")), mux)
	if err != nil {
		log.Fatal(err)
	}
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	var theme themes.Theme

	data, err := json.Marshal(&theme)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(data)
}

func (s *themeServer) getAllThemesHandler(w http.ResponseWriter, r *http.Request) {
	themes, err := s.themeStore.GetAllThemes()
	if err != nil {
		fmt.Println(r.Host, err)
		w.WriteHeader(500)
		return
	}
	data, err := json.Marshal(themes)
	if err != nil {
		fmt.Println(r.Host, err)
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}

func (s *themeServer) createThemeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(r.Host, err)
		w.WriteHeader(400)
		return
	}

	err = s.themeStore.AddTheme(string(data))
	if _, ok := err.(themes.ThemeTitleError); ok {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	if err != nil {
		fmt.Println(r.Host, err)
		w.WriteHeader(500)
		return
	}
}

func (s *themeServer) getThemeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	theme, err := s.themeStore.GetTheme(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(&theme)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(data)
}

func (s *themeServer) repeateThemeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	err = s.themeStore.ThemeRepeated(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}

func (s *themeServer) deleteThemeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	err = s.themeStore.RemoveTheme(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}
