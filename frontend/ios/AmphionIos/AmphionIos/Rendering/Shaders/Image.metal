//
//  Image.metal
//  AmphionIos
//
//  Created by Алексей Крицков on 24.08.2021.
//

#include <metal_stdlib>
using namespace metal;


struct ImageIn {
    float4 position [[attribute(0)]];
    float2 texCoord [[attribute(1)]];
};

struct ImageOut {
    float4 position [[position]];
    float2 texCoord;
};

vertex ImageOut image_vertex(const ImageIn ImageIn [[stage_in]]) {
    ImageOut ImageOut;
    ImageOut.position = ImageIn.position;
    ImageOut.texCoord = ImageIn.texCoord;
    
    return ImageOut;
}

fragment float4 image_fragment(ImageOut imageIn [[stage_in]], texture2d<float> texture [[texture(0)]], sampler mySampler [[sampler(0)]]) {
    return texture.sample(mySampler, imageIn.texCoord);
}
