# Pixie: A Blocking Proxy

Pixie is a proxy designed to block content. Based on URLs and content types.
Initially focusing on blocking by URL, but will add Content-Type based Blocking
later.

# Usecases

### Adblocking proxy

With https://easylist.to, we actually have a database of ad providers that we
can easily block!

We can use this proxy to block only based on the URL and also find providers who
are not on easylist.

### A General Blocking proxy

Sometimes we want to open only some URLs and block others, like image URLs.
We can use this to achieve that

# Usage

We have two lists in the config file.
* *Accept*: Any URL that matches any regex in this list will be allowed
* *Deny*: Any URL that doesn't match any accept regexes but matches a regex in
this list will be blocked with 403.

Any URL that doesn't match any of the regexes will be allowed.

# Running

Use the docker image:
```
docker run -v /path/to/config.json:/config.json -p 8080:8080 gouthamve/pixie:v0
```

### Credits:
* Built on https://github.com/vulcand/oxy.
* Some code taken from https://github.com/elazarl/goproxy.
