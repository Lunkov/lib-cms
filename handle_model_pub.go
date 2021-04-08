package cms

import (
  "net/http"
  "strconv"
  "strings"
  "errors"

  "github.com/gorilla/mux"
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-ui"
)

func (c *CMS) getItemModelPublic(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Methods", "GET")
  w.Header().Set("Content-Type", "application/json")
  
  _, user, ok := c.checkAuth(w, r)
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

func (c *CMS) postItemModelPublic(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Methods", "POST")
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
  
  var p map[string]interface{}

  err := ui.DecodeJSONBody(w, r, &p)
  if err != nil {
    var mr *malformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.msg, mr.status)
    } else {
      glog.Infof("ERR: %v", err.Error())
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }
  ok = c.DB.DBInsert(model_id, user, &p)
  if ok {
    w.WriteHeader(http.StatusOK)
  } else {
    w.WriteHeader(http.StatusOK)
  }
}

func (c *CMS) getTableModelPublic(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Methods", "GET")
  w.Header().Set("Content-Type", "application/json")
  
  glog.Infof("DBG: +++ '%s' getTableModelPublic", r.URL.Path)
  
  _, user, ok := c.checkAuth(w, r)
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
