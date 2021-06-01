#version 410 core

layout(location = 0) in vec2 pos;
layout(location = 1) in vec2 normal;

out vec2 pass_normal;

void main() {
    gl_Position = vec4(pos, 0.0, 1.0);
    pass_normal = normal;
}