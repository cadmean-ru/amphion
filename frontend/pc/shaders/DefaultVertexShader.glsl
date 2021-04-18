#version 330
layout (location = 0) in vec3 aPos;

out vec4 trashPos;

void main()
{
    trashPos = vec4(aPos.x, aPos.y, aPos.z, 1.0);
    gl_Position = trashPos;
}