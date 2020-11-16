package cms

import (
  "io/ioutil"
  "path/filepath"

  "gopkg.in/yaml.v2"

  "github.com/golang/glog"

  "github.com/Lunkov/lib-model"
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-auth"
  "github.com/Lunkov/lib-tr"
)

type UIInfo struct {
  CSS                string   `yaml:"css"`
  DefaultPage        string   `yaml:"default_page"`
  LoginPage          string   `yaml:"login_page"`
  AfterLoginPage     string   `yaml:"after_login_page"`
  ErrorLoginPage     string   `yaml:"error_login_page"`
  LogoutPage         string   `yaml:"logout_page"`
  PrivateZone        string   `yaml:"private_zone"`
  PublicZone         string   `yaml:"public_zone"`
  EnableFileWatcher  bool     `yaml:"enable_filewatcher"`
  EnableMinify       bool     `yaml:"enable_minify"`
}

type  APIInfo struct {
  PrivateZone        string   `yaml:"private_zone"`
  PublicZone         string   `yaml:"public_zone"`
  Health             string   `yaml:"health"`
}

type  WebSocketInfo struct {
  Url                string   `yaml:"url"`
  ExternalUrl        string   `yaml:"external_url"`
}

type MainInfo struct {
  CODE            string   `yaml:"code"`
  Title           string   `yaml:"title"`
  DefaultLang     string   `yaml:"lang"`
  HTTPPort        string   `yaml:"port"`
  EnableUI        bool     `yaml:"enable-ui"`
  Url             string   `yaml:"url"`
  Storage         string   `yaml:"storage"`
  StaticFiles     string   `yaml:"static_files"`
}

type ConfigInfo struct {
  ConfigPath      string

  Main            MainInfo                `yaml:"main"`
  Session         auth.SessionInfo        `yaml:"session"`
  UI              UIInfo                  `yaml:"ui"`
  API             APIInfo                 `yaml:"api"`
  WS              WebSocketInfo           `yaml:"websocket"`
  PostgresWrite   models.PostgreSQLInfo   `yaml:"postgres_write"`
  PostgresRead    models.PostgreSQLInfo   `yaml:"postgres_read"`
}

var globConf = ConfigInfo{}

func SetConfig(conf ConfigInfo) {
  globConf = conf
}

func GetConfig() *ConfigInfo {
  return &globConf
}

func LoadConfig(filename string, waittime int) ConfigInfo {
  var err error
  var conf = ConfigInfo{ Main: MainInfo{ Title: "" }, Session: auth.SessionInfo{ Mode: "memory", Expiry_time: 120 }}

  env.WaitFile(filename, waittime)

  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s)  #%v ", filename, err)
    return conf
  }

  err = yaml.Unmarshal(yamlFile, &conf)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s): YAML: %v", filename, err)
  }
  
  if conf.ConfigPath == "" {
    conf.ConfigPath = filepath.Dir(filename)
  }
  if conf.Main.StaticFiles == "" {
    conf.Main.StaticFiles = "static"
  }
  if conf.Main.Storage == "" {
    conf.Main.Storage = "/static/storage/"
  }
  if conf.Main.DefaultLang == "" {
    lang, ok := tr.GetLocale()
    if ok {
      conf.Main.DefaultLang = lang
    } else {
      conf.Main.DefaultLang = "ru_RU"
    }
    if glog.V(3) {
      glog.Infof("LOG: DefaultLang = '%s'", conf.Main.DefaultLang)
    }
  }

  return conf
}

