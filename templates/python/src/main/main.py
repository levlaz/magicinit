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
from dagger import dag, function, object_type


@object_type
class Python:
    
    @function
    def container_echo(self, string_arg: str) -> dagger.Container:
        """Returns a container that echoes whatever string argument is provided"""
        return dag.container().from_("alpine:latest").with_exec(["echo", string_arg])

    @function
    async def grep_dir(self, directory_arg: dagger.Directory, pattern: str) -> str:
        """Returns lines that match a pattern in the files of the provided Directory"""
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_mounted_directory("/mnt", directory_arg)
            .with_workdir("/mnt")
            .with_exec(["grep", "-R", pattern, "."])
            .stdout()
        )
