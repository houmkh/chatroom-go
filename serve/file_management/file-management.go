package file_management

import (
	"chatroom/cmn"
	"chatroom/serve"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
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
		//Fn: user,

		Path: "/file_management",
		Name: "file_management",

		Developer: developer,
	})
}

const filePath = "D:/GoProject/chatroom/user-files"

func UploadFile(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	uid, _ := strconv.Atoi(r.FormValue("uid"))
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("fail to upload file")
		return
	}
	defer file.Close()
	fmt.Println("uid:", uid)
	index := strings.Index(header.Filename, ".")
	filetype := header.Filename[index:len(header.Filename)]
	newFilePath := filePath + "/" + r.FormValue("uid") + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + filetype
	newFile, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println("fail to create a file")
		panic(err.Error())
		return
	}
	_, err = io.Copy(newFile, file)
	defer newFile.Close()
	if err != nil {
		fmt.Println("fail to copy file")
		return
	}
	//filepath 是存储的名字 filename是用户给的文件名
	sqlstr := `insert into fileinfo(uid,filepath,filename) values ($1,$2,$3)`
	_, err = conn.Exec(context.Background(), sqlstr, uid, newFilePath, header.Filename)
	if err != nil {
		fmt.Println("fail to insert into fileinfo")
		return
	}
	fmt.Println("successfully upload file")

}

func ShowFiles(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	jsonMap := make(map[string]interface{})
	buf, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(buf, &jsonMap)
	if err != nil {
		fmt.Println("fail to unmarshal")
		return
	}
	uid, _ := strconv.Atoi(jsonMap["uid"].(string))
	sqlstr := `select filename,filepath,fid from fileinfo where uid = $1`
	resultSet, err := conn.Query(context.Background(), sqlstr, uid)
	if err != nil {
		fmt.Println("fail to read files")
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
	fmt.Println(fileList)
	buf, _ = json.Marshal(fileList)
	_, err = w.Write(buf)
	if err != nil {
		fmt.Println("fail to write back")
		return
	}
	fmt.Println("successfully send file list back")
}

func DownloadFile(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
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
		fmt.Println(err.Error())
		return
	}

	filebuf, err := ioutil.ReadFile(filepath)
	fmt.Println(filepath)
	if err != nil {
		fmt.Println("fail to read file")
		fmt.Println(err.Error())
		return
	}
	w.Write(filebuf)
}
