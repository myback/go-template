import ctypes
import json
import os
from pathlib import Path

lib_path = os.path.join(os.path.dirname(__file__), 'bind/template.so')
lib = ctypes.cdll.LoadLibrary(lib_path)

_render = lib.Render
_render.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
_render.restype = ctypes.c_void_p

_render_from_values_file = lib.RenderFomValuesFile
_render_from_values_file.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
_render_from_values_file.restype = ctypes.c_void_p


def render(src_file: Path, values: dict) -> bytes:
    return ctypes.string_at(
        _render(
            str(src_file.resolve()).encode(),
            json.dumps(values).encode()
        )
    )


def render_from_values_file(src_file: Path, values_file: str) -> bytes:
    return ctypes.string_at(
        _render_from_values_file(
            str(src_file.resolve()).encode(),
            values_file.encode()
        )
    )
