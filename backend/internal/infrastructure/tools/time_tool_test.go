package tools

import (
	"context"
	"strings"
	"testing"
)

func TestGetCurrentTimeTool(t *testing.T) {
	def, handler := GetCurrentTimeTool()

	if def.Name != "get_current_time" {
		t.Errorf("expected name 'get_current_time', got %s", def.Name)
	}
	if handler == nil {
		t.Fatal("expected non-nil handler")
	}

	ctx := context.Background()

	t.Run("default UTC timezone", func(t *testing.T) {
		result, err := handler(ctx, map[string]any{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Content["timezone"] != "UTC" {
			t.Errorf("expected timezone 'UTC', got %v", result.Content["timezone"])
		}
		datetime, ok := result.Content["datetime"].(string)
		if !ok || datetime == "" {
			t.Error("expected non-empty datetime string")
		}
	})

	t.Run("valid timezone", func(t *testing.T) {
		result, err := handler(ctx, map[string]any{"timezone": "Asia/Tokyo"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Content["timezone"] != "Asia/Tokyo" {
			t.Errorf("expected timezone 'Asia/Tokyo', got %v", result.Content["timezone"])
		}
		datetime := result.Content["datetime"].(string)
		if !strings.Contains(datetime, "+09:00") {
			t.Errorf("expected Tokyo timezone offset +09:00 in %s", datetime)
		}
	})

	t.Run("invalid timezone", func(t *testing.T) {
		_, err := handler(ctx, map[string]any{"timezone": "Invalid/Timezone"})
		if err == nil {
			t.Error("expected error for invalid timezone")
		}
	})

	t.Run("empty timezone defaults to UTC", func(t *testing.T) {
		result, err := handler(ctx, map[string]any{"timezone": ""})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Content["timezone"] != "UTC" {
			t.Errorf("expected timezone 'UTC', got %v", result.Content["timezone"])
		}
	})
}
