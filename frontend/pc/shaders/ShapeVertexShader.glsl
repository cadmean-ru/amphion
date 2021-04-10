#version 330
//layout (location = 0) in ivec3 vPosition;
//layout (location = 1) in ivec4 vFillColor;
//layout (location = 2) in ivec4 vStrokeColor;
//layout (location = 3) in int vStrokeWeight;
//layout (location = 4) in int vCornerRadius;
//
//out vec4 fPosition;
//out vec4 fFillColor;
//out vec4 fStrokeColor;
//out float fStrokeWeight;
//out float fCornerRadius;
//
//uniform mat4 uProjection;
//
//const mat4 colorNormalizer = mat4(
//    1/255, 0, 0, 0,
//    0, 1/255, 0, 0,
//    0, 0, 1/255, 0,
//    0, 0, 0, 1/255
//);
//
//void main()
//{
//    fPosition = uProjection * vec4(vPosition.xyz, 1.0);
//    gl_Position = fPosition;
////    fFillColor = colorNormalizer * vFillColor;
//    fFillColor = vec4(uProjection[0][0], uProjection[1][1], 1, 1);
////    fStrokeColor = colorNormalizer * vStrokeColor;
////    fStrokeWeight = (uProjection * vec4(vStrokeWeight)).x;
////    fCornerRadius = (uProjection * vec4(vCornerRadius)).x;
//}

layout (location = 0) in vec3 vPosition;
layout (location = 1) in vec3 vTlPosition;
layout (location = 2) in vec3 vBrPosition;
layout (location = 3) in vec4 vFillCollor;
layout (location = 4) in float vStrokeWeight;
layout (location = 5) in vec4 vStrokeColor;
layout (location = 6) in float vCornerRadius;

//layout (location = 2) in vec3 vTest;

out vec4 fPosition;
flat out vec3 fTlPosition;
flat out vec3 fBrPosition;
out vec4 fFillColor;
out float fStrokeWeight;
out vec4 fStrokeColor;
out float fCornerRadius;

uniform mat4 uProjection;

void main() {
//    vec4 test = uProjection * vec4(vTest.xyz, 1.0);
    fPosition = vec4(vPosition.xyz, 1.0);
    gl_Position = fPosition;
    fTlPosition = vTlPosition;
    fBrPosition = vBrPosition;
    fFillColor = vFillCollor / 255;
    fStrokeWeight = vStrokeWeight;
    fStrokeColor = vStrokeColor / 255;
    fCornerRadius = vCornerRadius;
}