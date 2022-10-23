package service

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"zyj.cn/define"
	"zyj.cn/models"
)

// GetProblemList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("get SubmitBasic page parse error:", err)
		return
	}
	pageQuery := (page - 1) * size
	var count int64

	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	stauts, _ := strconv.Atoi(c.Query("stauts"))

	list := make([]*models.SubmitBasic, 0)

	tx := models.GetSubmitList(problemIdentity, userIdentity, stauts)
	err = tx.Count(&count).Offset(pageQuery).Limit(size).Find(&list).Error

	if err != nil {
		log.Println("get submit list error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get submit list error:" + err.Error(),
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
