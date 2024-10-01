"""A Dagger pipeline for Python generated via Magicinit™

                       _      _       _ _   
                      (_)    (_)     (_) |  
 _ __ ___   __ _  __ _ _  ___ _ _ __  _| |_ 
| '_ ` _ \ / _` |/ _` | |/ __| | '_ \| | __|
| | | | | | (_| | (_| | | (__| | | | | | |_ 
|_| |_| |_|\__,_|\__, |_|\___|_|_| |_|_|\__|
                  __/ |                     
                 |___/                      
 """

import dagger
from dagger import dag, function, object_type, DefaultPath
from typing import Annotated

@object_type
class Python:
    source: Annotated[dagger.Directory, DefaultPath("/"), Doc("Source directory of the project, defaults to the root of the repository")]
    python_version: Annotated[str, "latest", Doc("Python version to use, defaults to latest")]

    @function
    def base(self) -> dagger.Container:
        """Base container for Python project"""
        return (
            dag.
            container().
            from_(f"python:{self.python_version}").
            with_mounted_directory("/src", self.source).
            with_workdir("/src")
        )
    
    @function
    def build(self):
        """Build the project"""
        return (
            self.base().
            with_exec(["pip", "install", "-e", "."]) 
        )
    
    @function
    def lint(self):
        """Lint the project"""
        return (
            self.base().
            with_exec(["pip", "install", "ruff"]).
            with_exec(["ruff", "check", "."])
        )

    @function
    def test(self):
        """Test the project"""
        return (
            self.base().
            with_exec(["pip", "install", "pytest"]).
            with_exec(["pytest"])
        )