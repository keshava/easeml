import pluggy  # type: ignore

hookimpl = pluggy.HookimplMarker("easemlclient")
"""Marker to be imported and used in plugins (and for own implementations)"""