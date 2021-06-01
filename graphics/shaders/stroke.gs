#version 440 core

layout (lines) in;
layout (triangle_strip, max_vertices = 4) out;

in vec2 pass_normal[];

void main() {
    float width = 0.05;

    gl_Position = gl_in[0].gl_Position + vec4(pass_normal[0] * width, 0.0, 1.0); 
    EmitVertex();

    gl_Position = gl_in[0].gl_Position + vec4(pass_normal[0] * -width, 0.0, 1.0);
    EmitVertex();

    gl_Position = gl_in[1].gl_Position + vec4(pass_normal[1] * width, 0.0, 1.0);
    EmitVertex();

    gl_Position = gl_in[1].gl_Position + vec4(pass_normal[1] * -width, 0.0, 1.0);
    EmitVertex();
    
    EndPrimitive();
}