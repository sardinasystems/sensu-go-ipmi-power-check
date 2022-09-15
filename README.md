# sensu-go-ipmi-power-check

## Table of Contents
- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The sensu-go-ipmi-power-check is a [Sensu Check][6] that checks for the status of the power supply.
Needs freeipmi - ipmimonitoring.

## Files

- sensu-go-ipmi-power-check

## Usage examples

```
sensu-go-ipmi-power-check
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add sardinasystems/sensu-go-ipmi-power-check
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index](https://bonsai.sensu.io/assets/sardinasystems/sensu-go-ipmi-power-check).

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-go-ipmi-power-check
  namespace: default
spec:
  command: sensu-go-ipmi-power-check
  subscriptions:
  - system
  runtime_assets:
  - sensu-go-ipmi-power-check
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-go-ipmi-power-check repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sardinasystems/sensu-go-ipmi-power-check/blob/master/.github/workflows/release.yml
[5]: https://github.com/sardinasystems/sensu-go-ipmi-power-check/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[7]: https://github.com/sensu/check-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
