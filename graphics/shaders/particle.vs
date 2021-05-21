#version 410 core

layout(location = 0) in vec2 pos;
layout(location = 1) in vec2 uv;
layout(location = 2) in mat3 transform;
layout(location = 5) in float frame;

out float pass_frame;
out vec2 pass_uv;

uniform mat3 globalTransform;

void main() {
    gl_Position = vec4(globalTransform * transform * vec3(pos, 1.0), 1.0);
    pass_uv = uv;
    pass_frame = frame;
}