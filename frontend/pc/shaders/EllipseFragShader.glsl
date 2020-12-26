#version 330 core
out vec4 FragColor;

uniform vec4 ourColor;

uniform vec3 tlPos;
uniform vec3 brPos;

in vec4 trashPos;

void main()
{
    float a = (brPos.x - tlPos.x) / 2; //0.1
    float b = (brPos.y - tlPos.y) / 2; //0.1
    float x = trashPos.x;
    float y = trashPos.y;
    float xc = tlPos.x + a; // 0.3
    float yc = tlPos.y + b; // -0.3
    float x2 = x - xc;
    float y2 = y - yc;
    float c = x2 / a;
    float d = y2 / b;
    float res = c * c + d * d;
    if (res <= 1)
    FragColor = ourColor;
    else
    discard;
}