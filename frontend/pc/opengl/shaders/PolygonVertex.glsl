layout (location = 0) in vec3 vPosition;
layout (location = 1) in vec4 vFillCollor;

out vec4 fPosition;
out vec4 fFillColor;

void main() {
    fPosition = applyProjection(vPosition);
    gl_Position = fPosition;
    fFillColor = vFillCollor / 255;
}
