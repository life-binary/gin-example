package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"time"
)

type Test struct {
	Ip         string
	Port       int
	RemotePort int `yaml:"remote_port"`
}

type Person struct {
	Name    string `json:"name" form:"name" binding:"required"`
	Address string `json:"address" form:"address"`
}

func startPage(c *gin.Context) {
	var person Person

	// return error
	// ShouldBindQuery

	if err := c.ShouldBindJSON(&person); err == nil {
		log.Println("====== Only Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
		c.JSON(http.StatusOK, gin.H{"status": "Success"})
	} else {
		log.Printf("error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Failed"})
	}
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		log.Println(ip)

		// 每一个ip， 允许速率， 最大突发
		bucket := 100 // 令牌数
		key := ip
		interval := time.Second / 1 // 时间

		bol := GetToken(key, bucket, interval)
		if !bol {
			fmt.Printf("[RateLimit] rate limit trigger, key: %v", key)
			c.JSON(http.StatusBadRequest, gin.H{"message": "非法请求"})
			c.Abort()
			return
		}

		// handler请求处理函数
		c.Next()
	}
}

func main() {

	data0, err := yaml.Marshal(&Test{Ip: "0.0.0.0", Port: 80, Remote_Port: 8888})
	fmt.Printf("data:%s, err:%v\n", string(data0), err)

	str := `ip: 0.0.0.0
port: 80
remote_port: 8888`

	var t Test
	err = yaml.Unmarshal([]byte(str), &t)
	fmt.Printf("data:%v err:%v\n", t, err)

	for i := 0; i < 10; i++ {
	fla:
		for j := 0; j < 10; j++ {
			fmt.Printf("i=%d, j=%d\n", i, j)
			break fla
		}
	}

	InitLimit()

	gin.ForceConsoleColor()

	// Creates a router without any middleware by default
	r := gin.Default()
	r.Use(RateLimit())

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	//r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.GET("/panic", func(c *gin.Context) {
		// panic with a string -- the custom middleware could save this to a database or report it to the user
		panic("foo")
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world", "test": 2})
	})

	r.Any("/testing", startPage)

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
