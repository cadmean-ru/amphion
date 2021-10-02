in vec2 fTexCoord;
in vec4 fPosition;

out vec4 color;

uniform sampler2D uTexture;
uniform vec4 uTextColor;

void main()
{
    if (isOutsideClipArea(fPosition)) {
        discard;
    }

    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(uTexture, fTexCoord).r);
    color = uTextColor * sampled;
}