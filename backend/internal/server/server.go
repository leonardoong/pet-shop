package server

import (
	"fmt"
	"log"
	"strings"

	authhandler "petshop/internal/handler/auth"
	carthandler "petshop/internal/handler/cart"
	dashhandler "petshop/internal/handler/dashboard"
	orderhandler "petshop/internal/handler/order"
	paymenthandler "petshop/internal/handler/payment"
	producthandler "petshop/internal/handler/product"
	uploadhandler "petshop/internal/handler/upload"
	userhandler "petshop/internal/handler/user"
	"petshop/internal/middleware"
	authrepo "petshop/internal/repository/auth"
	cartrepo "petshop/internal/repository/cart"
	dashrepo "petshop/internal/repository/dashboard"
	orderrepo "petshop/internal/repository/order"
	productrepo "petshop/internal/repository/product"
	userrepo "petshop/internal/repository/user"
	authsvc "petshop/internal/service/auth"
	cartsvc "petshop/internal/service/cart"
	dashsvc "petshop/internal/service/dashboard"
	notificationsvc "petshop/internal/service/notification"
	ordersvc "petshop/internal/service/order"
	productsvc "petshop/internal/service/product"
	usersvc "petshop/internal/service/user"
	"petshop/pkg/config"
	"petshop/pkg/email"
	jwtpkg "petshop/pkg/jwt"
	"petshop/pkg/payment"

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
	addrRepo, adminUserRepo := userrepo.NewRepository(db)
	orderRepo, adminOrderRepo := orderrepo.NewRepository(db)
	dashboardRepo := dashrepo.NewRepository(db)

	// --- Services ---
	authService := authsvc.NewService(authRepo, jwtManager, s.cfg.JWT.AccessExpiryMinutes)
	categoryService := productsvc.NewCategoryService(categoryRepo)
	productService := productsvc.NewProductService(productRepo)
	cartService := cartsvc.NewService(cartRepo, productRepo)
	addrService := usersvc.NewAddressService(addrRepo)

	var paymentDrv payment.Driver
	if s.cfg.Payment.Driver == "mock" {
		paymentDrv = payment.NewMockDriver()
	}

	orderService := ordersvc.NewService(orderRepo, cartRepo, addrRepo, productRepo, paymentDrv)

	adminProductSvc := productsvc.NewAdminProductService(productRepo, categoryRepo)
	adminCategorySvc := productsvc.NewAdminCategoryService(categoryRepo, productRepo)
	inventorySvc := productsvc.NewInventoryService(productRepo, categoryRepo)
	adminOrderSvc := ordersvc.NewAdminService(adminOrderRepo)
	adminCustomerSvc := usersvc.NewAdminCustomerService(adminUserRepo)
	dashboardSvc := dashsvc.NewService(dashboardRepo)
	profileSvc := usersvc.NewProfileService(adminUserRepo, authRepo)
	resetSvc := authsvc.NewResetService(authRepo)
	reviewSvc := productsvc.NewReviewService(productRepo)

	emailClient := email.NewDebugClient()
	_ = notificationsvc.NewService(emailClient, "http://localhost:3000") // ready for email triggers

	// --- Handlers ---
	authH := authhandler.NewHandler(authService)
	categoryH := producthandler.NewCategoryHandler(categoryService)
	productH := producthandler.NewProductHandler(productService)
	cartH := carthandler.NewHandler(cartService)
	addrH := userhandler.NewAddressHandler(addrService)
	orderH := orderhandler.NewHandler(orderService)

	adminProductH := producthandler.NewAdminProductHandler(adminProductSvc)
	adminCategoryH := producthandler.NewAdminCategoryHandler(adminCategorySvc)
	inventoryH := producthandler.NewInventoryHandler(inventorySvc)
	adminOrderH := orderhandler.NewAdminOrderHandler(adminOrderSvc)
	adminCustomerH := userhandler.NewAdminCustomerHandler(adminCustomerSvc)
	dashboardH := dashhandler.NewHandler(dashboardSvc)
	paymentH := paymenthandler.NewHandler(orderRepo, adminOrderRepo)
	profileH := userhandler.NewProfileHandler(profileSvc)
	resetH := authhandler.NewResetHandler(resetSvc)
	reviewH := producthandler.NewReviewHandler(reviewSvc)
	uploadH := uploadhandler.NewHandler("./uploads", "http://localhost:8080/api/v1", 5*1024*1024)
	avatarH := userhandler.NewAvatarHandler(adminUserRepo, "./uploads/avatars", "http://localhost:8080/api/v1")

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

	v1.POST("/auth/forgot-password", resetH.ForgotPassword)
	v1.POST("/auth/reset-password", resetH.ResetPassword)
	v1.GET("/payment/mock/pay", paymentH.MockPay)

	// --- Public: Catalog ---
	v1.GET("/categories", categoryH.List)
	v1.GET("/categories/:slug", categoryH.GetBySlug)
	v1.GET("/products", productH.List)
	v1.GET("/products/:slug", productH.GetBySlug)
	v1.GET("/products/:slug/reviews", reviewH.ListByProduct)

	// --- Protected: Customer ---
	customer := v1.Group("/customer")
	customer.Use(middleware.RequireCustomer(jwtManager))
	{
		customer.GET("/me", profileH.Get)
		customer.PUT("/me", profileH.Update)
		customer.PUT("/me/email", profileH.ChangeEmail)
		customer.PUT("/me/password", profileH.ChangePassword)
		customer.POST("/me/avatar", avatarH.Upload)

		customer.POST("/reviews", reviewH.Create)

		customer.GET("/cart", cartH.Get)
		customer.POST("/cart/items", cartH.AddItem)
		customer.PUT("/cart/items/:productId", cartH.UpdateItem)
		customer.DELETE("/cart/items/:productId", cartH.RemoveItem)

		customer.GET("/addresses", addrH.List)
		customer.POST("/addresses", addrH.Create)
		customer.PUT("/addresses/:id", addrH.Update)
		customer.DELETE("/addresses/:id", addrH.Delete)
		customer.PATCH("/addresses/:id/default", addrH.SetDefault)

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

		admin.POST("/upload", uploadH.Upload, middleware.RequirePermission("products:create"))

		admin.GET("/dashboard", dashboardH.GetStats, middleware.RequirePermission("dashboard:read"))

		adminProducts := admin.Group("/products")
		adminProducts.Use(middleware.RequirePermission("products:read"))
		{
			adminProducts.GET("", adminProductH.List)
			adminProducts.GET("/:id", adminProductH.Get)
		}
		adminProductsCreate := admin.Group("/products")
		adminProductsCreate.Use(middleware.RequirePermission("products:create"))
		{
			adminProductsCreate.POST("", adminProductH.Create)
		}
		adminProductsUpdate := admin.Group("/products")
		adminProductsUpdate.Use(middleware.RequirePermission("products:update"))
		{
			adminProductsUpdate.PUT("/:id", adminProductH.Update)
		}
		adminProductsDelete := admin.Group("/products")
		adminProductsDelete.Use(middleware.RequirePermission("products:delete"))
		{
			adminProductsDelete.DELETE("/:id", adminProductH.Delete)
		}

		adminCategories := admin.Group("/categories")
		adminCategories.Use(middleware.RequirePermission("categories:read"))
		{
			adminCategories.GET("", categoryH.List)
		}
		adminCategoriesCreate := admin.Group("/categories")
		adminCategoriesCreate.Use(middleware.RequirePermission("categories:create"))
		{
			adminCategoriesCreate.POST("", adminCategoryH.Create)
		}
		adminCategoriesUpdate := admin.Group("/categories")
		adminCategoriesUpdate.Use(middleware.RequirePermission("categories:update"))
		{
			adminCategoriesUpdate.PUT("/:id", adminCategoryH.Update)
		}
		adminCategoriesDelete := admin.Group("/categories")
		adminCategoriesDelete.Use(middleware.RequirePermission("categories:delete"))
		{
			adminCategoriesDelete.DELETE("/:id", adminCategoryH.Delete)
		}

		adminOrders := admin.Group("/orders")
		adminOrders.Use(middleware.RequirePermission("orders:read"))
		{
			adminOrders.GET("", adminOrderH.List)
			adminOrders.GET("/:id", adminOrderH.Get)
		}
		adminOrdersStatus := admin.Group("/orders")
		adminOrdersStatus.Use(middleware.RequirePermission("orders:update"))
		{
			adminOrdersStatus.PATCH("/:id/status", adminOrderH.UpdateStatus)
		}

		adminInventory := admin.Group("/inventory")
		adminInventory.Use(middleware.RequirePermission("inventory:read"))
		{
			adminInventory.GET("", inventoryH.List)
		}
		adminInventoryUpdate := admin.Group("/inventory")
		adminInventoryUpdate.Use(middleware.RequirePermission("inventory:update"))
		{
			adminInventoryUpdate.PATCH("/:productId", inventoryH.AdjustStock)
		}

		adminCustomers := admin.Group("/customers")
		adminCustomers.Use(middleware.RequirePermission("customers:read"))
		{
			adminCustomers.GET("", adminCustomerH.List)
			adminCustomers.GET("/:id", adminCustomerH.Get)
		}
		adminCustomersUpdate := admin.Group("/customers")
		adminCustomersUpdate.Use(middleware.RequirePermission("customers:update"))
		{
		adminCustomersUpdate.PATCH("/:id", adminCustomerH.ToggleActive)
		}

		adminReviews := admin.Group("/reviews")
		adminReviews.Use(middleware.RequirePermission("products:read"))
		{
			adminReviews.GET("", reviewH.AdminList)
			adminReviews.PATCH("/:id/approve", reviewH.AdminToggle)
		}
	}

	s.router.Static("/api/v1/uploads", "./uploads")

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
