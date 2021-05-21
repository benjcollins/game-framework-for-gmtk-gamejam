#version 410 core

in vec2 pass_uv;

out vec4 frag_color;

uniform sampler2D textureSampler;

void main() {
    frag_color = texture(textureSampler, pass_uv);
}