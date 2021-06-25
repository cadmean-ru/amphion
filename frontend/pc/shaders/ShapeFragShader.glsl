#version 330

in vec4 fPosition;
flat in vec3 fTlPosition;
flat in vec3 fBrPosition;
in vec4 fFillColor;
in float fStrokeWeight;
in vec4 fStrokeColor;
in float fCornerRadius;

out vec4 resultColor;

bool isInsideRect(vec2 pos, vec2 tl, vec2 br) {
    return pos.x >= tl.x && pos.x <= br.x && pos.y <= tl.y && pos.y >= br.y;
}

bool isInsideCircle(vec2 pos, vec2 center, float r) {
    float a = pos.x - center.x;
    float b = pos.y - center.y;
    return sqrt(a*a + b*b) <= r;
}

void main() {
    vec2 pos2d = fPosition.xy;
    vec2 tl = fTlPosition.xy;
    vec2 br = fBrPosition.xy;
    vec2 size = abs(br - tl);

    vec2 bl = vec2(tl.x, br.y);
    vec2 tr = vec2(br.x, tl.y);

    if (isInsideRect(pos2d, tl, tl + vec2(fCornerRadius, -fCornerRadius))) {
        if (!isInsideCircle(pos2d, tl + vec2(fCornerRadius, -fCornerRadius), fCornerRadius)) {
            discard;
        }
    }
    if (isInsideRect(pos2d, tr - vec2(fCornerRadius, 0), tr - vec2(0, fCornerRadius))) {
        if (!isInsideCircle(pos2d, tr - vec2(fCornerRadius), fCornerRadius)) {
            discard;
        }
    }
    if (isInsideRect(pos2d, bl + vec2(0, fCornerRadius), bl + vec2(fCornerRadius, 0))) {
        if (!isInsideCircle(pos2d, bl + vec2(fCornerRadius), fCornerRadius)) {
            discard;
        }
    }
    if (isInsideRect(pos2d, br - vec2(fCornerRadius, -fCornerRadius), br)) {
        if (!isInsideCircle(pos2d, br - vec2(fCornerRadius, -fCornerRadius), fCornerRadius)) {
            discard;
        }
    }

    if (isInsideRect(pos2d, tl + vec2(fStrokeWeight), br - vec2(fStrokeWeight))) {
        resultColor = fFillColor;
    } else {
        resultColor = fStrokeColor;
    }
}