package service

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"zyj.cn/helper"
	"zyj.cn/models"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "user_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {

	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户唯一标识不能为空",
		})
		return
	}

	data := new(models.UserBasic)
	err := models.DB.Where("identity=?", identity).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetUserDetail Error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户或密码为空",
		})
		return
	}
	password = helper.GetMd5(password)
	println(username, password)

	data := new(models.UserBasic)
	err := models.DB.Where("name=? AND password=?", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户或密码错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get userbasic error:" + err.Error(),
		})
		return
	}

	token, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GenerateToken error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})

}

// Login
// @Tags 公共方法
// @Summary 发送验证码
// @Param email formData string false "toUserEmail"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")

	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}

	//code := "123456"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := strconv.Itoa(r.Intn(1000000))//随机生成六位数验证码
	err := helper.SendCode(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "send code error:" + err.Error(),
		})
		return
	}
	err = models.RDB.Set(email, code, time.Second*60*3).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "sava code to redis error:" + err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": "发送成功",
	})
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Param mail formData string true "mail"
// @Param code formData string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /register [post]
func Register(c *gin.Context) {

	name := c.PostForm("name")
	password := c.PostForm("password")
	phone := c.PostForm("phone")
	_ = phone
	mail := c.PostForm("mail")
	code := c.PostForm("code")
	if code == "" || password == "" || name == "" || mail == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
	}

	//验证验证码是否正确
	codeRedis, err := models.RDB.Get(mail).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get code error:" + err.Error(),
		})
		return
	}
	if codeRedis != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "输入验证码错误",
		})
		return
	}

	user := new(models.UserBasic)

	err = models.DB.First(&user).Where("mail=?", mail).Error
	if err != nil {

		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get user error:" + err.Error(),
		})
		return
	}
	if user.Mail == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号已存在",
		})
		return
	}

	Identity := helper.GetUUID()
	data := &models.UserBasic{
		Identity: Identity,
		Name:     name,
		Password: helper.GetMd5(password),
		Phone:    phone,
		Mail:     mail,
	}

	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "create user error:" + err.Error(),
		})
		return
	}

	//生成token
	token, err := helper.GenerateToken(Identity, name, 0)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GenerateToken:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": token,
	})

}
