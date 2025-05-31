package fork

import (
	"go.fork.vn/config"
	"go.fork.vn/di"
	"go.fork.vn/fork/adapter"
	"go.fork.vn/log"
)

// ServiceProvider là đối tượng thực hiện việc đăng ký và khởi tạo HTTP framework
// vào container dependency injection. Nó tuân theo interface ServiceProvider
// của package di để tích hợp với hệ thống service container.
type ServiceProvider struct{}

// NewServiceProvider tạo một instance mới của HTTP service provider.
// Provider này được sử dụng để đăng ký HTTP framework vào application container.
//
// Returns:
//   - di.ServiceProvider: Provider được sử dụng để đăng ký HTTP services
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký các binding liên quan đến HTTP framework vào container.
// Phương thức này tạo một HTTP WebApp instance và đăng ký nó vào container
// để các thành phần khác của ứng dụng có thể truy cập.
//
// Parameters:
//   - app: Instance của application chứa container DI
func (p *ServiceProvider) Register(app interface{}) {
	// Lấy container từ app
	if appWithContainer, ok := app.(interface {
		Container() *di.Container
	}); ok {
		c := appWithContainer.Container()

		// Đăng ký factory function để tạo HTTP WebApp
		c.Bind("http", func(container *di.Container) interface{} {
			// Tạo WebApp mặc định, config sẽ được load trong Boot()
			return NewWebApp()
		})

		// Đăng ký alias cho WebApp
		c.Alias("http.webapp", "http")
	}
}

// Boot được gọi sau khi tất cả các service provider đã được đăng ký.
// Phương thức này cấu hình HTTP WebApp bằng cách thiết lập adapter
// từ thông tin trong config và load WebApp config từ provider.
//
// Parameters:
//   - app: Instance của application chứa container DI
//
// Panics:
//   - Nếu không tìm thấy adapter trong config
//   - Nếu có lỗi trong việc load WebApp config
func (p *ServiceProvider) Boot(app interface{}) {
	if container, ok := app.(interface {
		Container() *di.Container
	}); ok {
		c := container.Container()
		// Lấy HTTP WebApp từ container
		httpApp := c.MustMake("http").(*WebApp)
		logger := c.MustMake("log").(log.Manager)
		// Lấy config manager từ container
		configManager := c.MustMake("config").(config.Manager)

		// Load WebApp config từ provider sử dụng LoadConfigFromProvider
		appConfig, err := LoadConfigFromProvider(configManager, "http")
		if err != nil {
			logger.Error("Failed to load HTTP WebApp config: " + err.Error())
			// Sử dụng config mặc định nếu có lỗi
			appConfig = DefaultWebAppConfig()
		}

		// Validate config trước khi set
		if err := appConfig.Validate(); err != nil {
			logger.Error("Invalid HTTP WebApp config: " + err.Error())
			// Sử dụng config mặc định nếu config không hợp lệ
			appConfig = DefaultWebAppConfig()
		}

		// Set config cho WebApp
		httpApp.SetConfig(appConfig)

		// Log thông tin config đã load
		logger.Info("HTTP WebApp config loaded successfully",
			"graceful_shutdown_enabled", appConfig.GracefulShutdown.Enabled,
			"graceful_shutdown_timeout", appConfig.GracefulShutdown.Timeout,
		)

		// Lấy tên adapter từ config
		adapterName, ok := configManager.GetString("http.adapter")
		if !ok {
			logger.Fatal("HTTP adapter not found in config")
		}

		// Lấy adapter instance từ container và thiết lập cho HTTP WebApp
		adapter := c.MustMake("http.adapter." + adapterName).(adapter.Adapter)
		if adapter == nil {
			logger.Fatal("HTTP adapter not found in container")
		} else {
			// Thiết lập adapter cho HTTP WebApp
			httpApp.SetAdapter(adapter)
			logger.Info("HTTP adapter set successfully", "adapter", adapterName)
		}

		// Setup graceful shutdown signal listening nếu được enable
		if appConfig.GracefulShutdown.Enabled {
			httpApp.ListenForShutdownSignals()
			logger.Info("Graceful shutdown enabled")
		}
	}
}

// Requires trả về danh sách các provider mà HTTP service provider phụ thuộc vào.
//
// Returns:
//   - []string: Mảng các tên providers được yêu cầu
func (p *ServiceProvider) Requires() []string {
	return []string{
		"log",    // Phụ thuộc vào log service
		"config", // Phụ thuộc vào config service
	}
}

// Providers trả về danh sách các service mà HTTP service provider đăng ký.
//
// Returns:
//   - []string: Mảng các tên services được đăng ký
func (p *ServiceProvider) Providers() []string {
	return []string{
		"http",        // HTTP WebApp chính
		"http.webapp", // Alias cho WebApp
	}
}
