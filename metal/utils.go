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

func NewMatrix_float4x4(m []Vector_float4) Matrix_float4x4 {
	return Matrix_float4x4(C.new_matrix_float4x4(
		C.vector_float4(m[0]),
		C.vector_float4(m[1]),
		C.vector_float4(m[2]),
		C.vector_float4(m[3]),
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

/*

matrix_float4x4 matrix_float4x4_translation(vector_float3 t)
{
    vector_float4 X = { 1, 0, 0, 0 };
    vector_float4 Y = { 0, 1, 0, 0 };
    vector_float4 Z = { 0, 0, 1, 0 };
    vector_float4 W = { t.x, t.y, t.z, 1 };

    matrix_float4x4 mat = { X, Y, Z, W };
    return mat;
}

matrix_float4x4 matrix_float4x4_uniform_scale(float scale)
{
    vector_float4 X = { scale, 0, 0, 0 };
    vector_float4 Y = { 0, scale, 0, 0 };
    vector_float4 Z = { 0, 0, scale, 0 };
    vector_float4 W = { 0, 0, 0, 1 };

    matrix_float4x4 mat = { X, Y, Z, W };
    return mat;
}

matrix_float4x4 matrix_float4x4_rotation(vector_float3 axis, float angle)
{
    float c = cos(angle);
    float s = sin(angle);

    vector_float4 X;
    X.x = axis.x * axis.x + (1 - axis.x * axis.x) * c;
    X.y = axis.x * axis.y * (1 - c) - axis.z * s;
    X.z = axis.x * axis.z * (1 - c) + axis.y * s;
    X.w = 0.0;

    vector_float4 Y;
    Y.x = axis.x * axis.y * (1 - c) + axis.z * s;
    Y.y = axis.y * axis.y + (1 - axis.y * axis.y) * c;
    Y.z = axis.y * axis.z * (1 - c) - axis.x * s;
    Y.w = 0.0;

    vector_float4 Z;
    Z.x = axis.x * axis.z * (1 - c) - axis.y * s;
    Z.y = axis.y * axis.z * (1 - c) + axis.x * s;
    Z.z = axis.z * axis.z + (1 - axis.z * axis.z) * c;
    Z.w = 0.0;

    vector_float4 W;
    W.x = 0.0;
    W.y = 0.0;
    W.z = 0.0;
    W.w = 1.0;

    matrix_float4x4 mat = { X, Y, Z, W };
    return mat;
}

matrix_float4x4 matrix_float4x4_perspective(float aspect, float fovy, float near, float far)
{
    float yScale = 1 / tan(fovy * 0.5);
    float xScale = yScale / aspect;
    float zRange = far - near;
    float zScale = -(far + near) / zRange;
    float wzScale = -2 * far * near / zRange;

    vector_float4 P = { xScale, 0, 0, 0 };
    vector_float4 Q = { 0, yScale, 0, 0 };
    vector_float4 R = { 0, 0, zScale, -1 };
    vector_float4 S = { 0, 0, wzScale, 0 };

    matrix_float4x4 mat = { P, Q, R, S };
    return mat;
}

*/
