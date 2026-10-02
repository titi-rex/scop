#version 410

layout (location = 0) in vec3 vp;
// layout (location = 1) in vec3 vc;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec3 vertexColor;

void main() {
    gl_Position = projection * view * model * vec4(vp, 1.0);
    vertexColor = vec3(1,0,0);
}
