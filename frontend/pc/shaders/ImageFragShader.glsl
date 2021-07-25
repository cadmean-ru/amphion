out vec4 fragColor;

in vec2 fTexCoord;
in vec4 fPosition;

uniform sampler2D uTexture;

void main()
{
    if (isOutsideClipArea(fPosition)) {
        discard;
    }

    fragColor = texture(uTexture, fTexCoord);
}