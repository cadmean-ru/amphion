#version 330

uniform vec4 uClipArea2dRect;
uniform int uClipArea2dShape;

bool isInsideRect(vec2 pos, vec2 tl, vec2 br) {
    return pos.x >= tl.x && pos.x <= br.x && pos.y <= tl.y && pos.y >= br.y;
}

bool isInsideCircle(vec2 pos, vec2 tl, vec2 br) {
    float r = (br.x - tl.x) / 2;
    vec2 center = vec2(tl.x + r, tl.y - r);

    float a = pos.x - center.x;
    float b = pos.y - center.y;
    return sqrt(a*a + b*b) <= r;
}

bool isInsideCircle(vec2 pos, vec2 center, float r) {
    float a = pos.x - center.x;
    float b = pos.y - center.y;
    return sqrt(a*a + b*b) <= r;
}

bool isInsideEllipse(vec2 pos, vec2 tl, vec2 br) {
    float x = pos.x;
    float y = pos.y;
    float a = (br.x - tl.x) / 2;
    float b = (tl.y - br.y) / 2;
    float xc = tl.x + a;
    float yc = br.y + b;
    float x2 = x - xc;
    float y2 = y - yc;
    float c = x2 / a;
    float d = y2 / b;
    float res = c * c + d * d;
    return res <= 1;
}

float _clipCornerRadius;

bool isInsideRoundedRect(vec2 pos2d, vec2 tl, vec2 br) {
    vec2 size = abs(br - tl);

    if (_clipCornerRadius > 0) {
        vec2 bl = vec2(tl.x, br.y);
        vec2 tr = vec2(br.x, tl.y);

        float cr = min(min(size.x/2, _clipCornerRadius), min(size.y/2, _clipCornerRadius));

        if (isInsideRect(pos2d, tl, tl + vec2(cr, -cr))) {
            if (!isInsideCircle(pos2d, tl + vec2(cr, -cr), cr)) {
                return false;
            }
        }
        if (isInsideRect(pos2d, tr - vec2(cr, 0), tr - vec2(0, cr))) {
            if (!isInsideCircle(pos2d, tr - vec2(cr), cr)) {
                return false;
            }
        }
        if (isInsideRect(pos2d, bl + vec2(0, cr), bl + vec2(cr, 0))) {
            if (!isInsideCircle(pos2d, bl + vec2(cr), cr)) {
                return false;
            }
        }
        if (isInsideRect(pos2d, br - vec2(cr, -cr), br)) {
            if (!isInsideCircle(pos2d, br - vec2(cr, -cr), cr)) {
                return false;
            }
        }
    }

    return isInsideRect(pos2d, tl, br);
}

bool isOutsideClipArea(vec2 fragPos2d) {
    vec2 tl = vec2(uClipArea2dRect.x, uClipArea2dRect.y);
    vec2 br = vec2(uClipArea2dRect.z, uClipArea2dRect.w);

    switch (uClipArea2dShape) {
        case 0:
            return false;
        case 2:
            return !isInsideRoundedRect(fragPos2d, tl, br);
        case 3:
            return !isInsideCircle(fragPos2d, tl, br);
        case 4:
            return !isInsideEllipse(fragPos2d, tl, br);
        default:
            return !isInsideRect(fragPos2d, tl, br);
    }
}

bool isOutsideClipArea(vec4 fragPos2d) {
    return isOutsideClipArea(fragPos2d.xy);
}

