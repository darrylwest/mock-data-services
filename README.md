# Automated Test Services

```
 __  __         _     ___       _          ___              _           
|  \/  |___  __| |__ |   \ __ _| |_ __ _  / __| ___ _ ___ _(_)__ ___ ___
| |\/| / _ \/ _| / / | |) / _` |  _/ _` | \__ \/ -_) '_\ V / / _/ -_|_-<
|_|  |_\___/\__|_\_\ |___/\__,_|\__\__,_| |___/\___|_|  \_/|_\__\___/__/
```

## Overview

Mock Data Services are designed to intercept, record and play-back standard request/response patterns, whether REST or SOAP.  The L4 proxy listens for connections, stores incoming requests then sends the request to the remote target and stores the target's response.  Stored responses can then be used to short-circuit the full round trip to the target by returning qualified responses to recognized requests.

## Installation

### Download

###### darryl.west | 2018.02.28

