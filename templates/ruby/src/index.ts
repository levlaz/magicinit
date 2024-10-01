/**
A Dagger pipeline for Python generated via Magicinit™

                       _      _       _ _   
                      (_)    (_)     (_) |  
 _ __ ___   __ _  __ _ _  ___ _ _ __  _| |_ 
| '_ ` _ \ / _` |/ _` | |/ __| | '_ \| | __|
| | | | | | (_| | (_| | | (__| | | | | | |_ 
|_| |_| |_|\__,_|\__, |_|\___|_|_| |_|_|\__|
                  __/ |                     
                 |___/    
                 
We wish this example was written in Ruby, please join us in 
discord to help us make this a reality: https://discord.com/invite/dagger-io
 */
import { dag, Container, Directory, object, func, argument } from "@dagger.io/dagger"

@object()
class Ruby {
  rubyVersion: string = "latest"
  source: Directory

  /**
   * Module level arguments using constructor
   * @param source location of source code, defaults to current working dir
   *
   * more info on defaultPath: https://docs.dagger.io/manuals/developer/functions/#directories-and-files
   * more info on constructor: https://docs.dagger.io/manuals/developer/entrypoint-function/
   */
  constructor(@argument({ defaultPath: "/" }) source: Directory) {
    this.source = source
  }

  /**
   * Base container for Ruby project
   */
  @func()
  base(): Container {
    return dag.
      container().
      from(`ruby:${this.rubyVersion}`).
      withMountedDirectory("/src", this.source).
      withWorkdir("/src")
  }

  /**
   * Build the project
   */
  @func()
  build(): Container {
    return this.base().withExec(["bundle", "install"])
  }

  /**
   * Lint the project
   */
  @func()
  lint(): Container {
    return this.base().withExec(["bundle", "exec", "rubocop"])
  }

  /**
   * Test the project
   */
  @func()
  test(): Container {
    return this.base().withExec(["bundle", "exec", "rspec"])
  }
}
