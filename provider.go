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
//
// Panics:
//   - Nếu app là nil
//   - Nếu container là nil
//   - Nếu không thể đăng ký binding hoặc alias
func (p *ServiceProvider) Register(app di.Application) {
	// Kiểm tra app không được nil
	if app == nil {
		panic("fork.ServiceProvider.Register: application cannot be nil")
	}

	// Lấy container từ app
	c := app.Container()
	if c == nil {
		panic("fork.ServiceProvider.Register: container cannot be nil")
	}

	// Đăng ký factory function để tạo HTTP WebApp
	c.Bind("http", func(container di.Container) interface{} {
		// Tạo WebApp mặc định, config sẽ được load trong Boot()
		return NewWebApp()
	})

	// Đăng ký alias cho WebApp
	c.Alias("http.webapp", "http")
}

// Boot được gọi sau khi tất cả các service provider đã được đăng ký.
// Phương thức này cấu hình HTTP WebApp bằng cách thiết lập adapter
// từ thông tin trong config và load WebApp config từ provider.
//
// Parameters:
//   - app: Instance của application chứa container DI
//
// Panics:
//   - Nếu app là nil
//   - Nếu container là nil
//   - Nếu không tìm thấy required services (http, log, config)
//   - Nếu type assertion thất bại cho services
//   - Nếu không tìm thấy adapter trong config
//   - Nếu config loading hoặc validation thất bại (handled by LoadConfigFromProvider)
//   - Nếu unmarshal config thất bại (handled by LoadConfigFromProvider)
func (p *ServiceProvider) Boot(app di.Application) {
	// Kiểm tra app không được nil
	if app == nil {
		panic("fork.ServiceProvider.Boot: application cannot be nil")
	}

	// Lấy container từ app
	c := app.Container()
	if c == nil {
		panic("fork.ServiceProvider.Boot: container cannot be nil")
	}

	// Lấy HTTP WebApp từ container với kiểm tra type assertion
	httpService := c.MustMake("http")
	if httpService == nil {
		panic("fork.ServiceProvider.Boot: http service not found in container")
	}
	httpApp, ok := httpService.(*WebApp)
	if !ok {
		panic("fork.ServiceProvider.Boot: http service is not a *WebApp type")
	}

	// Lấy logger từ container với kiểm tra type assertion
	logService := c.MustMake("log")
	if logService == nil {
		panic("fork.ServiceProvider.Boot: log service not found in container")
	}
	logger, ok := logService.(log.Manager)
	if !ok {
		panic("fork.ServiceProvider.Boot: log service is not a log.Manager type")
	}

	// Lấy config manager từ container với kiểm tra type assertion
	configService := c.MustMake("config")
	if configService == nil {
		panic("fork.ServiceProvider.Boot: config service not found in container")
	}
	configManager, ok := configService.(config.Manager)
	if !ok {
		panic("fork.ServiceProvider.Boot: config service is not a config.Manager type")
	}

	// Tạo config mặc định
	appConfig := DefaultWebAppConfig()

	// Thực hiện unmarshal với error handling
	if err := configManager.UnmarshalKey("http", appConfig); err != nil {
		panic("fork.ServiceProvider.Boot: failed to unmarshal http config: " + err.Error())
	}

	// Validate config sau khi unmarshal
	if err := appConfig.Validate(); err != nil {
		panic("fork.ServiceProvider.Boot: failed to validate http config: " + err.Error())
	}
	// Set config cho WebApp
	httpApp.SetConfig(appConfig)

	// Log thông tin config đã load
	logger.Info("HTTP WebApp config loaded successfully",
		"graceful_shutdown_enabled", appConfig.GracefulShutdown.Enabled,
		"graceful_shutdown_timeout", appConfig.GracefulShutdown.Timeout,
	)

	// Lấy tên adapter từ config - bắt buộc phải có
	adapterName, ok := configManager.GetString("http.adapter")
	if !ok {
		logger.Fatal("HTTP adapter not found in config")
		panic("fork.ServiceProvider.Boot: http.adapter not found in config")
	}
	if adapterName == "" {
		logger.Fatal("HTTP adapter name is empty in config")
		panic("fork.ServiceProvider.Boot: http.adapter name is empty in config")
	}

	// Lấy adapter instance từ container và thiết lập cho HTTP WebApp
	adapterKey := "http.adapter." + adapterName
	adapterService := c.MustMake(adapterKey)
	if adapterService == nil {
		logger.Fatal("HTTP adapter not found in container: " + adapterKey)
		panic("fork.ServiceProvider.Boot: HTTP adapter not found in container: " + adapterKey)
	}

	// Type assertion cho adapter
	adapterInstance, ok := adapterService.(adapter.Adapter)
	if !ok {
		logger.Fatal("HTTP adapter is not of type adapter.Adapter: " + adapterKey)
		panic("fork.ServiceProvider.Boot: HTTP adapter is not of type adapter.Adapter: " + adapterKey)
	}

	// Thiết lập adapter cho HTTP WebApp
	httpApp.SetAdapter(adapterInstance)
	logger.Info("HTTP adapter set successfully", "adapter", adapterName)

	// Setup graceful shutdown signal listening nếu được enable
	if appConfig.GracefulShutdown.Enabled {
		httpApp.ListenForShutdownSignals()
		logger.Info("Graceful shutdown enabled")
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
