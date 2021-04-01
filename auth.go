package cms

import (
  "github.com/golang/glog"

  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-auth"
)

func (c *CMS) InitAuth() {
  
  auth.SessionInit(c.Conf.Session.Mode, c.Conf.Session.Expiry_time, c.Conf.Session.Redis.Url, c.Conf.Session.Redis.Max_connections)
  if auth.SessionHasError() {
    glog.Warningf("WRN: SESSION: Init error")
  }
  
  env.LoadFromFiles(c.Conf.ConfigPath + "/auth/", "", auth.LoadYAML)
  if auth.Count() < 1 {
    glog.Warningf("WRN: AUTH: Not Found auth connectors")
  }
}

func (c *CMS) CloseAuth() {
  auth.Close()
  auth.SessionClose()
}

func (c *CMS) RestartAuth() {
  c.CloseAuth()
  c.InitAuth()
}
