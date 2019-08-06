package handler

import (
	"fmt"
	"github.com/filestore-server/util"
	"math"
	"net/http"
	"strconv"
	"time"

	rPool "github.com/filestore-server/cache/redis"
)

//初始化信息
type MultiparUploadInfo struct {
	FileHash    string
	FileSize    int
	UploadId    string
	ChunkSize   int
	ChunkConunt int
}

//初始化分块上传

func InitialMultipartUploadHandle(w http.ResponseWriter, r *http.Request) {
	//1 解析用户请求信息
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params intvalid", nil).JSONBytes())
		return
	}
	//2 获得redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//3 生成分块上传的初始化信息
	upInfo := MultiparUploadInfo{
		FileHash:    filehash,
		FileSize:    filesize,
		UploadId:    username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:   5 * 1024 * 1024, //5MB
		ChunkConunt: int(math.Ceil(float64(filesize) / 5 * 1024 * 1024)),
	}

	//4	将初始化信息写入到rides缓存中
	rConn.Do("HSET", "MP_"+upInfo.UploadId, "chunkcount", upInfo.ChunkConunt)
	rConn.Do("HSET", "MP_"+upInfo.UploadId, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadId, "filesize", upInfo.FileSize)

	//5 将相应初始化数据返回到客户端
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())

}
