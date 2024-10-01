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
 */
import { dag, Container, Directory, object, func, argument } from "@dagger.io/dagger"

@object()
class Typescript {
  typescriptVersion: string = "latest"
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
   * Base container for Typescript project
   */
  @func()
  base(): Container {
    return dag.
      container().
      from(`ruby:${this.typescriptVersion}`).
      withMountedDirectory("/src", this.source).
      withWorkdir("/src")
  }

  /**
   * Build the project
   */
  @func()
  build(): Container {
    return this.base().withExec(["npm", "install"])
  }

  /**
   * Lint the project
   */
  @func()
  lint(): Container {
    return this.base().withExec(["npm", "lint"])
  }

  /**
   * Test the project
   */
  @func()
  test(): Container {
    return this.base().withExec(["npm", "test"])
  }
}
