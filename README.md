# aapi-codegen

An Asyncapi code generator that generates a watermill plugin.

## Install
```
npm i @fair-bytes/aapi-codegen -g
```

## Usage
```
Usage: aapi-codegen [options] <asyncapi>

Options:
  -p, --package <package>        golang package
  -g, --generate <generate>      what to generate, possible are: models,operations,postbox (default: "models,operations")
  -m, --models <models>          models file
  -o, --operations <operations>  operations file
  -b, --postbox <postbox>        postbox file
  -h, --help                     display help for command
```

## Examples

Generate code for the streetlamp examples by executing the following command:  

Simple streetlight:
```bash
# Execute example in development
npm run dev -- -p asyncapi -m examples/streetlamp/asyncapi/models.gen.go -o examples/streetlamp/asyncapi/operations.gen.go examples/streetlamp/streetlamp.asyncapi.yml

# Execute example in production
aapi-codegen -p asyncapi -m examples/streetlamp/asyncapi/models.gen.go -o examples/streetlamp/asyncapi/operations.gen.go examples/streetlamp/streetlamp.asyncapi.yml
```

Streetlight using the postbox pattern and the watermill forwarder component:
```bash
# Execute example in development
npm run dev -- -p asyncapi -g models,operations,postbox -b examples/streetlamp_using_forwarder/asyncapi/postbox.gen.go -m examples/streetlamp_using_forwarder/asyncapi/models.gen.go -o examples/streetlamp_using_forwarder/asyncapi/operations.gen.go examples/streetlamp_using_forwarder/streetlamp.asyncapi.yml

# Execute example in production
aapi-codegen -p asyncapi -g models,operations,postbox -b examples/streetlamp_using_forwarder/asyncapi/postbox.gen.go -m examples/streetlamp_using_forwarder/asyncapi/models.gen.go -o examples/streetlamp_using_forwarder/asyncapi/operations.gen.go examples/streetlamp_using_forwarder/streetlamp.asyncapi.yml
```
