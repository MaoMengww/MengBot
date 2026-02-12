package utils

import (
	"errors"
	"math"
)

func F64ToF32(vec []float64) ([]float32, error) {
	if len(vec) == 0 {
		return nil, errors.New("vector is empty")
	}
	out := make([]float32, len(vec))
	for i, v := range vec {
		// 检查 NaN (Not a Number) 和 Inf (无穷大)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, errors.New("vector contains NaN or Inf")
		}
		out[i] = float32(v) // 强制类型转换
	}
	return out, nil
}
