<samp>

# telmon

## What is it?

Telmon is a simple golang daemon that will monitor a given telnet connection from a telnet client to a telnet server. 
It reports any loss in connection via `email` to a list of receivers.

## How to use it.

### 1. Installation.
Install telmon binary from GitHub.

```shell
wget <github-path>.tar.gz
```

Extract to `/usr/local/bin/`

```shell
tar -C /usr/local/bin -xzf <github-path>.tar.gz telmon
```

- visit the releases page for system specific downloads.

In order to use telmon, create a configuration file ``.telmon-config.yaml`` in any directory in which you will run the 
telmon binary.

### 2. Configurations

```shell
$ mkdir telmon
$ cd telmon
$ nano .telmon-config.yaml
```

The content of the configuration file should be as follows.
```yaml
---
addr: example.com
port: <telnet-port>
email: <sender-email>
password: <sender-email-password>
receivers:
  - user-1@example.com
  - user-2@example.com
```

### 3. Run

```shell
$ telmon
```

### 4. Supervisord

Supervisord helps to monitor and control running processes
in any linux machine. For linux users, add a supervisord configuration
file to monitor telmon.

> create a file /etc/supervisord/conf.d/telmon.conf

Add the content below.

```editorconfig
[program:telmon]
directory=/path/to/your/telmon/directory
command=telmon
autostart=true
autorestart=true
stderr_logfile=/var/log/telmon.err.log
stdout_logfile=/var/log/telmon.out.log
```

Reload supervisord configurations.
```shell
supervisorctl reread
supervisorctl update
```

OR .. reload supervisor
```shell
supervisorctl reload
```

## Credits

- The [viper project](https://github.com/spf13/viper)
- The [go-telnet project](https://github.com/reiver/go-telnet)

</samp>