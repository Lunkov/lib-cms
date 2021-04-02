package cms

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/golang/glog"

  "github.com/Lunkov/lib-ui"
  "github.com/Lunkov/lib-auth"
)

type CMS struct {
  Conf       ConfigInfo
  U         *ui.UI
  HasError   bool
}

func New() *CMS {
  return &CMS{}
}

func (c *CMS) InitUI() {
  c.U = ui.NewUI(c.Conf.UI.PathTemplates, &c.Conf.UI.CacheForms, &c.Conf.UI.CacheViews, &c.Conf.UI.CachePages, &c.Conf.UI.CacheRenders) 
  c.U.Init(c.Conf.ConfigPath, false, false)
}


func (c *CMS) CheckHealth(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if auth.SessionHasError() {
      glog.Errorf("ERR: Auth API: HasError")
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    next.ServeHTTP(w, r)
  })
}

func (c *CMS) Health(w http.ResponseWriter, r *http.Request) {
  status := "OK"
  if auth.SessionHasError() {
    status = "ERROR"
  }
  fmt.Fprintf(w, "{\"status\": \"%s\", \"auth\": %d, \"mode\": \"%s\", \"online\": %d}", status, auth.Count(), auth.SessionMode(), auth.SessionCount())
}


func (c *CMS) HandleFuncs(router *mux.Router) {
  if c.Conf.API.Health != "" {
    glog.Infof("LOG: Enable Health Check: %s", c.Conf.API.Health)
    router.HandleFunc(c.Conf.API.Health, c.Health)
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
