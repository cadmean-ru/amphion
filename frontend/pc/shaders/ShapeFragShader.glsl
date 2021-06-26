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

    float cr = min(min(size.x/2, size.y/2), fCornerRadius);
    float cr2 = cr - fStrokeWeight;

    if (isInsideRect(pos2d, tl, tl + vec2(cr, -cr))) {
        if (!isInsideCircle(pos2d, tl + vec2(cr, -cr), cr)) {
            discard;
        } else if (!isInsideCircle(pos2d, tl + vec2(cr, -cr), cr2)) {
            resultColor = fStrokeColor;
            return;
        }
    }
    if (isInsideRect(pos2d, tr - vec2(cr, 0), tr - vec2(0, cr))) {
        if (!isInsideCircle(pos2d, tr - vec2(cr), cr)) {
            discard;
        } else if (!isInsideCircle(pos2d, tr - vec2(cr), cr2)) {
            resultColor = fStrokeColor;
            return;
        }
    }
    if (isInsideRect(pos2d, bl + vec2(0, cr), bl + vec2(cr, 0))) {
        if (!isInsideCircle(pos2d, bl + vec2(cr), cr)) {
            discard;
        } else if (!isInsideCircle(pos2d, bl + vec2(cr), cr2)) {
            resultColor = fStrokeColor;
            return;
        }
    }
    if (isInsideRect(pos2d, br - vec2(cr, -cr), br)) {
        if (!isInsideCircle(pos2d, br - vec2(cr, -cr), cr)) {
            discard;
        } else if (!isInsideCircle(pos2d, br - vec2(cr, -cr), cr2)) {
            resultColor = fStrokeColor;
            return;
        }
    }

    if (isInsideRect(pos2d, tl + vec2(fStrokeWeight, -fStrokeWeight), br - vec2(fStrokeWeight, -fStrokeWeight))) {
        resultColor = fFillColor;
    } else {
        resultColor = fStrokeColor;
    }
}