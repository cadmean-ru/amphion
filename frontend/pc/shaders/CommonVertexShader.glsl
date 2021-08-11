#version 330

uniform mat4 uProjection;

vec4 applyProjection(vec4 v) {
    return uProjection * v;
}

vec4 applyProjection(vec3 v) {
    return applyProjection(vec4(v.xyz, 1));
}
