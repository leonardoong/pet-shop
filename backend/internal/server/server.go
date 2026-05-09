package server

import (
	"fmt"
	"log"
	"strings"

	"petshop/internal/auth"
	"petshop/internal/middleware"
	"petshop/pkg/config"
	jwtpkg "petshop/pkg/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
}

func New(cfg *config.Config, db *gorm.DB, jwtManager *jwtpkg.Manager) *Server {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// CORS
	origins := strings.Split(cfg.CORS.AllowedOrigins, ",")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	s := &Server{router: router, cfg: cfg}
	s.registerRoutes(db, jwtManager)
	return s
}

func (s *Server) registerRoutes(db *gorm.DB, jwtManager *jwtpkg.Manager) {
	// --- Auth ---
	authRepo := auth.NewRepository(db)
	authSvc := auth.NewService(authRepo, jwtManager, s.cfg.JWT.AccessExpiryMinutes)
	authHandler := auth.NewHandler(authSvc)

	v1 := s.router.Group("/api/v1")

	// Customer auth (public)
	customerAuth := v1.Group("/customer/auth")
	{
		customerAuth.POST("/register", authHandler.CustomerRegister)
		customerAuth.POST("/login", authHandler.CustomerLogin)
		customerAuth.POST("/refresh", authHandler.CustomerRefresh)
		customerAuth.POST("/logout", authHandler.Logout)
	}

	// Admin auth (public)
	adminAuth := v1.Group("/admin/auth")
	{
		adminAuth.POST("/login", authHandler.AdminLogin)
		adminAuth.POST("/refresh", authHandler.AdminRefresh)
		adminAuth.POST("/logout", authHandler.Logout)
	}

	// Protected customer routes (example — expanded in later phases)
	customerProtected := v1.Group("/customer")
	customerProtected.Use(middleware.RequireCustomer(jwtManager))
	{
		customerProtected.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"customer_id": c.GetString(middleware.ContextCustomerID)})
		})
	}

	// Protected admin routes (example — expanded in later phases)
	adminProtected := v1.Group("/admin")
	adminProtected.Use(middleware.RequireAdmin(jwtManager))
	{
		adminProtected.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"admin_id":    c.GetString(middleware.ContextAdminID),
				"permissions": c.MustGet(middleware.ContextAdminPerms),
			})
		})
	}

	// Health
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger UI — available at /swagger/index.html
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (s *Server) Run() {
	addr := fmt.Sprintf(":%s", s.cfg.App.Port)
	log.Printf("server listening on %s", addr)
	if err := s.router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
