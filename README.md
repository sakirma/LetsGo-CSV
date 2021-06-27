## General
### How to run
1. Execute GenerateCSVInput.exe | `$ GenerateCSVInput <amount> <type=1/2>`
2. An Input.csv file will be generated. This file has to be placed inside `GenerateCSVOutput`
3. Run GenerateCSVOutput.exe
### Architecture
My idea for the architecture was to split the tool in different modules. This way I could have the Reader module be
imported and used by any other application. It is straight and simple.

* GenerateCSVInput: Generates a CSV input file which can be used
* GenerateCSVOutput: Generates a CSV output file by using the Reader module
* Reader: Reads a CSV input file and provides method to get costs for usage.
    * Most important methods for this assignment are in this class: Init & GetNextCost

### Linter & Unit testing
* Linter has been implemented through Github Actions. This uses the script within the `tools` folder to get all the modules of the repository. <br> Unfortunately, `golangci-lint` does not support multiple modules within the subdirs. An issue exists of this on Github but it still is not fixed yet <br>https://github.com/golangci/golangci-lint/issues/828
* Unit tests have been done through the IDE (Goland). There is only one file that checks the core of the application, which is the Reader module.

## Assumptions

### CSV Input Generation

[Gas and Electricity usage source](https://www.engie.nl/product-advies/gemiddeld-energieverbruik) <br>
A random value will be chosen for each of the type readings:

* Gas: 2 people = 1710m3 ever year| `1710 / 525600 * 15 = 0.048` = m3 / minutes in year * 15 min
* Electricity: 2 people = 2860kWh ever year | `2860000 / 525600 * 15 * 1000 = 81.621` = Wh / minutes in year * 15 min

Will add or subtract some random value to these readings.
incorrect readings.

Another assumption that I have made is that the `metering_point_id` is always valid in the input file.

### CSV

There is a chance that the starting rows in the CSV has invalid readings. This gave me the reason to remove every
invalid readings until I get to find a valid one.  <br>
I assume removing the first invalid readings until the next valid readings wouldn't do harm to the result, because the
CSV files has millions of data for readings and removing the first 10 or 100 readings would still give an accurate
result at the end.

### Reading

* It doesn't mention how the metering points work exactly. So I have assumed that the metering points only count up the
  readings compared to the previous reading. This gives a usage or costs that makes much more sense.
* The missing usage is consumed linearly and I was afraid that this would cause the usage go further down below the 0 or
  above the 100. Which is the reason why I am clamping the reading between 0 and 100 whenever I encounter an invalid data.
