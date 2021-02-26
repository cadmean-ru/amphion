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
in vec4 fFillColor;

out vec4 resultColor;

void main() {
    resultColor = fFillColor;
}