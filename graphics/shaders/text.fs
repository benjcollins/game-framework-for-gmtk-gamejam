#version 410 core

in vec2 pass_uv;

out vec4 frag_color;

uniform sampler2D textureSampler;

uniform mat3 transform;

void main() {
    float width = 0.5;
    float thing = 0.0003;
    float alias = (inverse(transform) * vec3(thing, thing, 0.0)).x;

    vec4 color = texture(textureSampler, pass_uv);
    float brightness = color.r;

    if (brightness < width - alias) {

    } else if (brightness < width) {
        frag_color = vec4(0.0, 0.0, 0.0, (brightness - width + alias) / alias);
    } else {
        frag_color = vec4(0.0, 0.0, 0.0, 1.0);
    }
}