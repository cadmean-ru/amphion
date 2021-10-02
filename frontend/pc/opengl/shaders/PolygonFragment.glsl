in vec4 fPosition;
in vec4 fFillColor;

out vec4 fragmentColor;

void main() {
    if (isOutsideClipArea(fPosition)) {
        discard;
    }

    fragmentColor = fFillColor;
}
