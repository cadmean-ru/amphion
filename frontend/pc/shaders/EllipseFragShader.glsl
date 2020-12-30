#version 330 core

flat in vec3 fTopLeft;
flat in vec3 fBottomRight;
in vec4 fPosition;
in vec4 fFillColor;

out vec4 resultColor;

void main()
{
    float a = (fBottomRight.x - fTopLeft.x) / 2;
    float b = (fBottomRight.y - fTopLeft.y) / 2;
    float x = fPosition.x;
    float y = fPosition.y;
    float xc = fTopLeft.x + a;
    float yc = fTopLeft.y + b;
    float x2 = x - xc;
    float y2 = y - yc;
    float c = x2 / a;
    float d = y2 / b;
    float res = c * c + d * d;
    if (res <= 1)
        resultColor = fFillColor;
    else
        discard;
}