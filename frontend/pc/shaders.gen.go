
// +build darwin

package pc

// This file was auto-generated

var DefaultFragShaderStr = "#version 330 core\nout vec4 FragColor;\n\nuniform vec4 ourColor;\n\nvoid main()\n{\n    FragColor = ourColor;\n}\x00"

var DefaultVertexShaderStr = "#version 330 core\nlayout (location = 0) in vec3 aPos;\n\nout vec4 trashPos;\n\nvoid main()\n{\n    trashPos = vec4(aPos.x, aPos.y, aPos.z, 1.0);\n    gl_Position = trashPos;\n}\x00"

var EllipseFragShaderStr = "#version 330 core\nout vec4 FragColor;\n\nuniform vec4 ourColor;\n\nuniform vec3 tlPos;\nuniform vec3 brPos;\n\nin vec4 trashPos;\n\nvoid main()\n{\n    float a = (brPos.x - tlPos.x) / 2; //0.1\n    float b = (brPos.y - tlPos.y) / 2; //0.1\n    float x = trashPos.x;\n    float y = trashPos.y;\n    float xc = tlPos.x + a; // 0.3\n    float yc = tlPos.y + b; // -0.3\n    float x2 = x - xc;\n    float y2 = y - yc;\n    float c = x2 / a;\n    float d = y2 / b;\n    float res = c * c + d * d;\n    if (res <= 1)\n    FragColor = ourColor;\n    else\n    discard;\n}\x00"

var ImageFragShaderStr = "#version 330 core\n\nout vec4 FragColor;\n\nin vec2 texCoord;\n\nuniform sampler2D ourTexture;\n\nvoid main()\n{\n    FragColor = texture(ourTexture, texCoord);\n}\x00"

var ImageVertexShaderStr = "#version 330 core\nlayout (location = 0) in vec3 aPos;\nlayout (location = 1) in vec2 aTexCoord;\n\nout vec2 texCoord;\n\nvoid main()\n{\n    gl_Position = vec4(aPos, 1.0);\n    texCoord = aTexCoord;\n}\x00"

var TextVertexShaderStr = "#version 330 core\nlayout (location = 0) in vec3 aPos;\nlayout (location = 1) in vec2 aTexCoord;\n\nout vec2 TexCoords;\n\nuniform mat4 projection;\n\nvoid main()\n{\n    gl_Position = vec4(aPos, 1.0);\n    TexCoords = aTexCoord;\n}\x00"


