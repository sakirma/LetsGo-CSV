# Readme

## Assumptions

### CSV

There is a chance that the starting rows in the CSV has invalid readings. <br> 
I assume removing the first invalid readings until the next valid readings wouldn't do harm to the result, since that the CSV files has millions of data for readings and removing the first 10 or 100 readings would still give an accurate result at the end. 

First I created a csv generator in a module. <br>
The goal of this choice is to be able to:

* Find or learn a package that can generate the csv files
* Reuse this package for the output
* Use this package for the unit testing.

#### Generation

[Gas and Electricity usage source](https://www.engie.nl/product-advies/gemiddeld-energieverbruik) <br>
A random type and value will be chosen and the average usage of each type will be used for this:

* Gas: 2 people = 1710m3 ever year| `1710 / 525600 * 15 = 0.048` = m3 / minutes in year * 15 min
* Electricity: 2 people = 2860kWh ever year | `2860000 / 525600 * 15 * 1000 = 81.621` = Wh / minutes in year * 15 min

Will add or substract some random value to these readings. If an option is given to the program, then it will create
incorrect readings.

#### Unit testing

TODO (Will use random with a seed)

### Reading the Input

Since I am able to generate a CSV file, I will read this and convert it into usage using the following algorithm: <br>
`usage = reading * (i + 1) - reading * i`

This will be done in an other module with multiple packages.

* One package to read the CSV.
* One package will convert the readings into usage.
    * It will also check if it is valid.
