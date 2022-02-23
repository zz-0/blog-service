package main

import (
	"blog-service/global"
	"blog-service/internal/model"
	"blog-service/internal/routers"
	"blog-service/pkg/logger"
	"blog-service/pkg/setting"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	port    string
	runMode string
	config  string

	isVersion    bool
	buildTime    string
	buildVersion string
	gitCommitID  string
)

func init() {
	//输出设置初始化
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err:%v", err)
	}

	//输出数据库连接初始化日志
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err:%v", err)
	}

	//输出日志的初始化
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err:%v", err)
	}

	setupFlag()
}

// @title 博客系统
// @version 1.0
// @description 张震的第一个项目
// @termsOfService 储存仓库
func main() {

	if isVersion {
		fmt.Printf("build_time:%s\n", buildTime)
		fmt.Printf("build_Version:%s\n", buildVersion)
		fmt.Printf("git_commit_id:%s\n", gitCommitID)
		return
	}

	// 尝试程序段
	//logger:输出请求日志，并标准化日志的格式
	//Recovery:异常捕获,防止出现panic导致服务崩溃
	// r := gin.Default()  //返回路由引擎
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "pong"})
	// })
	// r.Run()
	gin.SetMode(gin.ReleaseMode) //设置程序的运行模式,不同模式对应不同的场景，默认是debug模式
	gin.SetMode(global.ServerSetting.RunMode)

	//路由启动接入
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

//配置初始化(详情见yaml文件)
func setupSetting() error {
	// setting, err := setting.NewSetting()
	s, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}

	//服务配置初始化
	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	//应用配置初始化
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	//数据库配置初始化
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	//鉴权配置初始化
	err = s.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	//邮箱配置初始化
	err = s.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	return nil
}

//数据库连接初始化
func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

//日志初始化
func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,  //日志文件所允许占用的最大空间为600mb
		MaxAge:    10,   //最长保留时间是10天
		LocalTime: true, //时间格式为本地时间
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}
