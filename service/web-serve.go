package service

import (
	"chatroom/cmn"
	"chatroom/dao"
	"fmt"
	"net/http"
	"sort"
)

//go:generate go run service-enroll-generate.go -a=annotation:(?P<name>.*)-service
var dbConn = dao.ConnDB()

//func reqProc(reqPath string, w http.ResponseWriter, r *http.Request) {
//	cmn.Services[reqPath].Fn(w, r, dbConn)
//}

func WebServe() {
	Enroll()

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
	//fmt.Println(pathList)
	for _, k := range pathList {
		reqPath := k

		//fmt.Println(reqPath)
		http.HandleFunc(reqPath, func(writer http.ResponseWriter, request *http.Request) {
			//reqProc(k, writer, request)
			//fmt.Println(request, dbConn)
			cmn.Services[reqPath].Fn(writer, request, dbConn)
		})

		//router.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
		//	reqProc(k, w, r)
		//})
	}
	//fmt.Println("end")
	//serv := &http.Server{
	//	Addr:    strconv.Itoa(8082),
	//	Handler: GzipHandler(router),
	//}

	//err := http.ListenAndServe(":8082", nil)
	//if err != nil {
	//	fmt.Println("listen port failed")
	//	return
	//}

	//host := "qnear.cn"
	//if viper.IsSet("webServe.serverName") {
	//	host = viper.GetString("webServe.serverName")
	//}

	//appLaunchPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//if err != nil {
	//	//z.Fatal(err.Error())
	//	return
	//}

	//certPath := appLaunchPath + "/certs"
	//var hostWhiteList string103.17
	//if viper.IsSet("webServe.hostWhiteList") {
	//	hostWhiteList = viper.GetString("webServe.hostWhiteList")
	//	names := strings.Split(hostWhiteList, ",")
	//	host := "qnear.cn"
	//	if viper.IsSet("webServe.serverName") {
	//		host = viper.GetString("webServe.serverName")
	//	}
	//	var exists bool
	//	for _, v := range names {
	//		if v == host {
	//			exists = true
	//			break
	//		}
	//	}
	//	if !exists {
	//		log.Fatal(fmt.Sprintf("webServe.serverName:%s must exists in webServe.hostWhiteList: %s",
	//			host, hostWhiteList))
	//	}
	//}

	//if hostWhiteList == "" {
	//	hostWhiteList = host
	//}
	//certManager := autocert.Manager{
	//	Prompt: autocert.AcceptTOS,
	//
	//	HostPolicy: autocert.HostWhitelist(
	//		strings.Split(hostWhiteList, ",")...), //Your domain here
	//
	//	Cache: autocert.DirCache(certPath), //Folder for storing certificates
	//}

	//getWxAccessToken(2)

	//httpListenPort := 8080
	//if viper.IsSet("webServe.httpListenPort") {
	//	httpListenPort = viper.GetInt("webServe.httpListenPort")
	//}

	//httpsListenPort := 8443
	//if viper.IsSet("webServe.httpsListenPort") {
	//	httpsListenPort = viper.GetInt("webServe.httpsListenPort")
	//}

	//var autoCert bool
	//if viper.IsSet("webServe.autoCert") {
	//	autoCert = viper.GetBool("webServe.autoCert")
	//}

	//var ep string
	//if autoCert {
	//	ep = fmt.Sprintf(":%v", httpsListenPort)
	//} else {
	//ep = fmt.Sprintf(":%v", httpListenPort)
	//}

	//s1 := "***********************************************************"
	//s2 := "   ************ app started ****************************"
	//s3 := fmt.Sprintf("                  db: %s@%s:%d/%s", viper.GetString("dbms.postgresql.user"),
	//	viper.GetString("dbms.postgresql.addr"),
	//	viper.GetInt32("dbms.postgresql.port"),
	//	viper.GetString("dbms.postgresql.db"))
	//s8 := fmt.Sprintf("             version: %s", cmn.GetBuildVer())
	//s4 := fmt.Sprintf("               redis: %s:%d", viper.GetString("dbms.redis.addr"),
	//	viper.GetInt32("dbms.redis.port"))
	//
	//s5 := "      web serve on *" + ep
	//
	//s6 := "   *****************************************************"
	//s7 := "***********************************************************"

	//z.Info(fmt.Sprintf("\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n", s1, s2, s3, s4, s5, s8, s6, s7))

	//serv := &http.Server{
	//	Addr:    ep,
	//	Handler: GzipHandler(router),
	//}

	//if autoCert {
	//	serv.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate}
	//	go func() { _ = http.ListenAndServe(":http", certManager.HTTPHandler(nil)) }()
	//	_ = serv.ListenAndServeTLS("", "")
	//	return
	//}

	//cmn.AppStartTime = time.Now()

	//z.Info(cmn.AppStartTime.Format(cmn.AppStartTimeLayout))
	//_ = serv.ListenAndServe()
}
