layout (location = 0) in vec3 vPosition;
layout (location = 1) in vec2 vTexCoord;

out vec2 fTexCoord;
out vec4 fPosition;

void main()
{
    gl_Position = vec4(vPosition, 1.0);
    fPosition = gl_Position;
    fTexCoord = vTexCoord;
}