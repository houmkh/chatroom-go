package service

import (
	"chatroom/cmn"
	"chatroom/dao"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"sort"
)

//go:generate go run service-enroll-generate.go -a=annotation:(?P<name>.*)-service
//var dbConn = dao.ConnDB()
var pool *pgxpool.Pool

//func reqProc(reqPath string, w http.ResponseWriter, r *http.Request) {
//	cmn.Services[reqPath].Fn(w, r, dbConn)
//}

func WebServe() {
	Enroll()
	pool = dao.ConnDB()

	//router := mux.NewRouter()

	var rootExists bool
	var pathList []string
	//pathList = append(pathList, "/")
	for k := range cmn.Services {
		if k == "/" {
			fmt.Println("root")
			rootExists = true
			continue
		}
		pathList = append(pathList, k)
	}
	sort.Strings(pathList)
	if rootExists {
		pathList = append(pathList, "/")
	}
	for _, k := range pathList {
		reqPath := k

		http.HandleFunc(reqPath, func(writer http.ResponseWriter, request *http.Request) {
			var dbConn = dao.GetConn()
			defer dbConn.Release()
			cmn.Services[reqPath].Fn(writer, request, dbConn)
		})

	}

}
