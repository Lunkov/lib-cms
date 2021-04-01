package cms

import (
  "net/http"
  "github.com/golang/glog"
  "github.com/gorilla/mux"

  "github.com/Lunkov/lib-auth"
)

func (c *CMS) GetLanguage(params map[string]string, defaultLang string) string {
  lang, ok := params["lang"]
  if !ok {
    if defaultLang == "" {
      return c.Conf.Main.DefaultLang
    }
    return defaultLang
  }
  return lang
}

func (c *CMS) GetPage(params map[string]string) string {
  page, ok := params["page"]
  if !ok {
    return c.Conf.UI.DefaultPage
  }
  return page
}


func (cm *CMS)_UILogin(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: LOGIN")
  }
  
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  params := mux.Vars(r)
  
  if auth.SessionHasError() || auth.Count() < 1 {
    if cm.Conf.Main.AuthRestart {
      cm.RestartAuth()
    }
    data := map[string]interface{}{"LANGS": (*cm.U.GetLangList()), "IS_AUTH": false, "AUTH_ERROR": "AUTH ERROR"}
    f := cm.U.RenderPage("error_login", cm.GetLanguage(params, cm.Conf.Main.DefaultLang), cm.Conf.UI.CSS, false, &data)
    w.Write([]byte(f))
    return
  }
  
  user, ok := auth.SessionHTTPUserInfo(w, r)
  if ok {
    if glog.V(9) {
      glog.Infof("DBG: GO TO AFTER LOGIN PAGE: %s", cm.Conf.UI.AfterLoginPage)
    }
    http.Redirect(w, r, cm.Conf.UI.AfterLoginPage + cm.GetLanguage(params, user.Language), http.StatusMovedPermanently)
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
        glog.Infof("DBG: GO TO AFTER LOGIN PAGE: %s", cm.Conf.UI.AfterLoginPage)
      }
      auth.SessionHTTPUserLogin(w, sessionID, &user)
      http.Redirect(w, r, cm.Conf.UI.AfterLoginPage + cm.GetLanguage(params, user.Language), http.StatusMovedPermanently)
      return
    }
  }
  if glog.V(9) {
    glog.Infof("DBG: RENDER LOGIN PAGE")
  }
  data := map[string]interface{}{"OAUTH_STATE": SHA1(sessionID), "AUTH_PWD_TYPES": (*auth.GetListPwd()), "AUTH_OAUTH_TYPES": (*auth.GetListOAuth()), "LANGS": (*cm.U.GetLangList()) }
  f := cm.U.RenderPage("login", cm.GetLanguage(params, cm.Conf.Main.DefaultLang), cm.Conf.UI.CSS, false, &data)
  w.Write([]byte(f))
}

func (c *CMS)_UILogout(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: LOGOUT")
  }
  sessionID := auth.SessionHTTPStart(w, r)
  auth.SessionHTTPUserLogout(w, sessionID)
  http.Redirect(w, r, c.Conf.UI.DefaultPage, http.StatusMovedPermanently)
}

func (c *CMS)_UIRedirect(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: HOME REDIRECT: %v", r.URL.String())
  }
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  params := mux.Vars(r)
  user, ok := auth.SessionHTTPUserInfo(w, r)
  if ok {
    http.Redirect(w, r, c.Conf.UI.AfterLoginPage + c.GetLanguage(params, user.Language), http.StatusMovedPermanently)
    return
  }
  http.Redirect(w, r, c.Conf.UI.DefaultPage + c.GetLanguage(params, c.Conf.Main.DefaultLang), http.StatusMovedPermanently)
}

func (c *CMS)_UIPage(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: PUBLIC PAGE")
  }
  
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  params := mux.Vars(r)
  user, ok := auth.SessionHTTPUserInfo(w, r)
  data := map[string]interface{}{"LANGS": (*c.U.GetLangList()), "IS_AUTH": false}
  if ok {
    data["USER"] = &user // maps.ConvertToMap(user)
    data["IS_AUTH"] = true
  }
  f := c.U.RenderPage(c.GetPage(params), c.GetLanguage(params, c.Conf.Main.DefaultLang), c.Conf.UI.CSS, false, &data)
  w.Write([]byte(f))
}

func (c *CMS)_UIPrivatePage(w http.ResponseWriter, r *http.Request)  {
  if glog.V(9) {
    glog.Infof("DBG: PRIVATE PAGE")
  }
  
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  params := mux.Vars(r)
  user, ok := auth.SessionHTTPUserInfo(w, r)
  if !ok {
    http.Redirect(w, r, c.Conf.UI.LoginPage + c.GetLanguage(params, c.Conf.Main.DefaultLang), http.StatusTemporaryRedirect)
    return
  }
  data := map[string]interface{}{"LANGS": (*c.U.GetLangList()), "IS_AUTH": true, "USER": &user} //maps.ConvertToMap(user)}
  f := c.U.RenderPage(c.GetPage(params), c.GetLanguage(params, c.Conf.Main.DefaultLang), c.Conf.UI.CSS, true, &data)
  w.Write([]byte(f))
}
