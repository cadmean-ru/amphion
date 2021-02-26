
// +build darwin linux windows

package pc

// This file was auto-generated

const DefaultFragShaderStr = "#version 330\nout vec4 FragColor;\n\nuniform vec4 ourColor;\n\nvoid main()\n{\n    FragColor = ourColor;\n}\x00"

const DefaultVertexShaderStr = "#version 330\nlayout (location = 0) in vec3 aPos;\n\nout vec4 trashPos;\n\nvoid main()\n{\n    trashPos = vec4(aPos.x, aPos.y, aPos.z, 1.0);\n    gl_Position = trashPos;\n}\x00"

const EllipseFragShaderStr = "#version 330\n\nflat in vec3 fTopLeft;\nflat in vec3 fBottomRight;\nin vec4 fPosition;\nin vec4 fFillColor;\n\nout vec4 resultColor;\n\nvoid main()\n{\n    float a = (fBottomRight.x - fTopLeft.x) / 2;\n    float b = (fBottomRight.y - fTopLeft.y) / 2;\n    float x = fPosition.x;\n    float y = fPosition.y;\n    float xc = fTopLeft.x + a;\n    float yc = fTopLeft.y + b;\n    float x2 = x - xc;\n    float y2 = y - yc;\n    float c = x2 / a;\n    float d = y2 / b;\n    float res = c * c + d * d;\n    if (res <= 1)\n        resultColor = fFillColor;\n    else\n        discard;\n}\x00"

const EllipseVertexShaderStr = "#version 330\n\nlayout (location = 0) in vec3 vPosition;\nlayout (location = 1) in vec4 vFillCollor;\nlayout (location = 2) in vec3 vTopLeft;\nlayout (location = 3) in vec3 vBottomRight;\n\nout vec4 fPosition;\nout vec4 fFillColor;\nflat out vec3 fTopLeft;\nflat out vec3 fBottomRight;\n\nuniform mat4 uProjection;\n\nconst float c1 = float(1)/float(255);\n\nconst mat4 colorNormalizer = mat4(\n    c1, 0, 0, 0,\n    0, c1, 0, 0,\n    0, 0, c1, 0,\n    0, 0, 0, c1\n);\n\nvoid main() {\n    fPosition = vec4(vPosition.xyz, 1.0);\n    gl_Position = fPosition;\n    fFillColor = colorNormalizer * vFillCollor;\n    fTopLeft = vTopLeft;\n    fBottomRight = vBottomRight;\n}\x00"

const ImageFragShaderStr = "#version 330\n\nout vec4 FragColor;\n\nin vec2 texCoord;\n\nuniform sampler2D ourTexture;\n\nvoid main()\n{\n    FragColor = texture(ourTexture, texCoord);\n}\x00"

const ImageVertexShaderStr = "#version 330\nlayout (location = 0) in vec3 aPos;\nlayout (location = 1) in vec2 aTexCoord;\n\nout vec2 texCoord;\n\nvoid main()\n{\n    gl_Position = vec4(aPos, 1.0);\n    texCoord = aTexCoord;\n}\x00"

const RectFragShaderStr = "\x00"

const ShapeFragShaderStr = "#version 330\n//\n//in vec4 fPosition;\n//in vec4 fFillColor;\n//in vec4 fStrokeColor;\n//in float fStrokeWeight;\n//in float fCornerRadius;\n//\n//out vec4 resultColor;\n//\n//void main()\n//{\n//    resultColor = fFillColor;\n//}\n\nin vec4 fPosition;\nin vec4 fFillColor;\n\nout vec4 resultColor;\n\nvoid main() {\n    resultColor = fFillColor;\n}\x00"

const ShapeVertexShaderStr = "#version 330\n//layout (location = 0) in ivec3 vPosition;\n//layout (location = 1) in ivec4 vFillColor;\n//layout (location = 2) in ivec4 vStrokeColor;\n//layout (location = 3) in int vStrokeWeight;\n//layout (location = 4) in int vCornerRadius;\n//\n//out vec4 fPosition;\n//out vec4 fFillColor;\n//out vec4 fStrokeColor;\n//out float fStrokeWeight;\n//out float fCornerRadius;\n//\n//uniform mat4 uProjection;\n//\n//const mat4 colorNormalizer = mat4(\n//    1/255, 0, 0, 0,\n//    0, 1/255, 0, 0,\n//    0, 0, 1/255, 0,\n//    0, 0, 0, 1/255\n//);\n//\n//void main()\n//{\n//    fPosition = uProjection * vec4(vPosition.xyz, 1.0);\n//    gl_Position = fPosition;\n////    fFillColor = colorNormalizer * vFillColor;\n//    fFillColor = vec4(uProjection[0][0], uProjection[1][1], 1, 1);\n////    fStrokeColor = colorNormalizer * vStrokeColor;\n////    fStrokeWeight = (uProjection * vec4(vStrokeWeight)).x;\n////    fCornerRadius = (uProjection * vec4(vCornerRadius)).x;\n//}\n\nlayout (location = 0) in vec3 vPosition;\nlayout (location = 1) in vec4 vFillCollor;\n//layout (location = 2) in vec3 vTest;\n\nout vec4 fPosition;\nout vec4 fFillColor;\n\nuniform mat4 uProjection;\n\nconst float c1 = float(1)/float(255);\n\nconst mat4 colorNormalizer = mat4(\n    c1, 0, 0, 0,\n    0, c1, 0, 0,\n    0, 0, c1, 0,\n    0, 0, 0, c1\n);\n\nvoid main() {\n//    vec4 test = uProjection * vec4(vTest.xyz, 1.0);\n    fPosition = vec4(vPosition.xyz, 1.0);\n    gl_Position = fPosition;\n    fFillColor = vFillCollor / 255;\n//    fFillColor = vec4(0, vFillCollor.g/255, 0, 1);\n}\x00"

const TextFragShaderStr = "#version 330\nin vec2 TexCoords;\nout vec4 color;\n\nuniform sampler2D text;\nuniform vec3 textColor;\n\nvoid main()\n{\n    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(text, TexCoords).r);\n    color = vec4(textColor, 1.0) * sampled;\n}\x00"

const TextVertexShaderStr = "#version 330\nlayout (location = 0) in vec3 aPos;\nlayout (location = 1) in vec2 aTexCoord;\n\nout vec2 TexCoords;\n\nuniform mat4 projection;\n\nvoid main()\n{\n    gl_Position = vec4(aPos, 1.0);\n    TexCoords = aTexCoord;\n}\x00"


