out vec4 FragColor;

uniform vec4 ourColor;

in vec4 trashPos;

void main()
{
    if (isOutsideClipArea(trashPos)) {
        discard;
    }

    FragColor = ourColor;
}