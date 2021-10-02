layout (location = 0) in vec3 vPosition;
layout (location = 1) in vec2 vTexCoord;

out vec4 fPosition;
out vec2 fTexCoord;

void main()
{
    gl_Position = applyProjection(vPosition);
    fPosition = gl_Position;
    fTexCoord = vTexCoord;
}