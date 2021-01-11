package cms

import (
  "github.com/golang/glog"

  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-auth"
)

func AuthInit() {
  
  auth.SessionInit(GetConfig().Session.Mode, GetConfig().Session.Expiry_time, GetConfig().Session.Redis.Url, GetConfig().Session.Redis.Max_connections)
  if auth.SessionHasError() {
    glog.Warningf("WRN: SESSION: Init error")
  }
  
  env.LoadFromYMLFiles(GetConfig().ConfigPath + "/auth/", auth.LoadYAML)
  if auth.Count() < 1 {
    glog.Warningf("WRN: AUTH: Not Found auth connectors")
  }
}

func AuthClose() {
  auth.Close()
  auth.SessionClose()
}

func AuthRestart() {
  AuthClose()
  AuthInit()
}
