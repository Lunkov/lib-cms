package cms

import (
  "github.com/golang/glog"

  "github.com/Lunkov/lib-env"
)

func (c *CMS) InitAuth() {
  
  c.Sessions.Init(c.Conf.Session.Mode, c.Conf.Session.Expiry_time, c.Conf.Session.Redis.Url, c.Conf.Session.Redis.Max_connections)
  if c.Sessions.HasError() {
    glog.Warningf("WRN: SESSION: Init error")
  }
  
  env.LoadFromFiles(c.Conf.ConfigPath + "/auth/", "", c.Auth.Load)
  if c.Auth.Count() < 1 {
    glog.Warningf("WRN: AUTH: Not Found auth connectors")
  }
}

func (c *CMS) CloseAuth() {
  c.Auth.Close()
  c.Sessions.Close()
}

func (c *CMS) RestartAuth() {
  c.CloseAuth()
  c.InitAuth()
}
