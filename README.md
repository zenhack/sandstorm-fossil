WIP port of [fossil][1] to [sandstorm][2]

Doesn't work yet.

In principle, Fossil should be a nice fit for sandstorm, but there's
some work to be done especially wrt. Integrating authentication. In
addition to the package manifest this repository contains the code
for a WIP reverse proxy that does the necessary voodoo to get auth
working.

Note that if you're viewing this on GitHub, this repository is an export
from a fossil repo that is the authoritative source. Once this app
is working, we can use it to host itself.

[1]: https://fossil-scm.org
[2]: https://sandstorm.io
