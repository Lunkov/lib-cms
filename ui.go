package cms

import (
  "net/http"
  "github.com/golang/glog"
  "github.com/gorilla/mux"

  "github.com/Lunkov/lib-auth"
  "github.com/Lunkov/lib-ui"
  "github.com/Lunkov/lib-tr"
)

func getLanguage(params map[string]string, defaultLang string) string {
  lang, ok := params["lang"]
  if !ok {
    if defaultLang == "" {
      return GetConfig().Main.DefaultLang
    }
    return defaultLang
  }
  return lang
}

func getPage(params map[string]string) string {
  page, ok := params["page"]
  if !ok {
    return GetConfig().UI.DefaultPage
  }
  return page
}

func UILogin(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: LOGIN")
  }
  
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  params := mux.Vars(r)
  
  if auth.SessionHasError() || auth.Count() < 1 {
    if GetConfig().Main.AuthRestart {
      AuthRestart()
    }
    data := map[string]interface{}{"LANGS": (*tr.GetList()), "IS_AUTH": false, "AUTH_ERROR": "AUTH ERROR"}
    f := ui.RenderPage(getLanguage(params, GetConfig().Main.DefaultLang), "error_login", GetConfig().UI.CSS, false, &data)
    w.Write([]byte(f))
    return
  }
  
  user, ok := auth.SessionHTTPUserInfo(w, r)
  if ok {
    if glog.V(9) {
      glog.Infof("DBG: GO TO AFTER LOGIN PAGE: %s", GetConfig().UI.AfterLoginPage)
    }
    http.Redirect(w, r, GetConfig().UI.AfterLoginPage + getLanguage(params, user.Language), http.StatusMovedPermanently)
    return
  }
  
  r.ParseForm()
  post := make(map[string]string)
  for key, value := range r.Form {
    post[key] = value[0]
  }
  authCode := r.Form.Get("auth_code")
  login := r.Form.Get("login")
  password := r.Form.Get("password")

  sessionID := auth.SessionHTTPStart(w, r)
  if login != "" && password != "" && authCode != "" {
    user, ok := auth.AuthUser(authCode, &post)
    if ok {
      if glog.V(9) {
        glog.Infof("DBG: GO TO AFTER LOGIN PAGE: %s", GetConfig().UI.AfterLoginPage)
      }
      auth.SessionHTTPUserLogin(w, sessionID, &user)
      http.Redirect(w, r, GetConfig().UI.AfterLoginPage + getLanguage(params, user.Language), http.StatusMovedPermanently)
      return
    }
  }
  if glog.V(9) {
    glog.Infof("DBG: RENDER LOGIN PAGE")
  }
  data := map[string]interface{}{"OAUTH_STATE": SHA1(sessionID), "AUTH_PWD_TYPES": (*auth.GetListPwd()), "AUTH_OAUTH_TYPES": (*auth.GetListOAuth()), "LANGS": (*tr.GetList()) }
  f := ui.RenderPage(getLanguage(params, GetConfig().Main.DefaultLang), "login", GetConfig().UI.CSS, false, &data)
  w.Write([]byte(f))
}

func UILogout(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: LOGOUT")
  }
  sessionID := auth.SessionHTTPStart(w, r)
  auth.SessionHTTPUserLogout(w, sessionID)
  http.Redirect(w, r, GetConfig().UI.DefaultPage, http.StatusMovedPermanently)
}

func UIRedirect(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: HOME REDIRECT: %v", r.URL.String())
  }
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  params := mux.Vars(r)
  user, ok := auth.SessionHTTPUserInfo(w, r)
  if ok {
    http.Redirect(w, r, GetConfig().UI.AfterLoginPage + getLanguage(params, user.Language), http.StatusMovedPermanently)
    return
  }
  http.Redirect(w, r, GetConfig().UI.DefaultPage + getLanguage(params, GetConfig().Main.DefaultLang), http.StatusMovedPermanently)
}

func UIPage(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: PUBLIC PAGE")
  }
  
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  params := mux.Vars(r)
  user, ok := auth.SessionHTTPUserInfo(w, r)
  data := map[string]interface{}{"LANGS": (*tr.GetList()), "IS_AUTH": false}
  if ok {
    data["USER"] = &user // maps.ConvertToMap(user)
    data["IS_AUTH"] = true
  }
  f := ui.RenderPage(getLanguage(params, GetConfig().Main.DefaultLang), getPage(params), GetConfig().UI.CSS, false, &data)
  w.Write([]byte(f))
}

func UIPrivatePage(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: PRIVATE PAGE")
  }
  
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  params := mux.Vars(r)
  user, ok := auth.SessionHTTPUserInfo(w, r)
  if !ok {
    http.Redirect(w, r, GetConfig().UI.LoginPage + getLanguage(params, GetConfig().Main.DefaultLang), http.StatusTemporaryRedirect)
    return
  }
  data := map[string]interface{}{"LANGS": (*tr.GetList()), "IS_AUTH": true, "USER": &user} //maps.ConvertToMap(user)}
  f := ui.RenderPage(getLanguage(params, GetConfig().Main.DefaultLang), getPage(params), GetConfig().UI.CSS, true, &data)
  w.Write([]byte(f))
}
