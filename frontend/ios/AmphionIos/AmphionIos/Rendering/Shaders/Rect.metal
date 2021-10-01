//
//  Rect.metal
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

#include "Common.metal"

struct RectIn {
    float3 position [[attribute(0)]];
    float4 color [[attribute(1)]];
};

struct RectOut {
    float4 position [[position]];
    float4 color;
};

vertex RectOut rect_vertex(const RectIn rectIn [[stage_in]], constant Uniform& uniform [[buffer(1)]]) {
    RectOut rectOut;
    rectOut.position = uniform.projection * float4(rectIn.position, 1);
    rectOut.color = rectIn.color / 255;
    
    return rectOut;
}

fragment float4 rect_fragment(RectOut rectIn [[stage_in]]) {
    return rectIn.color;
}
