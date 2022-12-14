<samp>

# telmon

## What is it?

Telmon is a simple golang daemon that will monitor a given telnet connection from a telnet client to a telnet server. 
It reports any loss in connection via `email` to a list of receivers.

## How to use it.

### 1. Installation.
Install telmon binary from GitHub. In this tutorial, we will be using the linux OS builds.

```shell
wget https://github.com/jackkweyunga/telmon/releases/download/v0.1/telmon_0.1_Linux_x86_64.tar.gz
```

Extract to `/usr/local/bin/`

```shell
sudo tar -C /usr/local/bin -xzf telmon_0.1_Linux_x86_64.tar.gz telmon
```

test it. If you see a message as one below, shows that it's working.
```shell
$ telmon
2022/10/07 09:49:57 Config File ".telmon-config" Not Found in ...
```

- visit the [releases page](https://github.com/jackkweyunga/telmon/releases) for system specific downloads.

In order to use telmon, create a configuration file ``.telmon-config.yaml`` in any directory in which you will run the 
telmon binary.

### 2. Configurations
create a `.telmon-config.yaml` configuration file in any directory accessible by the user. Home directory is better.
```shell
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
sample output
```shell
$ telmon
Monitoring Service started successfully ...
[2022-10-07T09:56:09+03:00] running
[TelnetClient] Trying to connect to telnet server.
[TelnetClient] Encountered errors while connecting:
dial tcp 192.241.152.67:23: connect: connection refused
Reporting this error
[TelnetClient] Report sent successfully
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
