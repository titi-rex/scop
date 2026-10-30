#version 410

layout (location = 0) in vec3 vp;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform mat4 texture_matrix;

out vec2 TexCoord;

void main() {
    vec4 modelPosition = model * vec4(vp, 1.0);

    gl_Position = projection * view * modelPosition;

    TexCoord = modelPosition.xz;
    // TexCoord = texture_matrix * modelPosition;
}
