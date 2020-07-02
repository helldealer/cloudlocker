package cloudlocker

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io/ioutil"
	"net/http"
)

func NewRouter(server *LockerServer) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/set", HandleSet(server)).Methods("POST")
	router.HandleFunc("/get", HandleGet(server)).Methods("POST")
	router.HandleFunc("/delete", HandleDelete(server)).Methods("POST")
	return router
}

func HandleSet(server *LockerServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var e Entry
		err = json.Unmarshal(body, &e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = server.db.Put(e.K, e.V, &opt.WriteOptions{Sync: false})
	}
}

func HandleGet(server *LockerServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		v, err := server.db.Get(body, nil)
		if err != nil {
			if err != leveldb.ErrNotFound {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		_, _ = w.Write(v)
	}
}

func HandleDelete(server *LockerServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//log.Println("r.Body", string(body))
		_ = server.db.Delete(body, nil)
	}
}
