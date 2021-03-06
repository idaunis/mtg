#include <metal_stdlib>
#include <metal_matrix>

using namespace metal;

struct Light
{
    float3 direction;
    float3 ambientColor;
    float3 diffuseColor;
    float3 specularColor;
};

constant Light light = {
    .direction = { 0.13, 0.72, 0.68 },
    .ambientColor = { 0.05, 0.05, 0.05 },
    .diffuseColor = { 0.9, 0.9, 0.9 },
    .specularColor = { 1, 1, 1 }
};

struct Material
{
    float3 ambientColor;
    float3 diffuseColor;
    float3 specularColor;
    float specularPower;
};

constant Material material = {
    .ambientColor = { 0.9, 0.1, 0 },
    .diffuseColor = { 0.9, 0.1, 0 },
    .specularColor = { 1, 1, 1 },
    .specularPower = 100
};

struct Vertex {
    float4 position;
    float4 normal;
};

struct ProjectedVertex
{
    float4 position [[position]];
    float3 eye;
    float3 normal;
};

struct Uniforms {
    float4x4 modelViewProjectionMatrix;
    float4x4 modelViewMatrix;
    float3x3 normalMatrix;
};

vertex ProjectedVertex vertex_project(const device Vertex *vertices [[buffer(0)]],
                                        constant Uniforms &uniforms [[buffer(1)]],
                                        uint vid [[vertex_id]])
{
    ProjectedVertex vertexOut;
    vertexOut.position = uniforms.modelViewProjectionMatrix * vertices[vid].position;
    vertexOut.eye =  -(uniforms.modelViewMatrix * vertices[vid].position).xyz;
    vertexOut.normal = uniforms.normalMatrix * vertices[vid].normal.xyz;

    return vertexOut;
}

fragment float4 fragment_light(ProjectedVertex vert [[stage_in]],
                        constant Uniforms &uniforms [[buffer(1)]])
{
    float3 ambientTerm = light.ambientColor * material.ambientColor;

    float3 normal = normalize(vert.normal);
    float diffuseIntensity = saturate(dot(normal, light.direction));
    float3 diffuseTerm = light.diffuseColor * material.diffuseColor * diffuseIntensity;

    float3 specularTerm(0);
    if (diffuseIntensity > 0)
    {
        float3 eyeDirection = normalize(vert.eye);
        float3 halfway = normalize(light.direction + eyeDirection);
        float specularFactor = pow(saturate(dot(normal, halfway)), material.specularPower);
        specularTerm = light.specularColor * material.specularColor * specularFactor;
    }

    return float4(ambientTerm + diffuseTerm + specularTerm, 1);
}
