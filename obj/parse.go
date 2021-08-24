package obj

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"test/metal"
)

type Model struct {
	Vertices      []metal.Vector_float4
	TexCoords     []metal.Vector_float2
	Normals       []metal.Vector_float4
	Indices       []metal.Uint16
	TexIndices    []metal.Uint16
	NormalIndices []metal.Uint16
}

func Parse(filename string) (*Model, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := Model{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")
		switch tokens[0] {
		case "v":
			x, _ := strconv.ParseFloat(tokens[1], 32)
			y, _ := strconv.ParseFloat(tokens[2], 32)
			z, _ := strconv.ParseFloat(tokens[3], 32)
			m.Vertices = append(m.Vertices, metal.Vector4(x, y, z, 1))
		case "vt":
			x, _ := strconv.ParseFloat(tokens[1], 32)
			y, _ := strconv.ParseFloat(tokens[2], 32)
			m.TexCoords = append(m.TexCoords, metal.Vector2(x, y))
		case "vn":
			x, _ := strconv.ParseFloat(tokens[1], 32)
			y, _ := strconv.ParseFloat(tokens[2], 32)
			z, _ := strconv.ParseFloat(tokens[3], 32)
			m.Normals = append(m.Normals, metal.Vector4(x, y, z, 0))
		case "f":
			for i := 1; i < len(tokens); i++ {
				v := strings.Split(tokens[i], "/")
				switch len(v) {
				case 1:
					vi, _ := strconv.Atoi(v[0])
					m.Indices = append(m.Indices, metal.Uint16(vi))
				case 2:
					vi, _ := strconv.Atoi(v[0])
					vt, _ := strconv.Atoi(v[1])
					m.Indices = append(m.Indices, metal.Uint16(vi))
					m.TexIndices = append(m.Indices, metal.Uint16(vt))
				case 3:
					vi, _ := strconv.Atoi(v[0])
					m.Indices = append(m.Indices, metal.Uint16(vi-1))
					if v[1] != "" {
						vt, _ := strconv.Atoi(v[1])
						m.TexIndices = append(m.Indices, metal.Uint16(vt))
					}
					vn, _ := strconv.Atoi(v[2])
					m.NormalIndices = append(m.NormalIndices, metal.Uint16(vn-1))
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &m, nil
}
