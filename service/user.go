package service

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"zyj.cn/define"
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
			"msg":  "get userBasic error:" + err.Error(),
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

// SendCode
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
	code := strconv.Itoa(r.Intn(1000000)) //随机生成六位数验证码
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
	var cnt int64
	err = models.DB.Debug().Where("mail = ?", mail).Model(&user).Count(&cnt).Error
	if err != nil {

		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get user error:" + err.Error(),
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已被注册",
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
	token, err := helper.GenerateToken(Identity, name, data.IsAdmin)
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

// GetRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {

	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("get ProblemBasic page parse error:", err)
		return
	}
	pageQuery := (page - 1) * size
	var count int64
	list := make([]*models.UserBasic, 0)
	err = models.DB.Model(new(models.UserBasic)).Count(&count).Order("pass_num DESC,submit_num ASC").
		Offset(pageQuery).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("get rank list error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get rank list error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": count,
			"page":  page,
			"size":  size,
			"data":  list,
		},
	})

}
