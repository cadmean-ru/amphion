layout (location = 0) in vec3 aPos;

out vec4 trashPos;

void main()
{
    trashPos = applyProjection(aPos);
    gl_Position = trashPos;
}