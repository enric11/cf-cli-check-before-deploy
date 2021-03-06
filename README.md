# PCF CLI Plugin Check Before Deploy
Check Before Deploy for PCF CLI Plugin

This repository provides a [Cloud Foundry CLI](https://github.com/cloudfoundry/cli) plugin for check the YAML file before launch deploy instructions.

## Building
You need to install GO.

The specific libraries used in this project are ready in the "vendor" folder.

After install external libraries, it's possible ejecute on windows "compile.bat" to generate plugin in different platforms (osx,windows,linux).

## Installing

To install the plugin in the `cf` CLI, first build it and then issue:
```bash
$ cf install-plugin -f <Your_OS>/check-before-deploy<Your_OS>

```

The plugin's commands may then be listed by issuing `cf help`.

To update the plugin, uninstall it as follows and then re-install the plugin as above:
```bash
$ cf uninstall-plugin check-before-deploy
```

## Testing
Run the tests as follows to check all commands (Actually 2):
```bash
$ cd cf-cli-check-before-deploy
$ cf check-before-deploy -file mta.yaml -all
```

## Results
![Image of execution](https://raw.githubusercontent.com/enric11/cf-cli-check-before-deploy/master/images/execution.png)



## License

The Check Before Deploy plugin for PCF CLI plugin is Open Source software released under the
[Apache 2.0 license](https://www.apache.org/licenses/LICENSE-2.0.html).
