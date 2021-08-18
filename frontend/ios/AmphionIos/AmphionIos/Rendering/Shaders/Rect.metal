//
//  Rect.metal
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

#include <metal_stdlib>
using namespace metal;

struct RectIn {
    float4 position [[attribute(0)]];
    float4 color [[attribute(1)]];
};

struct RectOut {
    float4 position [[position]];
    float4 color;
};

vertex RectOut rect_vertex(const RectIn rectIn [[stage_in]]) {
    RectOut rectOut;
    rectOut.position = rectIn.position;
    rectOut.color = rectIn.color;
    
    return rectOut;
}

fragment float4 rect_fragment(RectOut rectIn [[stage_in]]) {
    return rectIn.color;
}
