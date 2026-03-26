# Introduction
MSGGW(message gateway) is a Golang program that is filtering and transforming messages between different brokers.
Currently it is aimed to nats-internal and nats-public messages.

# File structure
```
.
├── Makefile                    
├── README.md
├── app
│   └── msggw        # code entrance
├── bin              # binary build output
├── build
│   └── docker       # Docker file
├── component.json   # product Version info
├── configs          # configuration file example
├── deployment       # k8s deployment script
├── doc              # related documents
├── go.mod
├── go.sum
├── pkg              # Go packages
│   ├── config       # for reading config file
│   ├── funcs        # Func support symbols with format `{func.*}`
│   ├── model        # basic interfaces and consts
│   ├── operator     # Operator support symbols with IS, NOT, MATCH ...
│   ├── parser       # Parser for parsing symbols
│   ├── service      # Generating business logic from a config file, and processing incoming messages with these configured logics.
│   └── symbol       # symbol definitions for parser
└── tests
```
# How to build
1. Binary
    Do `make build`, and find binary in ./bin.
2. Container
    Do `make container` after step 1

# How to run
1. Binary
    Run with command line `./bin/msggw ./configs/server_side.json`
2. Container locally
   `docker run msggw:latest --network host`
3. on K8s
   `kubectl apply -f ./deployment/msggw-deployment.yaml`
   `kubectl apply -f ./deployment/msggw-configmap.yaml`

# Configuration
1. Main structure
   ```json
    {
        "logLevel": string,  // currently it has options of `info` and `debug`.
        "props": object,     // properties that can be directly used in flow as env var. NOTE: K/V should be string.
        "brokers": arr,      // define brokers, name is used by flow source and send dest. NOTE: the name should be unique.
        "flows": arr         // processing logics
    }
   ```
2. Flows array
   Each flow means how to handle messages from a broker.
   ```json
   {
        "name": string,            // flow name
        "source": string,          // the name defined in `brokers.name` field
        "subcribes": string arr,   // topics and wildcards for filtering
        "payload": string,         // can be `msgbus`(msg from nats-internal) or `edgebus`(msg from edge side modules), 
        "branches": object arr     // each branch has filter, transform and send strategy, like if conditions in programming language
   }
   ``` 
3. Branch
   When a message is coming, the message will pass though each branch in array sequence. 
   Any match of a filter will let the message goes through the following process(transforms and sendTo), and won't process the following branches.
   ```json
   {
        "name": string,            // branch name
        "filters": arr,            // each string is a ternary expression(lval OP rval), result from expression will be logic AND together.
        "transforms": object arr,  // transform the data that will be send, optional keys: "{metadata.*}" or "topic" 
        "sendTo": object           // send to a broker(`dest` field) with specified payload(`payload` field)
   }
   ``` 
4. Ternary expression for `filter`
   the string is a form of `lval OP rVal`.
   lval can be type of:
   1. metadata
   2. topic
   OP can be value of:
   1. IS
   2. NOT
   3. MATCH
   rval can be of type:
   1. raw:      raw string
   2. topic:    message topic
   3. metadata: message metadata. Value is like `{metadata.pub_para}`
   4. func:     builtin functions. Value is like `{func.PubParaToTopic}`
   5. prop:     prop field of configuration. Value is like `{prop.deviceId}`
   6. keyword:  currently we only have `NULL` to imply an empty field.
   7. mix:      a mixed symbol has composed sub symbols. Value is like `{/.org.{metadata.oid}}`
5. Key and value for `transform`
   Key should be one of:
   1. metadata
   2. topic
   Value can be of type:(raw, topic, metadata, func, prop, keyword, mix) same as expression rvalue.
6. Key and value for `sendTo `
   ```json
   {
      "dest": string,     // find the value in broker name
      "payload": string   // can be edgebus(to edge side) or msgbus(to server side)
   }
   ```
