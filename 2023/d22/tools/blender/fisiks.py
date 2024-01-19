import bpy
from collections import defaultdict
from bpy import data as D
from bpy import context as C
import mathutils
from math import *

def get_axis_max_mins(obj):
    worldmat = obj.matrix_world
    bbox = [worldmat @ mathutils.Vector(corner) for corner in obj.bound_box]
    min_coords = mathutils.Vector((min([b[i] for b in bbox]) for i in range(3)))
    max_coords = mathutils.Vector((max([b[i] for b in bbox]) for i in range(3)))
    return min_coords, max_coords

def is_colliding(min1, max1, min2, max2):
    for i in range(3):
        if max1[i] <= min2[i] or min1[i] >= max2[i]:
            return False
    return True
        


# Current scene's fps
frame_rate = C.scene.render.fps
start_frame = 200

# Gravity acceleration
grav = -9.8
grav_per_frame = grav/frame_rate

# Selected collections
collection_name = 'Appear and scripted physics'
plane_name = 'StackGroundPlane'

# Get the objects
bricks = D.collections[collection_name].objects
plane = D.objects[plane_name]
collidables = defaultdict(list)
collidables[0].append(plane)

# kdcol = mathutils.kdtree.KDTree(1000)
# bbox = [plane.matrix_world @ p for p in plane.bound_box]
# for bp in bbox:
#     kdcol.insert(bp.co, )

# The objects still in motion will be selected
# Clear selection
for o in C.selected_objects:
    o.select_set(False)

# Select bricks
for brick in bricks:
    brick.select_set(True)

# Start at frame start
C.scene.frame_set(start_frame)

# Make initial keyframe for all selected
for brick in C.selected_objects:
    brick.keyframe_insert(data_path="location", index=2)
# Advance one frame
C.scene.frame_set(C.scene.frame_current + 1)

# The current speed (since all blocks will fall at the same speed)
current_velocity = 0.0
translation_vector = mathutils.Vector([0.0,0.0,0.0])
while len(C.selected_objects) > 0:
    # Accelerate and move
    current_velocity += grav_per_frame
    translation_vector[2] = current_velocity / frame_rate
    bpy.ops.transform.translate(value=translation_vector)
    print(f"Frame {C.scene.frame_current} - {len(C.selected_objects)} moving - {translation_vector[2]} translation")
    
    # Update?
    C.view_layer.update()

    # Collision resolution
    while True:
        collisions = list()
        # Brute force
        for moving_brick in C.selected_objects:
            minbrick, maxbrick = get_axis_max_mins(moving_brick)
            minbrickz = floor(minbrick[2]/30)
            for collidable in collidables[minbrickz]:
                cmin, cmax = get_axis_max_mins(collidable)
                if is_colliding(minbrick, maxbrick, cmin, cmax):
                    collisions.append((moving_brick, collidable))
                    break
        if not len(collisions):
            break

        # Adjust
        for brick, col in collisions:
            newz = round(brick.location[2])
            if brick.dimensions[2] % 2 != 0:
                newz += 0.5
            brick.location[2] = newz
            thething = floor(newz/30)
            brick.select_set(False)
            print(thething)
            collidables[thething].append(brick)
            brick.keyframe_insert(data_path="location", index=2)
        C.view_layer.update()
    
    # Update
    C.view_layer.update()
    for brick in C.selected_objects:
        brick.keyframe_insert(data_path="location", index=2)
    # Advance one frame
    C.scene.frame_set(C.scene.frame_current + 1)
