package fork_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	configMocks "go.fork.vn/config/mocks"
	diMocks "go.fork.vn/di/mocks"
	"go.fork.vn/fork"
	forkMocks "go.fork.vn/fork/mocks"
	logMocks "go.fork.vn/log/mocks"
)

// TestNewServiceProvider kiểm tra tạo ServiceProvider instance
func TestNewServiceProvider(t *testing.T) {
	provider := fork.NewServiceProvider()

	assert.NotNil(t, provider)
	assert.IsType(t, &fork.ServiceProvider{}, provider)
}

// TestServiceProvider_Requires kiểm tra danh sách dependencies
func TestServiceProvider_Requires(t *testing.T) {
	provider := &fork.ServiceProvider{}

	requires := provider.Requires()

	assert.Len(t, requires, 2)
	assert.Contains(t, requires, "log")
	assert.Contains(t, requires, "config")
}

// TestServiceProvider_Providers kiểm tra danh sách services được cung cấp
func TestServiceProvider_Providers(t *testing.T) {
	provider := &fork.ServiceProvider{}

	providers := provider.Providers()

	assert.Len(t, providers, 2)
	assert.Contains(t, providers, "http")
	assert.Contains(t, providers, "http.webapp")
}

// TestServiceProvider_Register kiểm tra đăng ký services
func TestServiceProvider_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		// Setup mocks
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)

		// Setup expectations
		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().Bind("http", mock.AnythingOfType("di.BindingFunc")).Return()
		mockContainer.EXPECT().Alias("http.webapp", "http").Return()

		// Test
		provider := &fork.ServiceProvider{}
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})
	})

	t.Run("panic when app is nil", func(t *testing.T) {
		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Register: application cannot be nil", func() {
			provider.Register(nil)
		})
	})

	t.Run("panic when container is nil", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockApp.EXPECT().Container().Return(nil)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Register: container cannot be nil", func() {
			provider.Register(mockApp)
		})
	})
}

