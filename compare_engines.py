#!/usr/bin/env python3
from ctypes import cdll, c_uint, c_int

lib = cdll.LoadLibrary('./libgoengine.so')

# Create engines and board
engineA = lib.NewAlphaBetaEngine()
engineB = lib.NewRandomEngine()
board = lib.NewBoard(c_int(9))

# Call CompareEngines with valid handles
result = lib.CompareEngines(
    c_uint(engineA),
    c_uint(engineB),
    c_uint(board),
    c_int(1),      # firstPlayer
    c_int(100)     # maxMoves
)

# Print the result and name the winner
if result == 1:
    print("Winner: Engine A (Alpha-Beta)")
elif result == 2:
    print("Winner: Engine B (Random)")
elif result == 0:
    print("Result: Draw")
