# CVM for Windows

Removing the hassle of changing Composer versions in the CLI on Windows.

This package has a much more niche use case than nvm does. When developing on Windows and using the integrated terminal, it's quite difficult to get those terminals to _actually_ listen to PATH changes.

This utility changes that.

## Installation

Download the latest CVM version from the [releases page](https://github.com/wujingquan/cvm/releases/latest) .

Create the folder `%UserProfile%\.cvm\bin` (e.g. `C:\Users\Example\.cvm\bin`) and drop the cvm exe in there. Add the folder to your PATH.

## Commands
```
cvm list
```
Will list out all the available Composer versions you have installed

```
cvm path
```
Will tell you what to put in your Path variable.

```
cvm use 2.7.7
```
> [!NOTE]  
> Versions must have major.minor specified in the *use* command. If a .patch version is omitted, newest available patch version is chosen.

Will switch your currently active Composer version to Composer 2.7.7

```
cvm install 2.7
```
> [!NOTE]  
> The install command will automatically determine the newest minor/patch versions if they are not specified

Will install Composer 2.7 at the latest patch.

## Thanks
- Based on [PVM](https://github.com/hjbdev/pvm)