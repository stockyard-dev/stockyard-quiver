package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-quiver/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){status:=r.URL.Query().Get("status");typ:=r.URL.Query().Get("type");list,_:=s.db.List(status,typ);if list==nil{list=[]store.Item{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var i store.Item;json.NewDecoder(r.Body).Decode(&i);if i.Title==""{writeError(w,400,"title required");return};if i.Type==""{i.Type="article"};if i.Status==""{i.Status="unread"};s.db.Create(&i);writeJSON(w,201,i)}
func(s *Server)handleUpdate(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Status string `json:"status"`;Progress int `json:"progress"`;Notes string `json:"notes"`};json.NewDecoder(r.Body).Decode(&req);s.db.Update(id,req.Status,req.Progress,req.Notes);writeJSON(w,200,map[string]string{"status":"updated"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
