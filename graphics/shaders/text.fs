#version 410 core

in vec2 pass_uv;

out vec4 frag_color;

uniform sampler2D textureSampler;

void main() {
    float width = 0.5;
    float alias = 0.08;

    vec4 color = texture(textureSampler, pass_uv);
    if (color.r < width - alias) {

    } else if (color.r < width) {
        frag_color = vec4(0.0, 0.0, 0.0, (color.r - width + alias) / alias);
    } else {
        frag_color = vec4(0.0, 0.0, 0.0, 1.0);
    }
}