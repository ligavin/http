package httphandler

import (
	. "http/comm"
	"os"
	"io"
	"regexp"
)

var fileSuffix *regexp.Regexp

var suffixWhiteList  map[string]interface{}

func checkUploadFile(handler Handler, path string)bool{

	fileSuffix,_ = regexp.Compile("^(./|D:\\down).*(.txt|.zip|.tgz|.exe)$")
	match := fileSuffix.MatchString(path)
	Debug(handler, "match :%v", match)
	return match
}

func Upload(handler Handler)(map[string]interface{}){

	r := handler.Request

	path := GetValue(handler, "path")

	if "" == path{
		return Result(-2,"path error")
	}

	file,fileHeader,err := r.FormFile("file")

	match := checkUploadFile(handler,path);

	if !match{
		return Result(-3, "path is invalid")
	}

	defer file.Close()

	if (err != nil){
		Debug(handler,"err:%s",err)
		return Result(-1,"interval error")
	}else {
		Debug(handler, "filename:%s,filesize:%vk", fileHeader.Filename, fileHeader.Size/1024)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		Error(handler,"create file error:%s",err)
	}

	written,err := io.Copy(f, file)

	Debug(handler,"written:%d,err:%v", written/1024, err)



	return Result(0,"ok")
}
