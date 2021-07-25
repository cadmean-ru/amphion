layout (location = 0) in vec3 vPosition;
layout (location = 1) in vec3 vTlPosition;
layout (location = 2) in vec3 vBrPosition;
layout (location = 3) in vec4 vFillCollor;
layout (location = 4) in float vStrokeWeight;
layout (location = 5) in vec4 vStrokeColor;
layout (location = 6) in float vCornerRadius;


out vec4 fPosition;
flat out vec3 fTlPosition;
flat out vec3 fBrPosition;
out vec4 fFillColor;
out float fStrokeWeight;
out vec4 fStrokeColor;
out float fCornerRadius;

void main() {
    fPosition = vec4(vPosition.xyz, 1.0);
    gl_Position = fPosition;
    fTlPosition = vTlPosition;
    fBrPosition = vBrPosition;
    fFillColor = vFillCollor / 255;
    fStrokeWeight = vStrokeWeight;
    fStrokeColor = vStrokeColor / 255;
    fCornerRadius = vCornerRadius;
}