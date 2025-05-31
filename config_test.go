package fork

import (
	"errors"
	"testing"
)

// MockConfigProvider for testing
type MockConfigProvider struct {
	data   map[string]interface{}
	errKey string // Key that should return error
}

func NewMockConfigProvider() *MockConfigProvider {
	return &MockConfigProvider{
		data: make(map[string]interface{}),
	}
}

func (m *MockConfigProvider) UnmarshalKey(key string, target interface{}) error {
	if key == m.errKey {
		return errors.New("mock error")
	}

	// Giả lập unmarshal thành công
	return nil
}

func (m *MockConfigProvider) SetError(key string) {
	m.errKey = key
}

func TestLoadConfigFromProvider(t *testing.T) {
	tests := []struct {
		name        string
		provider    interface{}
		key         string
		expectError bool
	}{
		{
			name:        "Valid provider and key",
			provider:    NewMockConfigProvider(),
			key:         "http",
			expectError: false,
		},
		{
			name: "Provider with unmarshal error",
			provider: func() *MockConfigProvider {
				p := NewMockConfigProvider()
				p.SetError("http")
				return p
			}(),
			key:         "http",
			expectError: true,
		},
		{
			name:        "Invalid provider",
			provider:    "invalid",
			key:         "http",
			expectError: false, // Should return default config
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadConfigFromProvider(tt.provider, tt.key)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if config == nil {
					t.Errorf("Expected config but got nil")
				}
			}
		})
	}
}

func TestNewWebAppWithConfig(t *testing.T) {
	tests := []struct {
		name        string
		provider    interface{}
		key         string
		expectError bool
	}{
		{
			name:        "Valid config provider",
			provider:    NewMockConfigProvider(),
			key:         "http",
			expectError: false,
		},
		{
			name:        "Empty config key",
			provider:    NewMockConfigProvider(),
			key:         "",
			expectError: false, // Should use default key "http"
		},
		{
			name:        "Invalid provider",
			provider:    "invalid",
			key:         "http",
			expectError: false, // Should work with default config
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewWebAppWithConfig(tt.provider, tt.key)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if app == nil {
					t.Errorf("Expected application but got nil")
				}

				// Check that config is set
				config := app.GetConfig()
				if config == nil {
					t.Errorf("Expected config to be set")
				}
			}
		})
	}
}

func TestReloadWebAppConfig(t *testing.T) {
	// Tạo application
	app := NewWebApp()
	provider := NewMockConfigProvider()

	// Test successful reload
	err := app.ReloadWebAppConfig(provider, "http")
	if err != nil {
		t.Errorf("Expected successful reload but got error: %v", err)
	}

	// Test reload with error
	provider.SetError("http")
	err = app.ReloadWebAppConfig(provider, "http")
	if err == nil {
		t.Errorf("Expected error during reload but got none")
	}

	// Test with empty key (should use default)
	provider.SetError("") // Reset error
	err = app.ReloadWebAppConfig(provider, "")
	if err != nil {
		t.Errorf("Expected successful reload with empty key but got error: %v", err)
	}
}

func TestWebAppConfigMerge(t *testing.T) {
	// Create base config
	baseConfig := DefaultWebAppConfig()

	// Create config to merge with graceful shutdown changes
	mergeConfig := &WebAppConfig{
		GracefulShutdown: GracefulShutdownConfig{
			Enabled:            false,
			Timeout:            60,
			WaitForConnections: false,
			SignalBufferSize:   2,
		},
	}

	// Merge configs
	baseConfig.MergeConfig(mergeConfig)

	// Check merged values
	if baseConfig.GracefulShutdown.Enabled != false {
		t.Errorf("Expected GracefulShutdown.Enabled to be false, got %v", baseConfig.GracefulShutdown.Enabled)
	}
	if baseConfig.GracefulShutdown.Timeout != 60 {
		t.Errorf("Expected GracefulShutdown.Timeout to be 60, got %d", baseConfig.GracefulShutdown.Timeout)
	}
	if baseConfig.GracefulShutdown.WaitForConnections != false {
		t.Errorf("Expected GracefulShutdown.WaitForConnections to be false, got %v", baseConfig.GracefulShutdown.WaitForConnections)
	}

	// Test merge with nil
	baseConfig.MergeConfig(nil)
	// Should not panic and values should remain the same
	if baseConfig.GracefulShutdown.Timeout != 60 {
		t.Errorf("Config should not change when merging with nil")
	}
}

func TestGracefulShutdownConfigMerge(t *testing.T) {
	baseConfig := &GracefulShutdownConfig{
		Enabled:            true,
		Timeout:            30,
		WaitForConnections: true,
		SignalBufferSize:   1,
	}

	mergeConfig := &GracefulShutdownConfig{
		Enabled:          false,
		Timeout:          60,
		SignalBufferSize: 2,
	}

	baseConfig.MergeConfig(mergeConfig)

	if baseConfig.Enabled != false {
		t.Errorf("Expected Enabled to be false, got %t", baseConfig.Enabled)
	}
	if baseConfig.Timeout != 60 {
		t.Errorf("Expected Timeout to be 60, got %d", baseConfig.Timeout)
	}
	if baseConfig.SignalBufferSize != 2 {
		t.Errorf("Expected SignalBufferSize to be 2, got %d", baseConfig.SignalBufferSize)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *WebAppConfig
		expectError bool
	}{
		{
			name:        "Valid config",
			config:      DefaultWebAppConfig(),
			expectError: false,
		},
		{
			name: "Invalid GracefulShutdown timeout",
			config: &WebAppConfig{
				GracefulShutdown: GracefulShutdownConfig{
					Timeout:          -1,
					SignalBufferSize: 1,
				},
			},
			expectError: true,
		},
		{
			name: "Invalid GracefulShutdown SignalBufferSize",
			config: &WebAppConfig{
				GracefulShutdown: GracefulShutdownConfig{
					Timeout:          30,
					SignalBufferSize: 0,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected validation error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error but got: %v", err)
				}
			}
		})
	}
}

func TestGracefulShutdownConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *GracefulShutdownConfig
		expectError bool
	}{
		{
			name: "Valid config",
			config: &GracefulShutdownConfig{
				Timeout:          30,
				SignalBufferSize: 1,
			},
			expectError: false,
		},
		{
			name: "Invalid timeout",
			config: &GracefulShutdownConfig{
				Timeout:          -1,
				SignalBufferSize: 1,
			},
			expectError: true,
		},
		{
			name: "Invalid SignalBufferSize",
			config: &GracefulShutdownConfig{
				Timeout:          30,
				SignalBufferSize: 0,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected validation error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error but got: %v", err)
				}
			}
		})
	}
}
