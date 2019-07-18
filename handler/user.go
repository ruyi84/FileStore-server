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
	if !pwdChecked {
		w.Write([]byte("FAILED:pwdChecker err"))
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
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	//
	////2验证token是否有效
	//isValidToekn := IstokenValid(token)
	//if isValidToekn {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	//3查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

//token是否有时效性
func IstokenValid(token string) bool {
	//TODO 判断token的时效性，是否过期
	if len(token) != 40 {
		return false
	}
	//TODO 从数据库tbl_user_token查询是否有当前token，
	//TODO 对比两个token是否一致

	return true
}

func GenToken(username string) string {
	//40位md5(usernmae +timestamp + token_solt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
