#version 410
in vec3 vp;

uniform mat4 transform;

void main() {
    gl_Position = transform * vec4(vp, 1.0);
}
