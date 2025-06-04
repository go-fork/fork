package fork

// WebAppConfig chứa các cấu hình bảo mật và hiệu suất cho WebApp
// Note: Some configurations have been moved to dedicated middleware packages:
// - MaxRequestBodySize -> bodylimit middleware
// - AllowedMethods -> method middleware
// - RequestTimeout -> timeout middleware
// - EnableSecurityHeaders -> helmet middleware
type WebAppConfig struct {
	// GracefulShutdown cấu hình graceful shutdown
	GracefulShutdown GracefulShutdownConfig `mapstructure:"graceful_shutdown" yaml:"graceful_shutdown"`
}

// GracefulShutdownConfig chứa cấu hình cho graceful shutdown
type GracefulShutdownConfig struct {
	// Enabled bật/tắt graceful shutdown
	// Mặc định: true
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`

	// Timeout thời gian tối đa để chờ graceful shutdown (seconds)
	// Mặc định: 30 seconds
	Timeout int `mapstructure:"timeout" yaml:"timeout"`

	// WaitForConnections có chờ tất cả connections kết thúc không
	// Mặc định: true
	WaitForConnections bool `mapstructure:"wait_for_connections" yaml:"wait_for_connections"`

	// SignalBufferSize kích thước buffer cho signal channel
	// Mặc định: 1
	SignalBufferSize int `mapstructure:"signal_buffer_size" yaml:"signal_buffer_size"`

	// OnShutdownStart callback được gọi khi bắt đầu shutdown
	OnShutdownStart func() `mapstructure:"-" yaml:"-"`

	// OnShutdownComplete callback được gọi khi shutdown hoàn thành
	OnShutdownComplete func() `mapstructure:"-" yaml:"-"`

	// OnShutdownError callback được gọi khi có lỗi trong quá trình shutdown
	OnShutdownError func(error) `mapstructure:"-" yaml:"-"`
}

// DefaultWebAppConfig trả về cấu hình mặc định cho WebApp
// Note: Middleware-specific configurations are now handled by their respective packages
func DefaultWebAppConfig() *WebAppConfig {
	return &WebAppConfig{
		GracefulShutdown: GracefulShutdownConfig{
			Enabled:            true,
			Timeout:            30, // 30 seconds
			WaitForConnections: true,
			SignalBufferSize:   1,
		},
	}
}

// MergeConfig hợp nhất cấu hình từ nhiều nguồn
// Note: Most configurations are now handled by middleware packages
func (c *WebAppConfig) MergeConfig(other *WebAppConfig) {
	if other == nil {
		return
	}

	c.GracefulShutdown.MergeConfig(&other.GracefulShutdown)
}

// MergeConfig hợp nhất cấu hình graceful shutdown
func (g *GracefulShutdownConfig) MergeConfig(other *GracefulShutdownConfig) {
	if other == nil {
		return
	}

	g.Enabled = other.Enabled

	if other.Timeout > 0 {
		g.Timeout = other.Timeout
	}

	g.WaitForConnections = other.WaitForConnections

	if other.SignalBufferSize > 0 {
		g.SignalBufferSize = other.SignalBufferSize
	}
}

// Validate kiểm tra tính hợp lệ của cấu hình
// Note: Most validations are now handled by middleware packages
func (c *WebAppConfig) Validate() error {
	return c.GracefulShutdown.Validate()
}

// Validate kiểm tra tính hợp lệ của cấu hình graceful shutdown
func (g *GracefulShutdownConfig) Validate() error {
	if g.Timeout < 0 {
		return ErrInvalidConfiguration
	}

	if g.SignalBufferSize < 1 {
		return ErrInvalidConfiguration
	}

	return nil
}
