in vec4 fPosition;
flat in vec3 fTlPosition;
flat in vec3 fBrPosition;
in vec4 fFillColor;
in float fStrokeWeight;
in vec4 fStrokeColor;
in float fCornerRadius;

out vec4 resultColor;

bool isInsideEllipse(float a, float b, float offset) {
    float x = fPosition.x;
    float y = fPosition.y;
    float xc = fTlPosition.x + offset + a;
    float yc = fTlPosition.y - offset + b;
    float x2 = x - xc;
    float y2 = y - yc;
    float c = x2 / a;
    float d = y2 / b;
    float res = c * c + d * d;
    return res <= 1;
}

void main()
{
    if (isOutsideClipArea(fPosition)) {
        discard;
    }

    float a = (fBrPosition.x - fTlPosition.x) / 2;
    float b = (fBrPosition.y - fTlPosition.y) / 2;
    float a1 = (fBrPosition.x - fTlPosition.x - fStrokeWeight - fStrokeWeight) / 2;
    float b1 = (fBrPosition.y - fTlPosition.y + fStrokeWeight + fStrokeWeight) / 2;

    if (isInsideEllipse(a, b, 0)) {
        if (fStrokeWeight > 0) {
            if (isInsideEllipse(a1, b1, fStrokeWeight))
                resultColor = fFillColor;
            else
                resultColor = fStrokeColor;
        } else {
            resultColor = fFillColor;
        }
    } else discard;
}