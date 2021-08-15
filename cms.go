package cms

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/golang/glog"

  "github.com/Lunkov/lib-ui"
  "github.com/Lunkov/lib-auth"
  "github.com/Lunkov/lib-model"
)

type CMS struct {
  Conf       ConfigInfo
  DB        *models.DBConn
  U         *ui.UI
  HasError   bool
  Sessions  *auth.Session
  Auth      *auth.Auth
}

func New() *CMS {
  return &CMS{Sessions: auth.NewSessions(), Auth: auth.New(), DB: models.New()}
}

func (c *CMS) InitDB() {
  c.DB.Init(models.ConnectStr(c.Conf.PostgresWrite), models.ConnectStr(c.Conf.PostgresRead), c.Conf.ConfigPath)
}

func (c *CMS) InitUI() {
  c.U = ui.NewUI(c.Conf.UI.PathTemplates, &c.Conf.UI.CacheForms, &c.Conf.UI.CacheViews, &c.Conf.UI.CachePages, &c.Conf.UI.CacheRenders) 
  c.U.Init(c.Conf.ConfigPath, c.Conf.UI.EnableFileWatcher, c.Conf.UI.EnableMinify)
}


func (c *CMS) CheckHealth(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if c.Sessions.HasError() {
      glog.Errorf("ERR: Auth API: HasError")
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    next.ServeHTTP(w, r)
  })
}

func (c *CMS) Health(w http.ResponseWriter, r *http.Request) {
  status := "OK"
  if c.Sessions.HasError() {
    status = "ERROR"
  }
  fmt.Fprintf(w, "{\"status\": \"%s\", \"auth\": %d, \"mode\": \"%s\", \"online\": %d}", status, c.Auth.Count(), c.Sessions.Mode(), c.Sessions.Count())
}


func (c *CMS) HandleFuncs(router *mux.Router) {
  if c.Conf.API.Health != "" {
    glog.Infof("LOG: Enable Health Check: %s", c.Conf.API.Health)
    router.HandleFunc(c.Conf.API.Health, c.Health)
  }
  
  // API PRIVATE MODELS
  if c.Conf.API.PrivateZone != "" {
    glog.Infof("LOG: Enable Private API: %s", c.Conf.API.PrivateZone)
    router.HandleFunc(c.Conf.API.PrivateZone+"{model_id}",             c.getTableModel).Methods("GET")
    router.HandleFunc(c.Conf.API.PrivateZone+"{model_id}",             c.postItemModel).Methods("POST")
    router.HandleFunc(c.Conf.API.PrivateZone+"{model_id}/{record_id}", c.updateItemModel).Methods("PUT")
    router.HandleFunc(c.Conf.API.PrivateZone+"{model_id}/{record_id}", c.getItemModel).Methods("GET")
    router.HandleFunc(c.Conf.API.PrivateZone+"{model_id}/{record_id}", c.deleteItemModel).Methods("DELETE")
  }
  if c.Conf.API.PublicZone != "" {
    glog.Infof("LOG: Enable Public API: %s", c.Conf.API.PublicZone)
    router.HandleFunc(c.Conf.API.PublicZone+"{model_id}",                  c.getTableModelPublic).Methods("GET")
    router.HandleFunc(c.Conf.API.PublicZone+"{model_id}/{record_id}",      c.getItemModelPublic).Methods("GET")
  }
  
  if c.Conf.Main.EnableUI {
    glog.Infof("LOG: Starting UI")
    // STATIC
    
    fs := http.FileServer(http.Dir("static"))
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    // UI
    router.HandleFunc(c.Conf.UI.LogoutPage  + "{lang}",        c.CheckHealth(c.UILogout))
    router.HandleFunc(c.Conf.UI.LoginPage   + "{lang}",        c.CheckHealth(c.UILogin))
    router.HandleFunc(c.Conf.UI.PrivateZone + "{page}/{lang}", c.CheckHealth(c.UIPrivatePage))
    router.HandleFunc(c.Conf.UI.PublicZone  + "{page}/{lang}", c.CheckHealth(c.UIPage))
    router.HandleFunc("/",                                     c.CheckHealth(c.UIRedirect))
  }

}
