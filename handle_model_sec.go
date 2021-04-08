package cms

import (
  "fmt"
  "os"
  "time"
  // "bytes"
  "net/http"
  "strconv"
  "strings"
  "encoding/json"

  "github.com/gorilla/mux"
  "github.com/golang/glog"
)

type ResultInfo struct {
  Status       string                 `json:"status"`
}

type malformedRequest struct {
  status int
  msg    string
}

func (mr *malformedRequest) Error() string {
  return mr.msg
}

func (c *CMS) getItemModel(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Methods", "GET")
  w.Header().Set("Content-Type", "application/json")
  
  glog.Infof("DBG: +++ '%s' getItemModel", r.URL.Path)
  
  _, user, ok := c.checkAuth(w, r)
  if !ok {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")

  params := mux.Vars(r)
  model_id, ok := params["model_id"]
  if !ok {
    glog.Errorf("ERR: URL '%s': DoN`t Set Model `%v`\n", r.URL.Path, params)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  record_id, okr := params["record_id"]
  if !okr {
    glog.Errorf("ERR: URL '%s': DoN`t Set RecordID `%v`\n", r.URL.Path, params)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  jsonRes, okg := c.DB.DBGetItemByID(model_id, user, record_id)
  
  if okg {  
    w.WriteHeader(http.StatusOK)
    w.Write(jsonRes)
  } else {
    w.WriteHeader(http.StatusBadRequest)
  }
    
}

func (c *CMS) deleteItemModel(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Methods", "DELETE")
  w.Header().Set("Content-Type", "application/json")
  
  _, user, ok := c.checkAuth(w, r)
  if !ok {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")

  params := mux.Vars(r)
  model_id, ok := params["model_id"]
  if !ok {
    glog.Errorf("ERR: URL '%s': DoN`t Set Model `%v`\n", r.URL.Path, params)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  record_id, okr := params["record_id"]
  if !okr {
    glog.Errorf("ERR: URL '%s': DoN`t Set Model `%v`\n", r.URL.Path, params)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  okd := c.DB.DBDeleteItemByID(model_id, user, record_id)
  
  if okd {  
    w.WriteHeader(http.StatusOK)
  } else {
    w.WriteHeader(http.StatusBadRequest)
  }
}

func (c *CMS) postItemModel(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Methods", "POST")
  w.Header().Set("Content-Type", "application/json")
  glog.Infof("DBG: +++ '%s' 111", r.URL.Path)
  _, user, ok := c.checkAuth(w, r)
  if !ok {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  glog.Infof("DBG: +++ '%s' 222", r.URL.Path)
  w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")

  params := mux.Vars(r)
  model_id, ok := params["model_id"]
  if !ok {
    glog.Errorf("ERR: URL '%s': DoN`t Set Model `%v`\n", r.URL.Path, params)
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  glog.Infof("DBG: +++ '%s' 333 %s === %v", r.URL.Path, model_id, user)
  /*
  var p map[string]interface{}

  err := c.DB.decodeJSONBody(w, r, &p)
  if err != nil {
    var mr *malformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.msg, mr.status)
    } else {
      glog.Infof("ERR: %v", err.Error())
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }*/
  
  t := time.Now()
  savePath := fmt.Sprintf("/%d/%d/%d/", t.Year(), t.Month(), t.Day())
  os.MkdirAll(savePath, os.ModePerm)
  parameters, ok := c.U.Forms.GetParameters(r, c.Conf.Main.Storage, savePath)
  glog.Infof("DBG: UploadFile URL '%s': `%v`", r.URL.Path, parameters)
  
  
  ok = c.DB.DBInsert(model_id, user, &parameters)
  if ok {
    w.WriteHeader(http.StatusOK)
    jsonRes, _ := json.Marshal(ResultInfo{Status:"OK"})
    w.Write(jsonRes)
  } else {
    w.WriteHeader(http.StatusOK)
    jsonRes, _ := json.Marshal(ResultInfo{Status:"ERROR"})
    w.Write(jsonRes)
  }
}

func (c *CMS) updateItemModel(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Methods", "PUT")
  w.Header().Set("Content-Type", "application/json")
  _, user, ok := c.checkAuth(w, r)
  if !ok {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")

  params := mux.Vars(r)
  model_id, ok := params["model_id"]
  if !ok {
    glog.Errorf("ERR: URL '%s': DoN`t Set Model `%v`\n", r.URL.Path, params)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  t := time.Now()
  savePath := fmt.Sprintf("/%d/%d/%d/", t.Year(), t.Month(), t.Day())
  os.MkdirAll(savePath, os.ModePerm)
  parameters, ok := c.U.Forms.GetParameters(r, c.Conf.Main.Storage, savePath)
  glog.Infof("DBG: FormGetParameters URL '%v': `%v`", ok, parameters)
  
  ok = c.DB.DBUpdate(model_id, user, &parameters)
  if ok {
    w.WriteHeader(http.StatusOK)
    jsonRes, _ := json.Marshal(ResultInfo{Status:"OK"})
    w.Write(jsonRes)
  } else {
    w.WriteHeader(http.StatusOK)
    jsonRes, _ := json.Marshal(ResultInfo{Status:"ERROR"})
    w.Write(jsonRes)
  }
}

func (c *CMS) getTableModel(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Methods", "GET")
  w.Header().Set("Content-Type", "application/json")
  
  _, user, ok := c.checkAuth(w, r)
  if !ok {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")

  params := mux.Vars(r)
  model_id, ok := params["model_id"]
  if !ok {
    glog.Errorf("ERR: URL '%s': DoN`t Set Model `%v`\n", r.URL.Path, params)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  r.ParseForm()
  offset, _ := strconv.Atoi(r.Form.Get("offset"))
  limit, _ := strconv.Atoi(r.Form.Get("limit"))
  fields := r.Form.Get("select")
  ar_fields := strings.Split(fields, ",")
  order := strings.ReplaceAll(r.Form.Get("order"), ".", " ")
  ar_order := strings.Split(order, ",")
  
  jsonRes, count, ok := c.DB.DBTableGet(model_id, user, ar_fields, ar_order, offset, limit)
  
  w.Header().Set("Range-Unit", "items")
  w.Header().Set("Content-Range", strconv.Itoa(count))

  if ok {  
    w.WriteHeader(http.StatusOK)
    w.Write(jsonRes)
  } else {
    glog.Errorf("ERR: REQUEST: model_id=%v, ar_fields=%v, ar_order=%v, offset=%v, limit=%v => jsonRes=%v, count=%v, ok=%v", model_id, ar_fields, ar_order, offset, limit, jsonRes, count, ok)
    w.WriteHeader(http.StatusBadRequest)
  }

}

