#version 330

in vec4 fPosition;
flat in vec3 fTlPosition;
flat in vec3 fBrPosition;
in vec4 fFillColor;
in float fStrokeWeight;
in vec4 fStrokeColor;
in float fCornerRadius;

out vec4 resultColor;

//void applyClippingArea2d(vec2 pos2d);

void main() {
//    applyClippingArea2d(fPosition.xy);

    if (fTlPosition.x + fStrokeWeight >= fPosition.x || fTlPosition.y - fStrokeWeight <= fPosition.y ||
        fBrPosition.x - fStrokeWeight <= fPosition.x || fBrPosition.y + fStrokeWeight >= fPosition.y) {

        resultColor = fStrokeColor;
    } else {
        resultColor = fFillColor;
    }
}