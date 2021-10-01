//
//  Text.metal
//  AmphionIos
//
//  Created by Алексей Крицков on 30.09.2021.
//

#include "Common.metal"

struct TextIn {
    float3 position [[attribute(0)]];
    float2 texCoord [[attribute(1)]];
    float4 color [[attribute(2)]];
};

struct TextOut {
    float4 position [[position]];
    float2 texCoord;
    float4 color;
};

vertex TextOut text_vertex(const TextIn textIn [[stage_in]], constant Uniform& uniform [[buffer(1)]]) {
    TextOut textOut;
    textOut.position = uniform.projection * float4(textIn.position, 1);
    textOut.texCoord = textIn.texCoord;
    textOut.color = textIn.color;
    
    return textOut;
}

fragment float4 text_fragment(TextOut textIn [[stage_in]], texture2d<uint> texture [[texture(0)]], sampler mySampler [[sampler(0)]]) {
    uint4 sampled = texture.sample(mySampler, textIn.texCoord);
//    if (sampled.x < 0.5) {
//        discard_fragment();
//    }
    
    return float4(textIn.color.xyz, float(sampled.x)/255*textIn.color.w);
}

