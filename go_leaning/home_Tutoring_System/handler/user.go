package handler

import (
	"home_Tutoring_System/database"
	"home_Tutoring_System/model"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"time"

	"home_Tutoring_System/config"
	"home_Tutoring_System/middleware"

	"github.com/golang-jwt/jwt/v5"
)

// RegisterRequest 对应前端注册时传过来的 JSON 数据。
// Gin 会根据 json tag，把请求体里的字段绑定到这个结构体里。
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,max=20"`
	Phone           string `json:"phone" binding:"required,len=11"`
	Password        string `json:"password" binding:"required,min=6"`
	Role            string `json:"role" binding:"required"`
	Subject         string `json:"subject"`
	TeacheringYears int    `json:"teaching_years"`
}

// LoginRequest 对应前端传来的json数据
type LoginRequest struct {
	Username string `json:"username" binding:"required_without=Phone"`
	Phone    string `json:"phone" binding:"required_without=Username"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateMeView struct {
	Username string `json:"new_name"`
}

// Register 用户注册接口。
// 学生只需要 username、phone、password、role。
// 老师除了这些字段，还必须传 subject 和 teaching_years。
func Register(c *gin.Context) {
	var req RegisterRequest

	// 1. 解析 JSON 请求体，并做基础校验。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
			"errors":  err.Error(),
		})
		return
	}

	// 2. 只允许注册 student 或 teacher。
	if req.Role != "student" && req.Role != "teacher" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "身份类型只能是 student 或 teacher",
		})
		return
	}

	// 3. 如果注册老师，必须填写授课科目和教龄。
	if req.Role == "teacher" {
		if req.Subject == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "教师必须填写授课科目",
			})
			return
		}

		if req.TeacheringYears < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "教龄不能小于 0",
			})
			return
		}
	}

	// 4. 检查用户名是否已经存在。
	var count int64
	database.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户名已经被占用了",
		})
		return
	}

	// 5. 检查手机号是否已经存在。
	database.DB.Model(&model.User{}).Where("phone = ?", req.Phone).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "手机号已经被注册了",
		})
		return
	}

	// 6. 密码不能明文存数据库，要先用 bcrypt 加密。
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "密码加密失败",
		})
		return
	}

	// 7. 创建用户。
	user := model.User{
		Username: req.Username,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "注册失败",
			"errors":  err.Error(),
		})
		return
	}

	// 8. 如果是老师，再创建老师档案。
	if req.Role == "teacher" {
		teacherProfile := model.TeacherProfile{
			UserID:          user.ID,
			Subject:         req.Subject,
			TeacheringYears: req.TeacheringYears,
		}

		if err := database.DB.Create(&teacherProfile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "教师信息创建失败",
				"errors":  err.Error(),
			})
			return
		}
	}

	// 9. 返回注册成功的信息。注意不要把密码返回给前端。
	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"phone":    user.Phone,
			"role":     user.Role,
		},
	})
}

func generateToken(user model.User) (string, error) {
	claims := middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(config.JWTSecret)
}

func Login(c *gin.Context) {
	var req LoginRequest

	// 1. 解析 JSON 请求体，并做基础校验。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
			"errors":  err.Error(),
		})
		return
	}

	var user model.User

	err := database.DB.Where("username= ? OR phone = ? ", req.Username, req.Phone).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(401, gin.H{"error": "账号不存在"})
		} else {
			c.JSON(500, gin.H{"error": "数据库错误"})
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "密码错误"})
		return
	}

	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "登陆成功",
		"data": gin.H{
			"token":    token,
			"id":       user.ID,
			"username": user.Username,
			"phone":    user.Phone,
			"role":     user.Role,
		},
	})

}

func MeView(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		userID := c.GetUint("userID")

		var user model.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(404, gin.H{"message": "用户不存在"})
			return
		}

		c.JSON(200, gin.H{
			"message": "获取个人身份信息成功",
			"data": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"phone":    user.Phone,
				"role":     user.Role,
			},
		})
		return
	}

	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
		var req UpdateMeView
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"message": "参数错误",
				"error":   err.Error(),
			})
			return
		}
		var user model.User
		userID := c.GetUint("userID")
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(404, gin.H{
				"message": "用户不存在",
				"error":   err.Error(),
			})
			return
		}
		if req.Username != "" {
			user.Username = req.Username
		}
		if err := database.DB.Save(&user).Error; err != nil {
			c.JSON(500, gin.H{"message": "修改失败"})
			return
		}
		c.JSON(200, gin.H{"message": "修改成功"})
		return
	}

	c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "请求方法不支持"})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "退出成功",
	})
}
