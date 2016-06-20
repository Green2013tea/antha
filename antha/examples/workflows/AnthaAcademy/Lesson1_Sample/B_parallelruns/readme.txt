antharun --driver

1. 
There are many additional flags which may be used with the antharun command. 
To see the full list type antharun --help on the command line.

2. 
To select which driver port to connect to add the --driver flag as shown above (making sure it matches the driver port you’ve served in a separate terminal). 
By default the driver port will be localhost:50051

'antharun --driver localhost:50051'

Before running this command you need to run the driver using the following command in a separate terminal:

'PipetMax'

3. The manualLiquidhandlingdriver would work in the same way
You can get this from source code before running in a separate terminal

go get github.com/antha-lang/manualLiquidHandler

Running it:


cd server
go build ./...
./server





