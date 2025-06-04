package fork_test

import (
	"testing"

	"go.fork.vn/config/mocks"
	"go.fork.vn/fork"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestDefaultWebAppConfig kiểm tra cấu hình mặc định
func TestDefaultWebAppConfig(t *testing.T) {
	config := fork.DefaultWebAppConfig()

	assert.NotNil(t, config)
	assert.True(t, config.GracefulShutdown.Enabled)
	assert.Equal(t, 30, config.GracefulShutdown.Timeout)
	assert.True(t, config.GracefulShutdown.WaitForConnections)
	assert.Equal(t, 1, config.GracefulShutdown.SignalBufferSize)
}

// TestWebAppConfig_Validate kiểm tra validation cấu hình
func TestWebAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *fork.WebAppConfig
		wantErr bool
	}{
		{
			name:    "valid default config",
			config:  fork.DefaultWebAppConfig(),
			wantErr: false,
		},
		{
			name: "valid custom config",
			config: &fork.WebAppConfig{
				GracefulShutdown: fork.GracefulShutdownConfig{
					Enabled:            true,
					Timeout:            60,
					WaitForConnections: true,
					SignalBufferSize:   2,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid timeout negative",
			config: &fork.WebAppConfig{
				GracefulShutdown: fork.GracefulShutdownConfig{
					Enabled:            true,
					Timeout:            -1, // Invalid
					WaitForConnections: true,
					SignalBufferSize:   1,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid signal buffer size zero",
			config: &fork.WebAppConfig{
				GracefulShutdown: fork.GracefulShutdownConfig{
					Enabled:            true,
					Timeout:            30,
					WaitForConnections: true,
					SignalBufferSize:   0, // Invalid
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, fork.ErrInvalidConfiguration, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGracefulShutdownConfig_Validate kiểm tra validation graceful shutdown
func TestGracefulShutdownConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  fork.GracefulShutdownConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: fork.GracefulShutdownConfig{
				Enabled:            true,
				Timeout:            30,
				WaitForConnections: true,
				SignalBufferSize:   1,
			},
			wantErr: false,
		},
		{
			name: "valid disabled config",
			config: fork.GracefulShutdownConfig{
				Enabled:            false,
				Timeout:            0,
				WaitForConnections: false,
				SignalBufferSize:   1,
			},
			wantErr: false,
		},
		{
			name: "invalid negative timeout",
			config: fork.GracefulShutdownConfig{
				Timeout:          -5,
				SignalBufferSize: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid zero signal buffer",
			config: fork.GracefulShutdownConfig{
				Timeout:          30,
				SignalBufferSize: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, fork.ErrInvalidConfiguration, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestWebAppConfig_MergeConfig kiểm tra merge cấu hình
func TestWebAppConfig_MergeConfig(t *testing.T) {
	t.Run("merge with nil config", func(t *testing.T) {
		config := fork.DefaultWebAppConfig()
		original := *config // Copy để so sánh

		config.MergeConfig(nil)

		// Config không thay đổi khi merge với nil
		assert.Equal(t, original, *config)
	})

	t.Run("merge with other config", func(t *testing.T) {
		config := fork.DefaultWebAppConfig()
		other := &fork.WebAppConfig{
			GracefulShutdown: fork.GracefulShutdownConfig{
				Enabled:            false,
				Timeout:            60,
				WaitForConnections: false,
				SignalBufferSize:   5,
			},
		}

		config.MergeConfig(other)

		assert.False(t, config.GracefulShutdown.Enabled)
		assert.Equal(t, 60, config.GracefulShutdown.Timeout)
		assert.False(t, config.GracefulShutdown.WaitForConnections)
		assert.Equal(t, 5, config.GracefulShutdown.SignalBufferSize)
	})
}

// TestGracefulShutdownConfig_MergeConfig kiểm tra merge graceful shutdown config
func TestGracefulShutdownConfig_MergeConfig(t *testing.T) {
	t.Run("merge with nil config", func(t *testing.T) {
		config := &fork.GracefulShutdownConfig{
			Enabled:            true,
			Timeout:            30,
			WaitForConnections: true,
			SignalBufferSize:   1,
		}
		original := *config

		config.MergeConfig(nil)

		assert.Equal(t, original, *config)
	})

	t.Run("merge with other config", func(t *testing.T) {
		config := &fork.GracefulShutdownConfig{
			Enabled:            true,
			Timeout:            30,
			WaitForConnections: true,
			SignalBufferSize:   1,
		}

		other := &fork.GracefulShutdownConfig{
			Enabled:            false,
			Timeout:            45,
			WaitForConnections: false,
			SignalBufferSize:   3,
		}

		config.MergeConfig(other)

		assert.False(t, config.Enabled)
		assert.Equal(t, 45, config.Timeout)
		assert.False(t, config.WaitForConnections)
		assert.Equal(t, 3, config.SignalBufferSize)
	})

	t.Run("merge with zero timeout keeps original", func(t *testing.T) {
		config := &fork.GracefulShutdownConfig{
			Timeout: 30,
		}

		other := &fork.GracefulShutdownConfig{
			Timeout: 0, // Zero value, should not override
		}

		config.MergeConfig(other)

		assert.Equal(t, 30, config.Timeout) // Original value preserved
	})

	t.Run("merge with zero signal buffer keeps original", func(t *testing.T) {
		config := &fork.GracefulShutdownConfig{
			SignalBufferSize: 2,
		}

		other := &fork.GracefulShutdownConfig{
			SignalBufferSize: 0, // Zero value, should not override
		}

		config.MergeConfig(other)

		assert.Equal(t, 2, config.SignalBufferSize) // Original value preserved
	})
}

// TestConfigLoadingWithMocks kiểm tra loading config với mocks
func TestConfigLoadingWithMocks(t *testing.T) {
	t.Run("successful config loading", func(t *testing.T) {
		mockManager := mocks.NewMockManager(t)

		// Setup expectations for graceful shutdown config
		mockManager.EXPECT().UnmarshalKey("http", mock.MatchedBy(func(config *fork.WebAppConfig) bool {
			// Validate that we received a WebAppConfig pointer
			return config != nil
		})).Run(func(key string, target interface{}) {
			// Type assert to get the actual config
			config := target.(*fork.WebAppConfig)
			// Simulate loading config từ YAML
			config.GracefulShutdown.Enabled = false
			config.GracefulShutdown.Timeout = 45
			config.GracefulShutdown.WaitForConnections = false
			config.GracefulShutdown.SignalBufferSize = 2
		}).Return(nil)

		// Test config loading
		config := fork.DefaultWebAppConfig()
		err := mockManager.UnmarshalKey("http", config)

		assert.NoError(t, err)
		assert.False(t, config.GracefulShutdown.Enabled)
		assert.Equal(t, 45, config.GracefulShutdown.Timeout)
		assert.False(t, config.GracefulShutdown.WaitForConnections)
		assert.Equal(t, 2, config.GracefulShutdown.SignalBufferSize)

		// Validate config after loading
		assert.NoError(t, config.Validate())
	})

	t.Run("config loading with unmarshal error", func(t *testing.T) {
		mockManager := mocks.NewMockManager(t)

		// Setup expectation for error
		mockManager.EXPECT().UnmarshalKey("http", mock.Anything).Return(assert.AnError)

		config := fork.DefaultWebAppConfig()
		err := mockManager.UnmarshalKey("http", config)

		assert.Error(t, err)
	})
}

// TestConfigIntegrationWithYAML kiểm tra tích hợp với YAML config
func TestConfigIntegrationWithYAML(t *testing.T) {
	t.Run("realistic config scenario", func(t *testing.T) {
		mockManager := mocks.NewMockManager(t)

		// Mô phỏng config từ app.example.yaml
		mockManager.EXPECT().UnmarshalKey("http", mock.AnythingOfType("*fork.WebAppConfig")).
			Run(func(key string, target interface{}) {
				// Type assert to get the actual config
				config := target.(*fork.WebAppConfig)
				// Mô phỏng values từ configs/app.example.yaml
				config.GracefulShutdown.Enabled = true
				config.GracefulShutdown.Timeout = 30
				config.GracefulShutdown.WaitForConnections = true
				config.GracefulShutdown.SignalBufferSize = 1
			}).Return(nil)

		// Load và validate config
		config := fork.DefaultWebAppConfig()
		err := mockManager.UnmarshalKey("http", config)

		assert.NoError(t, err)
		assert.NoError(t, config.Validate())

		// Verify các giá trị match với app.example.yaml
		assert.True(t, config.GracefulShutdown.Enabled)
		assert.Equal(t, 30, config.GracefulShutdown.Timeout)
		assert.True(t, config.GracefulShutdown.WaitForConnections)
		assert.Equal(t, 1, config.GracefulShutdown.SignalBufferSize)
	})
}

// BenchmarkDefaultWebAppConfig benchmark cho tạo config mặc định
func BenchmarkDefaultWebAppConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fork.DefaultWebAppConfig()
	}
}

// BenchmarkWebAppConfig_Validate benchmark cho validation
func BenchmarkWebAppConfig_Validate(b *testing.B) {
	config := fork.DefaultWebAppConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

// BenchmarkWebAppConfig_MergeConfig benchmark cho merge config
func BenchmarkWebAppConfig_MergeConfig(b *testing.B) {
	config := fork.DefaultWebAppConfig()
	other := &fork.WebAppConfig{
		GracefulShutdown: fork.GracefulShutdownConfig{
			Enabled:            false,
			Timeout:            60,
			WaitForConnections: false,
			SignalBufferSize:   2,
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		config.MergeConfig(other)
	}
}
