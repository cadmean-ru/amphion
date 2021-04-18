#version 330

uniform vec4 uClippingArea2d;

void applyClippingArea2d(vec2 pos2d) {
    float xMin = uClippingArea2d.x;
    float yMin = uClippingArea2d.y;
    float xMax = uClippingArea2d.z;
    float yMax = uClippingArea2d.w;

    if (pos2d.x < xMin || pos2d.x > xMax || pos2d.y < yMin || pos2d.y > yMax) {
        discard;
    }
}