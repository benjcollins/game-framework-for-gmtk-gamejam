#version 410 core

layout (lines) in;
layout (triangle_strip, max_vertices = 3) out;
// layout (line_strip, max_vertices = 3) out;

in vec2 pass_normal[];

uniform vec2 origin;

void main() {
    gl_Position = vec4(origin, 0.0, 1.0);
    EmitVertex();
    gl_Position = gl_in[0].gl_Position;
    EmitVertex();
    gl_Position = gl_in[1].gl_Position;
    EmitVertex();
    EndPrimitive();
}