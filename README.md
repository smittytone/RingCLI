# RingCLI

CLI access to data stored on the Colmi R02 smart ring.

![Colmi R02 smart ring](./images/r02_001.webp)

## Compilation

1. `cd RingCLI`
1. `go mod tidy`
1. `go build -o build/ringcli`

## Usage

### Utilities

At the moment, you require your ring’s BLE address to issue most sub-commands. Obtain it with:

```shell
ringcli utils scan --first
```

If your ring is in range and powered, you should see it listed and you can copy its address.

**Note** The `--first` in the command above causes `ringcli` to halt scanning on the first ring it finds. If you have multiple rings, do not include the `--first` switch. `ringcli` will now list all of them — or, at least, those it can detect in the 60-second scan window.

With your ring’s address you can now obtain more information about it, including its battery state:

```shell
ringlci utils info --address {your ring BLE address}
```

To locate your ring, if you’re unsure where it is, issue:

```shell
ringlci utils find --address {your ring BLE address}
```

This will flash the ring’s green LED twice. If that’s not sufficient, add the `--continuous` switch:

```shell
ringlci utils find --address {your ring BLE address} --continuous
```

To flash the LED every two seconds until you cancel. The ring will cease flashing after 200 seconds to preserve ring battery power.

To shut the ring down, issue:

```shell
ringlci utils shutdown --address {your ring BLE address}
```

This powers down the ring. The ring can be restarted by placing it on its charger.

### Data

The above sub-commands are provided by the `utils` command. `ringcloi` also has a `data` command:

```shell
ringlci data steps --address {your ring BLE address}
```

This will output your current daily step count, activity based calorie burn and the distance you have travelled (estimate).


## Copyright and Licence

`ringcli` is copyright © 2025, Tony Smith (@smittytone). The code is made available under the terms of the MIT licence.
