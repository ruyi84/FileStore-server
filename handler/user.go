package handler

import (
	"fmt"
	dblayer "github.com/filestore-server/db"
	"github.com/filestore-server/util"
	"io/ioutil"
	"net/http"
)

const pwd_salt = "*#890"

//处理用户注册请求
func SignupHander(w http.ResponseWriter, r *http.Request) {
	fmt.Println(1)
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))
	suc := dblayer.UserSignup(username, enc_passwd)
	if suc {
		w.Write([]byte("Succsess"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

//登录接口
func SignlnHandler(w http.ResponseWriter, r *http.Request) {
	//1 校验用户名及密码

	//2 生成访问凭证token

	//3 登录成功后重定向到首页
}
