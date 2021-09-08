package metal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVector3(t *testing.T) {
	v1 := Vector3(1, 2, 3)
	v2 := Vector3(4, 5, 6)

	c1 := Vector3_cross(v1, v2)
	assert.Equal(t, c1, Vector3(-3, 6, -3))

	a := Matrix_float4x4_extract_linear(NewMatrix_float4x4(
		[]Vector_float4{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		}))
	assert.Equal(t, a, NewMatrix_float3x3(
		[]Vector_float3{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		}))

	v4 := Vector4(1, 2, 3, 4)
	n4 := Vector4_normalize(v4)
	assert.Equal(t, n4, Vector4(0.18257417, 0.36514834, 0.5477225, 0.7302967))
}
