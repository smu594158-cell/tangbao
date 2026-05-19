package main

import (
	"backend/internal/domain"
	"context"
	"log"
	"os"

	"backend/pkg/middleware"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	httpDelivery "backend/internal/delivery/http"
	"backend/internal/repository"
	"backend/internal/usecase"
	"backend/pkg/mcp"

	"github.com/joho/godotenv"
)

func initChatModel(baseURL string) model.ChatModel {
	chatModel, err := ollama.NewChatModel(context.Background(), &ollama.ChatModelConfig{
		BaseURL: baseURL,
		Model:   "qwen:4b",
	})
	if err != nil {
		log.Printf("Warning: failed to init ollama chat model: %v (AI features disabled)", err)
		return nil
	}
	return chatModel
}

func main() {
	// 0. 加载 .env 环境变量文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 1. 初始化数据库连接 (依赖docker-compose 启动mysql)
	dsn := getEnvOrDefault("DB_DSN", "hztour_user:hztour_password@tcp(127.0.0.1:3306)/hztour_db?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移并初始化种子数据
	db.AutoMigrate(&domain.User{}, &domain.ChatHistory{}, &domain.Attraction{}, &domain.GeneratedText{}, &domain.AdminLog{})
	seedAttractions(db)

	// 2. 依赖注入 (Dependency Injection)
	// 2.1 实例化Repository
	chatRepo := repository.NewMysqlChatRepository(db)
	userRepo := repository.NewMysqlUserRepository(db)
	tourRepo := repository.NewMySQLTourRepository(db)
	adminLogRepo := repository.NewMysqlAdminLogRepository(db)

	// 2.2 实例化UseCase
	jwtSecret := getEnvOrDefault("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Println("WARNING: JWT_SECRET not set, using hardcoded default (insecure!)")
		jwtSecret = "super_secret_key_change_me_in_prod"
	}
	authUC := usecase.NewAuthUseCase(userRepo, jwtSecret)
	userUC := usecase.NewUserUseCase(userRepo)

	// 假设本地 Ollama 跑在 11434 端口
	ollamaBaseURL := getEnvOrDefault("OLLAMA_BASE_URL", "http://127.0.0.1:11434")
	llmModel := initChatModel(ollamaBaseURL)
	chatUC := usecase.NewChatUseCase(chatRepo, llmModel)

	tourUC := usecase.NewTourUsecase(tourRepo, llmModel)

	// 2.3 实例化Geo MCP Server & UseCase
	amapKey := getEnvOrDefault("AMAP_WEB_KEY", "f02253e396de81301f40da343d8bb16b")
	mcpServer := mcp.NewMCPServer(amapKey)
	geoUC := usecase.NewGeoUseCase(mcpServer)

	// 3. 初始化Gin 引擎
	r := gin.Default()
	r.Use(middleware.RateLimit())

	// 添加简单的健康检查接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 4. 注册路由
	authMiddleware := middleware.JWTAuth(jwtSecret)
	adminMiddleware := middleware.AdminAuth()
	adminLogMiddleware := middleware.AdminOperationLog(adminLogRepo)

	httpDelivery.NewAuthHandler(r, authUC)
	httpDelivery.NewCaptchaHandler(r)
	httpDelivery.NewChatHandler(r, chatUC, authMiddleware)
	httpDelivery.NewGeoHandler(r, geoUC) // 暂不加鉴权以便于测试
	httpDelivery.NewTourHandler(r, tourUC, authMiddleware)

	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(authMiddleware, adminMiddleware, adminLogMiddleware)
	httpDelivery.NewAdminHandler(adminGroup, userUC, tourUC)

	// 5. 启动 HTTP 服务
	port := getEnvOrDefault("SERVER_PORT", ":8080")
	log.Printf("Server is running on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func seedAttractions(db *gorm.DB) {
	var count int64
	db.Model(&domain.Attraction{}).Count(&count)
	if count == 0 {
		attractions := []domain.Attraction{
			{Name: "西湖", Description: "西湖，位于浙江省杭州市西湖区，是中国大陆首批国家重点风景名胜区和中国十大风景名胜之一。", Address: "浙江省杭州市西湖区", HeatLevel: 100},
			{Name: "灵隐寺", Description: "灵隐寺，又名云林寺，位于浙江省杭州市，背靠北高峰，面朝飞来峰，始建于东晋咸和元年（326年），占地面积约87000平方米。", Address: "浙江省杭州市西湖区法云弄1号", HeatLevel: 95},
			{Name: "千岛湖", Description: "千岛湖，即新安江水库，位于浙江省杭州市淳安县境内，小部分连接杭州市建德市西北。", Address: "浙江省杭州市淳安县", HeatLevel: 90},
			{Name: "雷峰塔", Description: "雷峰塔，又名皇妃塔、西关砖塔，位于浙江省杭州市西湖风景区南岸夕照山的雷峰上。", Address: "浙江省杭州市西湖区南山路15号", HeatLevel: 88},
		}
		for _, a := range attractions {
			db.Create(&a)
		}
		log.Println("Seeded initial attractions data.")
	}
}
