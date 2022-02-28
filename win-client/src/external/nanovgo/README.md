NanoVGo-GL4
=============

Modified version of `NanoVGo <https://github.com/shibukawa/vg4go>`_. These are the main changes:
* Removed support for mobile devices
* Replaced github.com/goxjs/gl with github.com/go-gl/gl, and github.com/goxjs/glfw with github.com/go-gl/glfw
* Upgraded to OpenGL 4.1 Core Profile
* Reduced some wrapping code used in gl_backend
* Removed the original documentation since the code has diverged. Look at the original documentation for more info
* Exposed the GL initialization to the app
* Added support for golang modules






---------------------------------------------------------------------
The original project's readme follows


NanoVGo
=============

Pure golang implementation of `NanoVG <https://github.com/memononen/nanovg>`_. NanoVG is a vector graphics engine inspired by HTML5 Canvas API.

`DEMO <https://shibukawa.github.io/nanovgo/>`_

API Reference
---------------

See `GoDoc <https://godoc.org/github.com/shibukawa/nanovgo>`_

Porting Memo
--------------

* Root folder ``.go`` files

  Ported from NanoVG.

* ``fontstashmini/fontstash_mini.go``

  Ported from `fontstash <https://github.com/memononen/fontstash>`_. It includes only needed functions.

* ``fontstashmini/truetype``

  Copy from ``https://github.com/TheOnly92/fontstash.go`` (Public Domain)

License
----------

zlib license

Original (NanoVG) Author
---------------------------

* `Mikko Mononen <https://github.com/memononen>`_

Author
---------------

* `Yoshiki Shibukawa <https://github.com/shibukawa>`_

Contribution
----------------

* Moriyoshi Koizumi
* @hnakamur2
* @mattn_jp
* @hagat
* @h_doxas
* FSX
