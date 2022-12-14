package file_management

import (
	"chatroom/cmn"
	"chatroom/serve"
	"chatroom/serve/reply_msg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//annotation:file_management-service
//author:{"name":"file_management","tel":"15521212871","email":"jiaying.hou@qq.com"}

func Enroll(author string) {
	var developer *cmn.ModuleAuthor

	if author != "" {
		var d cmn.ModuleAuthor
		err := json.Unmarshal([]byte(author), &d)
		if err != nil {
			return
		}
		developer = &d
	}

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: UploadFile,

		Path: "/user/upload_file",
		Name: "user/upload_file",

		Developer: developer,
	})
	cmn.AddService(&cmn.ServeEndPoint{
		Fn: ShowFiles,

		Path: "/user/show_files",
		Name: "user/show_files",

		Developer: developer,
	})
	cmn.AddService(&cmn.ServeEndPoint{
		Fn: DownloadFile,

		Path: "/user/download_file",
		Name: "user/download_file",

		Developer: developer,
	})
}

const filePath = "D:/GoProject/chatroom/user-files"

var (
	msg serve.ReplyMsg
)

func UploadFile(w http.ResponseWriter, r *http.Request, conn *pgxpool.Conn) {
	uid, _ := strconv.Atoi(r.FormValue("uid"))
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("fail to upload file")
		return
	}
	defer file.Close()
	//fmt.Println("uid:", uid)
	index := strings.Index(header.Filename, ".")
	filetype := header.Filename[index:len(header.Filename)]
	newFilePath := filePath + "/" + r.FormValue("uid") + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + filetype
	newFile, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println("fail to create a file")
		//panic(err.Error())
		fmt.Println(err)
		return
	}
	_, err = io.Copy(newFile, file)
	defer newFile.Close()
	if err != nil {
		fmt.Println("fail to copy file")
		fmt.Println(err)
		return
	}
	//filepath 是存储的名字 filename是用户给的文件名
	sqlstr := `insert into fileinfo(uid,filepath,filename) values ($1,$2,$3)`
	_, err = conn.Exec(context.Background(), sqlstr, uid, newFilePath, header.Filename)
	if err != nil {
		fmt.Println("fail to insert into fileinfo")
		fmt.Println(err)
		return
	}
	//fmt.Println("successfully upload file")

}

func ShowFiles(w http.ResponseWriter, r *http.Request, conn *pgxpool.Conn) {
	jsonMap := make(map[string]interface{})
	buf, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(buf, &jsonMap)
	if err != nil {
		fmt.Println("fail to unmarshal")
		fmt.Println(err)
		return
	}
	uid, _ := strconv.Atoi(jsonMap["uid"].(string))
	sqlstr := `select filename,filepath,fid from fileinfo where uid = $1`
	resultSet, err := conn.Query(context.Background(), sqlstr, uid)
	if err != nil {
		fmt.Println("fail to read files")
		fmt.Println(err)
		return
	}
	err = resultSet.Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	var fileList []serve.File
	for resultSet.Next() {
		var filename string
		var filepath string
		var fid int
		err := resultSet.Scan(&filename, &filepath, &fid)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var file serve.File
		file.FileName = filename
		file.FilePath = filepath
		file.Fid = fid
		fileList = append(fileList, file)
	}
	//fmt.Println(fileList)
	msg = serve.ReplyMsg{ServeStatus: 200, ResponseMessage: "success", Data: fileList}

	buf, _ = json.Marshal(msg)
	//_, err = w.Write(buf)
	reply_msg.Response(w, &msg)
	if err != nil {
		fmt.Println("fail to write back")
		return
	}
	//fmt.Println("successfully send file list back")
}

func DownloadFile(w http.ResponseWriter, r *http.Request, conn *pgxpool.Conn) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("get request body fail")
		return
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(buf, &jsonMap)
	if err != nil {
		fmt.Println("unmarshal fail")
		return
	}
	fid := int(jsonMap["fid"].(float64))
	sqlstr := `select filepath from fileinfo where fid = $1`
	result := conn.QueryRow(context.Background(), sqlstr, fid)
	var filepath string
	err = result.Scan(&filepath)
	if err == pgx.ErrNoRows {
		fmt.Println(err)
		return
	}

	filebuf, err := ioutil.ReadFile(filepath)
	fmt.Println(filepath)
	if err != nil {
		fmt.Println("fail to read file")
		fmt.Println(err)
		return
	}
	w.Write(filebuf)
}
