#version 330

uniform mat4 uProjection;

vec4 applyProjection(vec4 v) {
    return uProjection * v;
}

vec4 applyProjection(vec3 v) {
    return applyProjection(vec4(v.xyz, 1));
}

float applyProjectionScalar(float f) {
    vec4 res = applyProjection(vec4(f, f, f, 1));
    res += vec4(1, 1, 1, 0);
    return res.x;
}
