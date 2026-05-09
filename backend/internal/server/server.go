package server

import (
	"fmt"
	"log"
	"strings"

	authhandler "petshop/internal/handler/auth"
	carthandler "petshop/internal/handler/cart"
	orderhandler "petshop/internal/handler/order"
	producthandler "petshop/internal/handler/product"
	userhandler "petshop/internal/handler/user"
	"petshop/internal/middleware"
	authrepo "petshop/internal/repository/auth"
	cartrepo "petshop/internal/repository/cart"
	orderrepo "petshop/internal/repository/order"
	productrepo "petshop/internal/repository/product"
	userrepo "petshop/internal/repository/user"
	authsvc "petshop/internal/service/auth"
	cartsvc "petshop/internal/service/cart"
	ordersvc "petshop/internal/service/order"
	productsvc "petshop/internal/service/product"
	usersvc "petshop/internal/service/user"
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
	// --- Repositories ---
	authRepo := authrepo.NewRepository(db)
	categoryRepo, productRepo := productrepo.NewRepository(db)
	cartRepo := cartrepo.NewRepository(db)
	addrRepo := userrepo.NewRepository(db)
	orderRepo := orderrepo.NewRepository(db)

	// --- Services ---
	authService := authsvc.NewService(authRepo, jwtManager, s.cfg.JWT.AccessExpiryMinutes)
	categoryService := productsvc.NewCategoryService(categoryRepo)
	productService := productsvc.NewProductService(productRepo)
	cartService := cartsvc.NewService(cartRepo, productRepo)
	addrService := usersvc.NewAddressService(addrRepo)
	orderService := ordersvc.NewService(orderRepo, cartRepo, addrRepo, productRepo)

	// --- Handlers ---
	authH := authhandler.NewHandler(authService)
	categoryH := producthandler.NewCategoryHandler(categoryService)
	productH := producthandler.NewProductHandler(productService)
	cartH := carthandler.NewHandler(cartService)
	addrH := userhandler.NewAddressHandler(addrService)
	orderH := orderhandler.NewHandler(orderService)

	v1 := s.router.Group("/api/v1")

	// --- Public: Auth ---
	customerAuth := v1.Group("/customer/auth")
	{
		customerAuth.POST("/register", authH.CustomerRegister)
		customerAuth.POST("/login", authH.CustomerLogin)
		customerAuth.POST("/refresh", authH.CustomerRefresh)
		customerAuth.POST("/logout", authH.Logout)
	}

	adminAuth := v1.Group("/admin/auth")
	{
		adminAuth.POST("/login", authH.AdminLogin)
		adminAuth.POST("/refresh", authH.AdminRefresh)
		adminAuth.POST("/logout", authH.Logout)
	}

	// --- Public: Catalog ---
	v1.GET("/categories", categoryH.List)
	v1.GET("/categories/:slug", categoryH.GetBySlug)
	v1.GET("/products", productH.List)
	v1.GET("/products/:slug", productH.GetBySlug)

	// --- Protected: Customer ---
	customer := v1.Group("/customer")
	customer.Use(middleware.RequireCustomer(jwtManager))
	{
		// Cart
		customer.GET("/cart", cartH.Get)
		customer.POST("/cart/items", cartH.AddItem)
		customer.PUT("/cart/items/:productId", cartH.UpdateItem)
		customer.DELETE("/cart/items/:productId", cartH.RemoveItem)

		// Addresses
		customer.GET("/addresses", addrH.List)
		customer.POST("/addresses", addrH.Create)
		customer.PUT("/addresses/:id", addrH.Update)
		customer.DELETE("/addresses/:id", addrH.Delete)
		customer.PATCH("/addresses/:id/default", addrH.SetDefault)

		// Orders
		customer.POST("/orders", orderH.Checkout)
		customer.GET("/orders", orderH.List)
		customer.GET("/orders/:id", orderH.GetByID)
	}

	// --- Protected: Admin ---
	admin := v1.Group("/admin")
	admin.Use(middleware.RequireAdmin(jwtManager))
	{
		admin.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"admin_id":    c.GetString(middleware.ContextAdminID),
				"permissions": c.MustGet(middleware.ContextAdminPerms),
			})
		})
	}

	// --- Misc ---
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (s *Server) Run() {
	addr := fmt.Sprintf(":%s", s.cfg.App.Port)
	log.Printf("server listening on %s", addr)
	if err := s.router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
