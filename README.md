# aapi-codegen

An Asyncapi code generator that generates a watermill plugin.

## Examples

Generate code for streetlamp example by executing the following command:  

```
npm run dev -- -p asyncapi -m examples/streetlamp/asyncapi/models.gen.go -o examples/streetlamp/asyncapi/operations.gen.go examples/streetlamp/streetlamp.asyncapi.yml
```