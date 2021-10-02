in vec4 fPosition;
flat in vec3 fTlPosition;
flat in vec3 fBrPosition;
in vec4 fFillColor;
in float fStrokeWeight;
in vec4 fStrokeColor;
in float fCornerRadius;

out vec4 resultColor;

void main() {
    if (isOutsideClipArea(fPosition)) {
        discard;
    }

    resultColor = fFillColor;
}
