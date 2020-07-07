package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"runtime"
	"strings"
)

func EncryptionArithmetic(username, token string) string {
	var signature string
	str := fmt.Sprintf("%s-%s", username, token)
	h := md5.New()
	h.Write([]byte(str))
	signature = hex.EncodeToString(h.Sum(nil))
	signature = strings.ToLower(signature)
	return signature
}

type ShellLogger struct{
	Username string
	TraceID  string
}

func (s ShellLogger) Info(msg ...interface{})  {

	if s.Username == ""{
		log.Printf("[unknown user] %s \n", msg)
		return
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}else{
		file = shortFile(file)
	}

	log.Printf("%s:%d [%s] [%s] %s ", file,line, s.TraceID, s.Username, fmt.Sprintln(msg...))
}

func shortFile(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short
}
