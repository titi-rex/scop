#version 410

layout (location = 0) in vec3 i_position;
layout (location = 1) in vec2 i_text_crd;

uniform mat4 u_model;
uniform mat4 u_view;
uniform mat4 u_projection;

out vec2 v_text_crd;

void main() {
    gl_Position = u_projection * u_view * u_model * vec4(i_position, 1.0);
    v_text_crd = i_text_crd;
    // v_text_crd = i_position.yz;
}
