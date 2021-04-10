#version 330
//
//in vec4 fPosition;
//in vec4 fFillColor;
//in vec4 fStrokeColor;
//in float fStrokeWeight;
//in float fCornerRadius;
//
//out vec4 resultColor;
//
//void main()
//{
//    resultColor = fFillColor;
//}

in vec4 fPosition;
flat in vec3 fTlPosition;
flat in vec3 fBrPosition;
in vec4 fFillColor;
in float fStrokeWeight;
in vec4 fStrokeColor;
in float fCornerRadius;

out vec4 resultColor;

void main() {
    if (fTlPosition.x + fStrokeWeight >= fPosition.x || fTlPosition.y - fStrokeWeight <= fPosition.y ||
        fBrPosition.x - fStrokeWeight <= fPosition.x || fBrPosition.y + fStrokeWeight >= fPosition.y) {

        resultColor = fStrokeColor;
    } else {
        resultColor = fFillColor;
    }
}