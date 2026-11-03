#version 410

in vec2 v_text_crd;

uniform uint u_mode;
uniform sampler2D u_texture1;
uniform float u_col_ratio;

out vec4 color;

void main() {
  float r = (gl_PrimitiveID % 255) / 255.0;
  vec4 base_color = vec4(r, r, r, 1);

  if (u_mode == 0) {
    color = base_color;
  } else if (u_mode == 1) {
    vec2 inv_text_crd = vec2(v_text_crd.x, 1.0 - v_text_crd.y);
    vec4 text_color = texture(u_texture1, inv_text_crd);

    color = mix(base_color, text_color, u_col_ratio);
  }
}
