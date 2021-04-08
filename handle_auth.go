package cms

import (
  "net/http"

  "github.com/golang/glog"

  "github.com/Lunkov/lib-auth/base"
  "github.com/Lunkov/lib-auth"
)

type jwt_struct struct {
	Token string
}

func (c *CMS) checkAuth(w http.ResponseWriter, r *http.Request) (string, *base.User, bool) {
  w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding")
  w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
  sessionID := auth.SessionHTTPStart(w, r)
  if glog.V(9) {
    glog.Infof("DBG: REQ: %s: (session=%s)", r.URL, sessionID)
  }
  user, ok := auth.SessionGetUserInfo(sessionID)
  if glog.V(9) {
    glog.Infof("DBG: REQ: %s: SessionUser: %v, (%s) ok=%t\n", r.URL, user, sessionID, ok)
  }
  return sessionID, user, ok
}

func (c *CMS) SessionUser(w http.ResponseWriter, r *http.Request) {
  _, user, ok := c.checkAuth(w, r)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
  if !ok {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
  w.Write(([]byte)(user.ToJSON()))
}

func (c *CMS) Groups(w http.ResponseWriter, r *http.Request) {
  _, user, ok := c.checkAuth(w, r)
	w.Header().Set("Access-Control-Allow-Methods", "GET")
  if !ok {
    glog.Errorf("ERR: Roles: user(%v): %v\n", ok, user)
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
}
