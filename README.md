# aapi-codegen

An Asyncapi code generator that generates a watermill plugin.

## Examples

Generate code for streetlamp example by executing the following command:  

```
npm run dev -- -p asyncapi -m examples/streetlamp/asyncapi/models.gen.go -o examples/streetlamp/asyncapi/operations.gen.go examples/streetlamp/streetlamp.asyncapi.yml
```
```
npm run dev -- -p asyncapi -g models,operations,postbox -b examples/streetlamp_using_forwarder/asyncapi/postbox.gen.go -m examples/streetlamp_using_forwarder/asyncapi/models.gen.go -o examples/streetlamp_using_forwarder/asyncapi/operations.gen.go examples/streetlamp_using_forwarder/streetlamp.asyncapi.yml
```
