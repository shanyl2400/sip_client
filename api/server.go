package api

import (
	"fmt"
	"sipsimclient/config"
	"sipsimclient/service/devices"
	"sipsimclient/service/user"

	"github.com/gin-gonic/gin"
)

func Start() {
	server := NewServer(config.Get().HttpServerPort)
	server.Serve()
}

type Server struct {
	deviceManager devices.DeviceManager
	userService   user.User
	port          int
}

func (s *Server) Serve() {
	//初始化数据
	router := gin.Default()
	s.route(router)
	router.Run(fmt.Sprintf(":%v", s.port))
}
func (s *Server) route(r *gin.Engine) {
	//user
	r.POST("/user/login", s.login)
	r.POST("/user", s.registerUser)
	r.GET("/users", s.mustLogin, s.listUsers)
	r.PUT("/user/password", s.mustLogin, s.updateUserPassword)
	r.PUT("/user/password/reset/:user_name", s.mustAdmin, s.resetUserPassword)
	r.DELETE("/user/:user_name", s.mustAdmin, s.deleteUser)

	//device
	r.GET("/devices", s.mustLogin, s.listDevices)
	r.POST("/device", s.mustLogin, s.addDevice)
	r.PUT("/device", s.mustLogin, s.updateDevice)
	r.PUT("/device/connect/:device_name", s.mustLogin, s.connectDevice)
	r.PUT("/device/disconnect/:device_name", s.mustLogin, s.disconnectDevice)
	r.DELETE("/device/:device_name", s.mustLogin, s.deleteDevice)
	r.POST("/device/send/:device_name", s.mustLogin, s.sendMessage)
	r.GET("/device/logs/:device_name", s.mustLogin, s.deviceLogs)
}

func NewServer(port int) *Server {
	manager := devices.NewDeviceManager()
	user := user.NewUser()
	manager.Init()
	user.Init()
	return &Server{
		port:          port,
		userService:   user,
		deviceManager: manager,
	}
}
