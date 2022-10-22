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
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {

	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("get ProblemBasic page parse error:", err)
		return
	}
	pageQuery := (page - 1) * size
	var count int64

	keyword := c.Query("keyword")

	list := make([]*models.ProblemBasic, 0)

	//tx := models.GetProblemList(keyword)
	//err = tx.Count(&count).Offset(pageQuery).Limit(size).Find(&list).Error

	err = models.DB.Model(new(models.ProblemBasic)).Where("title like ? or content like ?", "%"+keyword+"%", "%"+keyword+"%").
		Count(&count).Offset(pageQuery).Limit(size).Find(&list).Error

	if err != nil {
		log.Println("get problem list error:", err)
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
