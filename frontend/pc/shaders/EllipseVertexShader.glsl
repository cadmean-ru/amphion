#version 330 core

layout (location = 0) in vec3 vPosition;
layout (location = 1) in vec4 vFillCollor;
layout (location = 2) in vec3 vTopLeft;
layout (location = 3) in vec3 vBottomRight;

out vec4 fPosition;
out vec4 fFillColor;
flat out vec3 fTopLeft;
flat out vec3 fBottomRight;

uniform mat4 uProjection;

const float c1 = float(1)/float(255);

const mat4 colorNormalizer = mat4(
    c1, 0, 0, 0,
    0, c1, 0, 0,
    0, 0, c1, 0,
    0, 0, 0, c1
);

void main() {
    fPosition = vec4(vPosition.xyz, 1.0);
    gl_Position = fPosition;
    fFillColor = colorNormalizer * vFillCollor;
    fTopLeft = vTopLeft;
    fBottomRight = vBottomRight;
}