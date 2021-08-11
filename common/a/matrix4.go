package a

type Matrix4 [16]float32

func (m Matrix4) MulVector(v Vector4) Vector4 {
	var result [4]float32
	vec := [4]float32 {v.X, v.Y, v.Z, v.W}
	for i := 0; i < 4; i++ {
		var sum float32
		for j := 0; j < 4; j++ {
			sum += m[i*4 + j] * vec[j]
		}
		result[i] = sum
	}
	return NewVector4(result[0], result[1], result[2], result[3])
}