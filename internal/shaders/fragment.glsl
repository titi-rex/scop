#version 410

in vec2 v_text_crd;

uniform int u_mode;
uniform sampler2D u_texture1;

out vec4 color;

void main() {
  vec2 text_crd_fixed = vec2(v_text_crd.x, 1.0 - v_text_crd.y);
  color = texture(u_texture1, text_crd_fixed);
  // color = texture(u_texture1, v_text_crd);

}


// void main() {
//   if (u_mode == 0) {
//     float r = (gl_PrimitiveID % 255) / 255.0;
//     color = vec4(r, r, r, 1);
//   } else if (u_mode == 1) {
//     // vec2 text_crd_fixed = vec2(v_text_crd.x, 1.0 - v_text_crd.y);
//     // color = texture(u_texture1, text_crd_fixed);
//     color = texture(u_texture1, v_text_crd);
//   }
// }
