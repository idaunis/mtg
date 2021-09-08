package obj

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/idaunis/mtg/metal"
)

type Model struct {
	vertices  []metal.Vector_float4
	texCoords []metal.Vector_float2
	normals   []metal.Vector_float4
	groups    []*Group

	currentGroup *Group
}

type Group struct {
	Name     string
	Indices  []metal.Uint16
	Vertices []PackedVertex

	faceVertexToIndexMap map[FaceVertex]metal.Uint16
}

func (g *Group) GenerateNormals() {
	for i := 0; i < len(g.Vertices); i++ {
		g.Vertices[i].Normal = metal.Vector4(0, 0, 0, 0)
	}

	for i := 0; i < len(g.Indices); i += 3 {
		i0 := int(g.Indices[i])
		i1 := int(g.Indices[i+1])
		i2 := int(g.Indices[i+2])

		p0 := g.Vertices[i0].Position.XYZ()
		p1 := g.Vertices[i1].Position.XYZ()
		p2 := g.Vertices[i2].Position.XYZ()

		cross := metal.Vector3_cross(p1.Diff(p0), p2.Diff(p0))
		cross4 := metal.Vector4(float64(cross[0]), float64(cross[1]), float64(cross[2]), 0)

		g.Vertices[i0].Normal = g.Vertices[i0].Normal.Add(cross4)
		g.Vertices[i1].Normal = g.Vertices[i1].Normal.Add(cross4)
		g.Vertices[i2].Normal = g.Vertices[i2].Normal.Add(cross4)
	}

	for i := 0; i < len(g.Vertices); i++ {
		g.Vertices[i].Normal = metal.Vector4_normalize(g.Vertices[i].Normal)
	}
}

type PackedVertex struct {
	Position metal.Vector_float4
	TexCoord metal.Vector_float2
	Normal   metal.Vector_float4
}

func (g *Group) EachVertex(fn func(p PackedVertex)) {
	for i := range g.Vertices {
		fn(g.Vertices[i])
	}
}

type FaceVertex struct {
	vertexIndex  metal.Uint16
	normalIndex  metal.Uint16
	textureIndex metal.Uint16
}

func ternary(condition bool, a, b int) metal.Uint16 {
	if condition {
		return metal.Uint16(a)
	}
	return metal.Uint16(b)
}

func (m *Model) AddIndexFromFaceIndex(fv FaceVertex) {
	g := m.currentGroup

	if g.faceVertexToIndexMap == nil {
		g.faceVertexToIndexMap = make(map[FaceVertex]metal.Uint16)
	}

	if gIndex, found := g.faceVertexToIndexMap[fv]; found {
		g.Indices = append(g.Indices, gIndex)
	} else {
		position := m.vertices[fv.vertexIndex]
		normal := metal.Vector_float4{0, 1, 0, 0}
		texture := metal.Vector_float2{0, 0}
		if fv.normalIndex != 0xFFFF {
			normal = m.normals[fv.normalIndex]
		}
		if fv.textureIndex != 0xFFFF {
			texture = m.texCoords[fv.textureIndex]
		}
		g.Vertices = append(g.Vertices, PackedVertex{
			Position: position,
			TexCoord: texture,
			Normal:   normal,
		})
		gIndex = metal.Uint16(len(g.Vertices) - 1)
		g.Indices = append(g.Indices, gIndex)

		g.faceVertexToIndexMap[fv] = gIndex
	}
}

func (m *Model) GetGroup(index int) *Group {
	if index < 0 || index >= len(m.groups) {
		// err := fmt.Errorf("invalid group index, there are %d groups", len(m.groups))
		return m.currentGroup
	}
	return m.groups[index]
}

func (m *Model) NewGroup(name string) {
	group := &Group{Name: name}
	m.currentGroup = group
	m.groups = append(m.groups, group)
}

func Parse(filename string) (*Model, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := Model{}
	m.NewGroup("(unnamed)")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")
		switch tokens[0] {
		case "g":
			name := tokens[1]
			m.NewGroup(name)
		case "v":
			x, _ := strconv.ParseFloat(tokens[1], 32)
			y, _ := strconv.ParseFloat(tokens[2], 32)
			z, _ := strconv.ParseFloat(tokens[3], 32)
			m.vertices = append(m.vertices, metal.Vector4(x, y, z, 1))
		case "vt":
			x, _ := strconv.ParseFloat(tokens[1], 32)
			y, _ := strconv.ParseFloat(tokens[2], 32)
			m.texCoords = append(m.texCoords, metal.Vector2(x, y))
		case "vn":
			x, _ := strconv.ParseFloat(tokens[1], 32)
			y, _ := strconv.ParseFloat(tokens[2], 32)
			z, _ := strconv.ParseFloat(tokens[3], 32)
			m.normals = append(m.normals, metal.Vector4(x, y, z, 0))
		case "f":
			faces := []FaceVertex{}
			for i := 1; i < len(tokens); i++ {
				var vi, ti, ni int
				v := strings.Split(tokens[i], "/")
				switch len(v) {
				case 1:
					vi, _ = strconv.Atoi(v[0])
				case 2:
					vi, _ = strconv.Atoi(v[0])
					ti, _ = strconv.Atoi(v[1])
				case 3:
					vi, _ = strconv.Atoi(v[0])
					ti, _ = strconv.Atoi(v[1])
					ni, _ = strconv.Atoi(v[2])
				}

				faces = append(faces, FaceVertex{
					vertexIndex:  ternary(vi < 0, len(m.vertices)+vi-1, vi-1),
					textureIndex: ternary(ti < 0, len(m.texCoords)+ti-1, ti-1),
					normalIndex:  ternary(ni < 0, len(m.vertices)+ni-1, ni-1),
				})
			}

			for i := 0; i < len(faces)-2; i++ {
				m.AddIndexFromFaceIndex(faces[i])
				m.AddIndexFromFaceIndex(faces[i+1])
				m.AddIndexFromFaceIndex(faces[i+2])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &m, nil
}
