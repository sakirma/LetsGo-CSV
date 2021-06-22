# Readme

## Assumptions

### CSV

My first assumption was to be able to create a csv generator with a custom package. <br>
The goal of this choice is to be able to:

* Find or learn a package that can generate the csv files
* Reuse this package for the output
* Use this package for the unit testing.

#### Generation

[Gas and Electricity usage source](https://www.engie.nl/product-advies/gemiddeld-energieverbruik) <br>
A random type and value will be chosen and the average usage of each type will be used for this:

* Gas: 2 people = 1710m3 ever year| `1710 / 525600 * 15 = 0.052` = m3 / minutes in year * 15 min
* Electricity: 2 people = 2860kWh ever year| `2860000 / 525600 * 15 * 1000 = 81.621` = kWh / minutes in year * 15 min

Will add or substract some random value to these readings. If an option is given to the program, then it will create
incorrect readings
