#version 410 core

layout(location = 0) in vec2 pos;
layout(location = 1) in vec2 uv;

out vec2 pass_uv;

uniform mat3 transform;

void main() {
    gl_Position = vec4(transform * vec3(pos, 1.0), 1.0);
    pass_uv = uv;
}