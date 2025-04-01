# MOD+MOI Converter

This is a helpful util that converts MOD+MOI files into a desired format, while keeping their exif data

## Installation

If you have go setup, simply do `go install github.com/shadiestgoat/moiconv@latest`

If you do not, go to the github releases and download the latest release for your platform

## Usage

> [!IMPORTANT]
> #### For Windows Users
> Whenever you see `moiconv` in this section, you have to use `moiconv.exe`.
> Also, the recommendation is to either use wsl, or Powershell

Convert `./directory` full of MOV & MOI files and place them into `path/to/output`: `moiconv -o path/to/output ./directory`

Convert `./directory` which is structured as a bunch of sub folders full of MOV & MOI files, keeping the same file structure in output directory: `moiconv --recursive -o path/to/output ./directory`

Convert `./directory` which is structured as a bunch of sub folders full of MOV & MOI files, putting all the files into `path/to/output`: `moiconv --recursive --flat -o path/to/output ./directory`

Specify the format of output files: `moiconv --format 'mov' -o path/to/output ./directory`
