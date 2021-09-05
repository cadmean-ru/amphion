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

struct Uniform {
    float4x4 projection;
};

vertex RectOut rect_vertex(const RectIn rectIn [[stage_in]], constant Uniform& uniform [[buffer(1)]]) {
    RectOut rectOut;
    rectOut.position = uniform.projection * rectIn.position;
    rectOut.color = rectIn.color;
    
    return rectOut;
}

fragment float4 rect_fragment(RectOut rectIn [[stage_in]]) {
    return rectIn.color;
}
