Autonity Docs
============

All user facing documentation for Autonity is held under this directory.

The content is held in ``.rst`` files and we use Sphinx to generate a static
site from all of the content.

Building the documentation
==========================


Install Prerequisites
---------------------

The project is configured with a Makefile, to use it you will need to install
make.

Once make is installed, you can install further system prerequisites.

Linux::

    make install-prerequisites-linux

Mac::

    make install-prerequisites-mac

If the above does not work for you, you will need to have installed python3
with venv and pip3.

Build site
----------

::

    make serve

Will install python module depnenencies, build the site and serve it.

If you update any files under ``source`` or modify any python depnenencies in
``requirements.txt`` running ``make serve`` again will update requirements and
rebuild the site as needed.
