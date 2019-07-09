package handler

import (
	"fmt"
	dblayer "github.com/filestore-server/db"
	"github.com/filestore-server/util"
	"io/ioutil"
	"net/http"
	"time"
)

const pwd_salt = "*#890"

//处理用户注册请求
func SignupHander(w http.ResponseWriter, r *http.Request) {
	fmt.Println(1)
	if r.Method == http.MethodGet {
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
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encpasswd := util.Sha1([]byte(password + pwd_salt))

	//1 校验用户名及密码
	pwdChecked := dblayer.UserSignin(username, encpasswd)
	if pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	//2 生成访问凭证token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}

	//3 登录成功后重定向到首页
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))

}

func GenToken(username string) string {
	//40位md5(usernmae +timestamp + token_solt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
