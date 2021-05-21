#version 410 core

in vec2 pass_uv;
in float pass_frame;

out vec4 frag_color;

uniform sampler2D textureSampler;
uniform int hFrames;
uniform int vFrames;

vec4 sampleFrame(int frame) {
    float u = pass_uv.x / float(hFrames) + (frame % 4) / float(hFrames);
    float v = pass_uv.y / float(vFrames) + (frame / 4) / float(vFrames);
    return texture(textureSampler, vec2(u, v));
}

void main() {
    int frame = int(pass_frame);
    vec4 prevFrame = sampleFrame(frame);
    vec4 nextFrame = sampleFrame(frame+1);

    frag_color = mix(prevFrame, nextFrame, fract(pass_frame));
}