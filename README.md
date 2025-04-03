# RingCLI 0.1.0

CLI access to data stored on the Colmi R02 smart ring.

![Colmi R02 smart ring](./images/r02_001.webp)

## Compilation

1. `cd RingCLI`
1. `go mod tidy`
1. `go build -o build/ringcli`

## Usage

### Utilities

#### Scan for Rings

You require your ring’s BLE address to issue most sub-commands (but see [Bind a Ring](#bind-a-ring)). Obtain it by scanning for rings:

```shell
ringcli utils scan --first
```

If your ring is in range and powered on (ie. it has been placed the the charger at least once), you should see it listed and you can copy its address.

**Note** The `--first` in the command above causes `ringcli` to halt scanning on the first ring it finds. If you have multiple rings, do not include the `--first` switch. `ringcli` will now list all of them — or, at least, those it can detect within its 60-second scan window.

#### Set Ring Time

If you haven’t used another app to set the date and time on your ring, run:

```shell
ringlci utils settime --address {your ring BLE address}
```

You will need to do this to initialise your ring for use if you have not done so already.

#### Get Ring Info

With your ring’s address you can now obtain more information about it, including its battery state:

```shell
ringlci utils info --address {your ring BLE address}
```

#### Get Battery State

To just get the ring’s battery state, issue:

```shell
ringlci utils battery --address {your ring BLE address}
```

#### Get and Set Periodic Heart Rate Sampling

To enable periodic heart rate readings, issue:

```shell
ringlci utils setheartrate --address {your ring BLE address} --period 60 --enable
```

The period is in minutes and must be in the range of 1 to 255. Setting the period to zero disables periodic readings, as does using the `--disable` switch (unless `--enable` has been included too).

```shell
ringlci utils setheartrate --address {your ring BLE address} --disable
```

This call gets the current state:

```shell
ringlci utils getheartrate --address {your ring BLE address}
```

#### Locate a Ring

To locate your ring, if you’re unsure where it is, issue:

```shell
ringlci utils find --address {your ring BLE address}
```

This will flash the ring’s green LED twice. If that’s not sufficient, add the `--continuous` switch:

```shell
ringlci utils find --address {your ring BLE address} --continuous
```

To flash the LED every two seconds until you cancel. The ring will cease flashing after 200 seconds to preserve ring battery power.

#### Bind a Ring

To save having to enter the `--address` option every time, you can ‘bind’ your ring to your system. This retains the ring’s BLE address across runs.

```shell
ringlci utils bind --address {your ring BLE address}
```

To check a binding, run:

```shell
ringlci utils bind --show
```

You can only bind one ring: to check on other rings, pass in a temporary BLE address with the `--address` option.

#### Shutdown a Ring

To shut the ring down, issue:

```shell
ringlci utils shutdown --address {your ring BLE address}
```

This powers down the ring. The ring can be restarted by placing it on its charger.

### Data

#### Daily Steps, Calories Burned, Distance Moved

The above sub-commands are provided by the `utils` command. `ringcli` also has a `data` command:

```shell
ringlci data steps --address {your ring BLE address}
```

This will output your current daily step count, activity based calorie burn and the distance you have travelled (estimate).

#### Daily Heart Rate Log

The `data` sub-command `heartrate` will retrieve the day’s heart rate readings (from midnight to the current time):

```shell
ringlci data heartrate --address {your ring BLE address}
```

## Copyright and Licence

`ringcli` is copyright © 2025, Tony Smith (@smittytone). The code is made available under the terms of the MIT licence.
