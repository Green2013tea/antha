Other antharun flags:


antharun --parameters --workflow

By default the antharun command uses a parameters file named parameters.yml and a workflow file named workflow.json. 
If these files are named differently you’ll need to use the --parameters and/or --workflow flags to specify which files to use.

1.
To run the parameters found in this folder you'll need to run this:

antharun --parameters parameters.json --myamazingworkflow.json

_____________


antharun --inputPlateType

2. e.g. antharun --inputPlateType greiner384
This allows the type of input plate to be specified from the list of available plate types in github.com/antha-lang/antha/microArch/factory/make_plate_library.go

 
_____________

antharun --inputPlates 

3. e.g. antharun --inputPlates inputplate.csv 
This allows user defined input plates to be defined. If this is not chosen antha will decide upon the layout.
More than one inputplate can be defined: this waould be done like so:
antharun --inputPlates assemblyreagents.csv --inputPlates assemblyparts.csv