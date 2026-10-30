#version 410

in vec2 TexCoord;

uniform sampler2D texture1;

out vec4 color;

void main() {
    float r;

    // r = (gl_PrimitiveID % 255) / 255.0f;
    // color = vec4(r, r, r, 1);
    color = texture(texture1, TexCoord);
}
