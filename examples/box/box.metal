using namespace metal;

struct Vertex {
    float4 position [[position]];
    float4 color; 
};

struct Uniforms {
    float4x4 modelViewMatrix;
    float4x4 projectionMatrix; 
};

vertex Vertex vertex_project(device Vertex *vertices [[buffer(0)]], constant Uniforms *uniforms [[buffer(1)]], uint vid [[vertex_id]])
{
    Vertex vertexOut;
    vertexOut.position = uniforms->projectionMatrix * uniforms->modelViewMatrix * vertices[vid].position;
    vertexOut.color = vertices[vid].color;
    return vertexOut;
}

fragment half4 fragment_flatcolor(Vertex vertexIn [[stage_in]]) 
{
    return half4(vertexIn.color);
}
