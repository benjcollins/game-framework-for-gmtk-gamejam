#version 410 core

layout (lines) in;
layout (triangle_strip, max_vertices = 8) out;
// layout (line_strip, max_vertices = 8) out;

in vec2 pass_normal[];

uniform mat3 transform;
uniform float width;

out float invWidth;

void calculatePosition(int index, float w) {
    gl_Position = gl_in[index].gl_Position + vec4(transform * vec3(pass_normal[index] * w, 0.0), 0.0);
    invWidth = width - abs(w);
    EmitVertex();
}

void main() {
    calculatePosition(0, width);
    calculatePosition(0, 0);
    calculatePosition(1, width);
    calculatePosition(1, 0);
    
    EndPrimitive();

    calculatePosition(0, 0);
    calculatePosition(0, -width);
    calculatePosition(1, 0);
    calculatePosition(1, -width);

    EndPrimitive();
}