// TestServiceProvider_Boot kiểm tra boot process
func TestServiceProvider_Boot(t *testing.T) {
	t.Run("successful boot with valid config", func(t *testing.T) {
		// Setup mocks
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		mockConfig := configMocks.NewMockManager(t)
		mockAdapter := forkMocks.NewMockAdapter(t)

		// Setup expectations for container calls
		mockApp.EXPECT().Container().Return(mockContainer)

		// Mock HTTP service
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)

		// Mock log service
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)

		// Mock config service
		mockContainer.EXPECT().MustMake("config").Return(mockConfig)

		// Mock config loading
		mockConfig.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Run(func(key string, target interface{}) {
				config := target.(*fork.WebAppConfig)
				config.GracefulShutdown.Enabled = true
				config.GracefulShutdown.Timeout = 30
			}).Return(nil)

		// Mock logging call for config loaded
		mockLogger.EXPECT().Info("HTTP WebApp config loaded successfully",
			"graceful_shutdown_enabled", true,
			"graceful_shutdown_timeout", 30).Return()

		// Mock adapter config
		mockConfig.EXPECT().GetString("http.adapter").Return("test", true)

		// Mock adapter loading
		mockContainer.EXPECT().MustMake("http.adapter.test").Return(mockAdapter)

		// Mock adapter SetHandler call
		mockAdapter.EXPECT().SetHandler(mock.Anything).Return()

		// Mock logging calls
		mockLogger.EXPECT().Info("HTTP adapter set successfully", "adapter", "test").Return()
		mockLogger.EXPECT().Info("Graceful shutdown enabled").Return()

		// Test
		provider := &fork.ServiceProvider{}
		assert.NotPanics(t, func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when app is nil", func(t *testing.T) {
		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: application cannot be nil", func() {
			provider.Boot(nil)
		})
	})

	t.Run("panic when container is nil", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockApp.EXPECT().Container().Return(nil)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: container cannot be nil", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when http service not found", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(nil)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: http service not found in container", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when http service wrong type", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		wrongType := "not a webapp"

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(wrongType)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: http service is not a *WebApp type", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when log service not found", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(nil)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: log service not found in container", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when log service wrong type", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		wrongType := "not a logger"

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(wrongType)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: log service is not a log.Manager type", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when config service not found", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(nil)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: config service not found in container", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when config service wrong type", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		wrongType := "not a config"

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(wrongType)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: config service is not a config.Manager type", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when config unmarshal fails", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		mockConfig := configMocks.NewMockManager(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(mockConfig)

		// Mock unmarshal error
		mockConfig.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Return(assert.AnError)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: failed to unmarshal http config: "+assert.AnError.Error(), func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when config validation fails", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		mockConfig := configMocks.NewMockManager(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(mockConfig)

		// Mock config with invalid values
		mockConfig.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Run(func(key string, target interface{}) {
				config := target.(*fork.WebAppConfig)
				config.GracefulShutdown.Timeout = -1 // Invalid value
			}).Return(nil)

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: failed to validate http config: "+fork.ErrInvalidConfiguration.Error(), func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when adapter not found in config", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		mockConfig := configMocks.NewMockManager(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(mockConfig)

		// Mock valid config
		mockConfig.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Run(func(key string, target interface{}) {
				config := target.(*fork.WebAppConfig)
				config.GracefulShutdown.Enabled = true
				config.GracefulShutdown.Timeout = 30
			}).Return(nil)

		// Mock logging call for config loaded
		mockLogger.EXPECT().Info("HTTP WebApp config loaded successfully",
			"graceful_shutdown_enabled", true,
			"graceful_shutdown_timeout", 30).Return()

		// Mock adapter not found
		mockConfig.EXPECT().GetString("http.adapter").Return("", false)
		mockLogger.EXPECT().Fatal("HTTP adapter not found in config").Return()

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: http.adapter not found in config", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when adapter name is empty", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		mockConfig := configMocks.NewMockManager(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(mockConfig)

		// Mock valid config
		mockConfig.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Run(func(key string, target interface{}) {
				config := target.(*fork.WebAppConfig)
				config.GracefulShutdown.Enabled = true
				config.GracefulShutdown.Timeout = 30
			}).Return(nil)

		// Mock logging call for config loaded
		mockLogger.EXPECT().Info("HTTP WebApp config loaded successfully",
			"graceful_shutdown_enabled", true,
			"graceful_shutdown_timeout", 30).Return()

		// Mock empty adapter name
		mockConfig.EXPECT().GetString("http.adapter").Return("", true)
		mockLogger.EXPECT().Fatal("HTTP adapter name is empty in config").Return()

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: http.adapter name is empty in config", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when adapter not found in container", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		mockConfig := configMocks.NewMockManager(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(mockConfig)

		// Mock valid config
		mockConfig.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Run(func(key string, target interface{}) {
				config := target.(*fork.WebAppConfig)
				config.GracefulShutdown.Enabled = true
				config.GracefulShutdown.Timeout = 30
			}).Return(nil)

		// Mock logging call for config loaded
		mockLogger.EXPECT().Info("HTTP WebApp config loaded successfully",
			"graceful_shutdown_enabled", true,
			"graceful_shutdown_timeout", 30).Return()

		// Mock adapter config
		mockConfig.EXPECT().GetString("http.adapter").Return("test", true)

		// Mock adapter not found in container
		mockContainer.EXPECT().MustMake("http.adapter.test").Return(nil)
		mockLogger.EXPECT().Fatal("HTTP adapter not found in container: http.adapter.test").Return()

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: HTTP adapter not found in container: http.adapter.test", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("panic when adapter wrong type", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		mockConfig := configMocks.NewMockManager(t)
		wrongType := "not an adapter"

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(mockConfig)

		// Mock valid config
		mockConfig.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Run(func(key string, target interface{}) {
				config := target.(*fork.WebAppConfig)
				config.GracefulShutdown.Enabled = true
				config.GracefulShutdown.Timeout = 30
			}).Return(nil)

		// Mock logging call for config loaded
		mockLogger.EXPECT().Info("HTTP WebApp config loaded successfully",
			"graceful_shutdown_enabled", true,
			"graceful_shutdown_timeout", 30).Return()

		// Mock adapter config
		mockConfig.EXPECT().GetString("http.adapter").Return("test", true)

		// Mock wrong adapter type
		mockContainer.EXPECT().MustMake("http.adapter.test").Return(wrongType)
		mockLogger.EXPECT().Fatal("HTTP adapter is not of type adapter.Adapter: http.adapter.test").Return()

		provider := &fork.ServiceProvider{}

		assert.PanicsWithValue(t, "fork.ServiceProvider.Boot: HTTP adapter is not of type adapter.Adapter: http.adapter.test", func() {
			provider.Boot(mockApp)
		})
	})

	t.Run("successful boot with graceful shutdown disabled", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		mockConfig := configMocks.NewMockManager(t)
		mockAdapter := forkMocks.NewMockAdapter(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(mockConfig)

		// Mock config with graceful shutdown disabled
		mockConfig.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Run(func(key string, target interface{}) {
				config := target.(*fork.WebAppConfig)
				config.GracefulShutdown.Enabled = false
				config.GracefulShutdown.Timeout = 30
			}).Return(nil)

		mockConfig.EXPECT().GetString("http.adapter").Return("test", true)
		mockContainer.EXPECT().MustMake("http.adapter.test").Return(mockAdapter)

		// Mock adapter SetHandler call
		mockAdapter.EXPECT().SetHandler(mock.Anything).Return()

		// Mock logging calls (no graceful shutdown enabled log)
		mockLogger.EXPECT().Info("HTTP WebApp config loaded successfully",
			"graceful_shutdown_enabled", false,
			"graceful_shutdown_timeout", 30).Return()
		mockLogger.EXPECT().Info("HTTP adapter set successfully", "adapter", "test").Return()
		// No graceful shutdown enabled log expected

		provider := &fork.ServiceProvider{}
		assert.NotPanics(t, func() {
			provider.Boot(mockApp)
		})
	})
}

// TestServiceProvider_Integration kiểm tra integration scenarios
func TestServiceProvider_Integration(t *testing.T) {
	t.Run("complete registration and boot cycle", func(t *testing.T) {
		// Setup mocks
		mockApp := diMocks.NewMockApplication(t)
		mockContainer := diMocks.NewMockContainer(t)
		mockWebApp := &fork.WebApp{}
		mockLogger := logMocks.NewMockManager(t)
		mockConfig := configMocks.NewMockManager(t)
		mockAdapter := forkMocks.NewMockAdapter(t)

		// Setup expectations for Register
		mockApp.EXPECT().Container().Return(mockContainer).Times(2) // Called in both Register and Boot
		mockContainer.EXPECT().Bind("http", mock.AnythingOfType("di.BindingFunc")).Return()
		mockContainer.EXPECT().Alias("http.webapp", "http").Return()

		// Setup expectations for Boot
		mockContainer.EXPECT().MustMake("http").Return(mockWebApp)
		mockContainer.EXPECT().MustMake("log").Return(mockLogger)
		mockContainer.EXPECT().MustMake("config").Return(mockConfig)

		mockConfig.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Run(func(key string, target interface{}) {
				config := target.(*fork.WebAppConfig)
				config.GracefulShutdown.Enabled = true
				config.GracefulShutdown.Timeout = 45
				config.GracefulShutdown.WaitForConnections = false
				config.GracefulShutdown.SignalBufferSize = 2
			}).Return(nil)

		mockConfig.EXPECT().GetString("http.adapter").Return("http", true)
		mockContainer.EXPECT().MustMake("http.adapter.http").Return(mockAdapter)

		// Mock adapter SetHandler call
		mockAdapter.EXPECT().SetHandler(mock.Anything).Return()

		mockLogger.EXPECT().Info("HTTP WebApp config loaded successfully",
			"graceful_shutdown_enabled", true,
			"graceful_shutdown_timeout", 45).Return()
		mockLogger.EXPECT().Info("HTTP adapter set successfully", "adapter", "http").Return()
		mockLogger.EXPECT().Info("Graceful shutdown enabled").Return()

		// Test complete cycle
		provider := &fork.ServiceProvider{}

		// Test Register
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		})

		// Test Boot
		assert.NotPanics(t, func() {
			provider.Boot(mockApp)
		})

		// Verify provider interface methods
		requires := provider.Requires()
		providers := provider.Providers()

		assert.Len(t, requires, 2)
		assert.Contains(t, requires, "log")
		assert.Contains(t, requires, "config")

		assert.Len(t, providers, 2)
		assert.Contains(t, providers, "http")
		assert.Contains(t, providers, "http.webapp")
	})
}

// BenchmarkServiceProvider_Register benchmark cho Registration
func BenchmarkServiceProvider_Register(b *testing.B) {
	mockApp := diMocks.NewMockApplication(b)
	mockContainer := diMocks.NewMockContainer(b)

	mockApp.EXPECT().Container().Return(mockContainer).Times(b.N)
	mockContainer.EXPECT().Bind("http", mock.AnythingOfType("di.BindingFunc")).Return().Times(b.N)
	mockContainer.EXPECT().Alias("http.webapp", "http").Return().Times(b.N)

	provider := &fork.ServiceProvider{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		provider.Register(mockApp)
	}
}

// BenchmarkServiceProvider_Requires benchmark cho Requires method
func BenchmarkServiceProvider_Requires(b *testing.B) {
	provider := &fork.ServiceProvider{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = provider.Requires()
	}
}

// BenchmarkServiceProvider_Providers benchmark cho Providers method
func BenchmarkServiceProvider_Providers(b *testing.B) {
	provider := &fork.ServiceProvider{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = provider.Providers()
	}
}
