package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage(t *testing.T) {
	tests := []struct {
		name          string
		updateFunc    func(*MemStorage)
		getFunc       func(*MemStorage) (interface{}, bool)
		expectedValue interface{}
	}{
		{
			name: "UpdateGauge",
			updateFunc: func(s *MemStorage) {
				s.UpdateGauge("test_gauge", 10.0)
			},
			getFunc: func(s *MemStorage) (interface{}, bool) {
				return s.GetGauge("test_gauge")
			},
			expectedValue: 10.0,
		},
		{
			name: "UpdateCounter",
			updateFunc: func(s *MemStorage) {
				s.UpdateCounter("test_counter", 5)
			},
			getFunc: func(s *MemStorage) (interface{}, bool) {
				return s.GetCounter("test_counter")
			},
			expectedValue: int64(5),
		},
	}

	storage := NewMemStorage()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.updateFunc(storage)

			value, exists := test.getFunc(storage)

			assert.True(t, exists)
			assert.Equal(t, test.expectedValue, value)
		})
	}
}
