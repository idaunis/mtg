package metal

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Metal -framework MetalKit
#include "renderer2.h"
#include <simd/matrix.h>
*/
import "C"

import (
	"math"
)

func Matrix_multiply(a Matrix_float4x4, b Matrix_float4x4) Matrix_float4x4 {
	return Matrix_float4x4(C.simd_matrix_multiply(C.matrix_float4x4(a), C.matrix_float4x4(b)))
}

func Vector3_cross(a Vector_float3, b Vector_float3) Vector_float3 {
	return Vector_float3(C.simd_vector3_cross(C.vector_float3(a), C.vector_float3(b)))
}

func Vector4_normalize(a Vector_float4) Vector_float4 {
	return Vector_float4(C.simd_vector4_normalize(C.vector_float4(a)))
}

func NewMatrix_float4x4(m []Vector_float4) Matrix_float4x4 {
	return Matrix_float4x4(C.new_matrix_float4x4(
		C.vector_float4(m[0]),
		C.vector_float4(m[1]),
		C.vector_float4(m[2]),
		C.vector_float4(m[3]),
	))
}

func NewMatrix_float3x3(m []Vector_float3) Matrix_float3x3 {
	// fmt.Println("**", m[0], m[1], m[2])
	// since m[0] is an array of 3 floats we have to pass the location of the initial element
	return Matrix_float3x3(C.new_matrix_float3x3(
		&m[0][0],
		&m[1][0],
		&m[2][0],
	))
}

func Matrix_float4x4_translation(t Vector_float3) Matrix_float4x4 {
	X := C.vector_float4{1, 0, 0, 0}
	Y := C.vector_float4{0, 1, 0, 0}
	Z := C.vector_float4{0, 0, 1, 0}
	W := C.vector_float4{t[0], t[1], t[2], 1}

	return Matrix_float4x4(C.new_matrix_float4x4(X, Y, Z, W))
}

func Matrix_float4x4_uniform_scale(scale float32) Matrix_float4x4 {
	return Matrix_float4x4(C.new_matrix_float4x4(
		C.vector_float4{C.float(scale), 0, 0, 0},
		C.vector_float4{0, C.float(scale), 0, 0},
		C.vector_float4{0, 0, C.float(scale), 0},
		C.vector_float4{0, 0, 0, 1},
	))
}

func Matrix_float4x4_rotation(axis Vector_float3, angle float32) Matrix_float4x4 {
	c := C.float(math.Cos(float64(angle)))
	s := C.float(math.Sin(float64(angle)))

	X := C.vector_float4{
		axis[0]*axis[0] + (1-axis[0]*axis[0])*c,
		axis[0]*axis[1]*(1-c) - axis[2]*s,
		axis[0]*axis[2]*(1-c) + axis[1]*s,
		0.0,
	}

	Y := C.vector_float4{
		axis[0]*axis[1]*(1-c) + axis[2]*s,
		axis[1]*axis[1] + (1-axis[1]*axis[1])*c,
		axis[1]*axis[2]*(1-c) - axis[0]*s,
		0.0,
	}

	Z := C.vector_float4{
		axis[0]*axis[2]*(1-c) - axis[1]*s,
		axis[1]*axis[2]*(1-c) + axis[0]*s,
		axis[2]*axis[2] + (1-axis[2]*axis[2])*c,
		0.0,
	}

	W := C.vector_float4{0.0, 0.0, 0.0, 1.0}

	return Matrix_float4x4(C.new_matrix_float4x4(X, Y, Z, W))
}

func Matrix_float4x4_perspective(aspect, fovy, near, far float32) Matrix_float4x4 {
	yScale := float32(1 / math.Tan(float64(fovy)*0.5))
	xScale := yScale / aspect
	zRange := far - near
	zScale := -(far + near) / zRange
	wzScale := -2 * far * near / zRange

	return Matrix_float4x4(C.new_matrix_float4x4(
		C.vector_float4{C.float(xScale), 0, 0, 0},
		C.vector_float4{0, C.float(yScale), 0, 0},
		C.vector_float4{0, 0, C.float(zScale), -1},
		C.vector_float4{0, 0, C.float(wzScale), 0},
	))
}

func Matrix_float4x4_extract_linear(m Matrix_float4x4) Matrix_float3x3 {
	return NewMatrix_float3x3([]Vector_float3{
		Vector_float4(m.columns[0]).XYZ(),
		Vector_float4(m.columns[1]).XYZ(),
		Vector_float4(m.columns[2]).XYZ(),
	})
}
