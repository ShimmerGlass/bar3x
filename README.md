# bar3x

![image](https://user-images.githubusercontent.com/1712219/86497905-c0216180-bd83-11ea-85e5-e4ed926d2d50.png)

`bar3x` is an easy to use and powerful status bar for your Linux desktop written in Golang.
Linux status bars typically choose between highly customizable text and rich graphics. `bar3x` takes a slightly different approach by providing a simple markup language for its configuration that allows for :

- Easy addition of new built-in modules
- Rich external command modules
- Customizable graphics

## Installing

Download the [latest release](https://github.com/ShimmerGlass/bar3x/releases/latest)

Or, if you have the Golang toolchain installed :

```
go get github.com/shimmerglass/bar3x
```

### Building

Building bar3x requires [Go](https://golang.org) and [go-bindata](https://github.com/jteeuwen/go-bindata).  go-bindata is go-get-able:

```
go get github.com/jteeuwen/go-bindata
git clone github.com/ShimmerGlass/bar3x
cd bar3x
go generate
go build .
# Copy ./bar3x to a directory in your $PATH
```

### Dependencies

- [libcairo](https://www.cairographics.org/) : should already installed as it is used by GTK, otherwise:
  - Debian/Ubuntu: `apt install libcairo2`
  - Fedora: `yum install cairo`
  - Arch: `pacman -S cairo`

## Quick Start

`bar3x` comes with a default configuration, and can be customized using `bar3x -config config.yaml` to change part or all of these parameters.

Find the list of **available modules** in the [Wiki](https://github.com/ShimmerGlass/bar3x/wiki/Modules).

Here is an example config :

```yaml

# colors configuration
bg_color:            "#17191e" # bar background
text_color:          "#d4e5f7" # general text
accent_color:        "#1ebce8" # icons and UI elements such as bars
neutral_color:       "#37393e" # background elements such as module separators and background graphs
neutral_light_color: "#90949d" # used for less important text such as units

# modules can be placed on the left, center and right of the bar
# choose the modules you want in each <ModuleRow>
bar_left: |
  <ModuleRow>
    <Volume />
  </ModuleRow>

bar_center: |
  <ModuleRow>
    <DateTime />
  </ModuleRow>

bar_right: |
  <ModuleRow>
    <Interface Iface="enp3s0" />
    <CPU />
    <RAM />
    <DiskUsage MountPoint="/" />
  </ModuleRow>
```

## Customizing

Much more that the colors can be customized, see this [Wiki page](https://github.com/ShimmerGlass/bar3x/wiki/Customizing) on how to cutomize your bar.
