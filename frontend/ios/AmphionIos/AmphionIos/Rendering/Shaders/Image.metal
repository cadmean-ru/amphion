//
//  Image.metal
//  AmphionIos
//
//  Created by Алексей Крицков on 24.08.2021.
//

#include "Common.metal"

struct ImageIn {
    float4 position [[attribute(0)]];
    float2 texCoord [[attribute(1)]];
};

struct ImageOut {
    float4 position [[position]];
    float2 texCoord;
};

vertex ImageOut image_vertex(const ImageIn imageIn [[stage_in]], constant Uniform& uniform [[buffer(1)]]) {
    ImageOut imageOut;
    imageOut.position = uniform.projection * imageIn.position;
    imageOut.texCoord = imageIn.texCoord;
    
    return imageOut;
}

fragment float4 image_fragment(ImageOut imageIn [[stage_in]], texture2d<uint> texture [[texture(0)]], sampler mySampler [[sampler(0)]]) {
    float4 sampled = float4(texture.sample(mySampler, imageIn.texCoord));
    return float4(sampled.x/255, sampled.y/255, sampled.z/255, sampled.w/255);
}
