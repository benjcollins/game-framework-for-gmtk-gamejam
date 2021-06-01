#version 410 core

in float invWidth;

out vec4 frag_color;

uniform vec4 color;

void main() {
    float threshold = 0.004;

    if (invWidth > threshold) {
        frag_color = color;
    } else {
        frag_color = invWidth / threshold * color;
    }
}