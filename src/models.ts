import { MessageInterface, OperationsInterface, SchemaInterface } from "@asyncapi/parser"
import { TemplateResult, catchDiagnostics, diagnostic, template } from "./templates"
import { messageId, toMixedCaps } from "./utils"
import { uniq } from "underscore"

export type ModelsOptions = {
    package: string,
    version: string
}

export type Model = {
    implementation: string,
    dependencies: string[]
}

export type Models = Model[]

export type ModelsData = {
    package: string
    version: string

    dependencies: string[]
    models: string[]
    messages: {
        id: string
        name: string
        payload: string
    }[]
}

export default function modelsTemplate(operations: OperationsInterface, opts?: ModelsOptions): TemplateResult<ModelsData> {
    try {

        var data: ModelsData = {
            package: opts?.package ?? "asyncapi",
            version: opts?.version ?? "v0.0.0",
            dependencies: [],
            models: [],
            messages: []
        }

        var messagesPayloads = operations
            .all()
            .map(operation => operation.messages()
                .all()
                .map(message => message.payload())
                .filter(message => !!message) as SchemaInterface[]
            )
            .flat()

        var anonymousPayloads = messagesPayloads.filter(payload => isAnonymousSchema(payload))
        if (anonymousPayloads.length > 0) {
            throw anonymousPayloads.map(payload => diagnostic(`payload '${payload.id()}' is missing 'x-parser-schema-id'`))
        }

        var models = uniq(messagesPayloads
            .map(msg => implementModel(msg))
            .flat(),
            false, model => model.implementation 
        )
        data.dependencies = uniq(models
            .map(model => model.dependencies)
            .flat()
        )
        data.models = models
            .map(model => model.implementation)

        // Based on the operations of this spec, we
        // generate the necessary Message from the
        // different channels. Messages that are
        // send to or received by multiple channel
        // are rendered once.
        data.messages = uniq(
            operations
            .all()
            .map(operation =>
                operation
                    .messages()
                    .all()
                    .map(message => {
                        var id = messageId(message)

                        return {
                            id,
                            name: toMixedCaps(id),
                            payload: toMixedCaps(message.payload()!.id())
                        }
                    })
            )
            .flat()
        , false, message => message.id)

        // assert content type
        var unsupportedContentTypes = operations
            .all()
            .map(operation => operation.messages().all())
            .flat()
            .map(message => {return {contentType: message.contentType(), message: message.id()}})
            .filter(message => message.contentType !== "application/json")
        if (unsupportedContentTypes.length > 0) {
            throw unsupportedContentTypes.map(ct => diagnostic(!!ct.contentType ? `unsupported content type '${ct}' in '${ct.message}'`: `missing content type in ${ct.message}`))
        }

        return template("./models", data)

    } catch(e) {
        return catchDiagnostics(e)
    }
}

function implementModel(schema: SchemaInterface): Models {
    var type = schema.type()
    if (type !== "object" && type !== "array") {
        throw diagnostic(`expected model ${schema.id()} to be object or array, but got ${type}`)
    }

    if (isAnonymousSchema(schema)) {
        throw diagnostic(`schema is missing x-parser-schema-id`)
    }

    var implementation: string[] = []
    var dependencies: string[] = []
    var dependentModels: Models = []
    switch(type) {
        case "object":
            // iterate over properties and implement each property
            var objProps = schema.properties()
            if (objProps === undefined) {
                throw diagnostic(`object ${schema.id()} is missing properties`)
            }

            var requiredProps = schema.required() ?? []
            var props = Object.keys(objProps)
                .map(propName => {
                    return {
                        propName: propName,
                        id: isAnonymousSchema(objProps![propName]) ? propName : objProps![propName].id()
                    }
                })
                .map(prop => implementProperty(prop.id, objProps![prop.propName], requiredProps.includes(prop.propName)))

            props.forEach(prop => {
                dependentModels = dependentModels.concat(prop.models)
                dependencies = dependencies.concat(prop.dependencies)
            })

            implementation = [
                `type ${toMixedCaps(schema.id())} struct {`,
                props
                    .map(prop =>
                        prop.lines
                            .map(line => "\t" + line)
                            .join("\n")
                    )
                    .join("\n"),
                "}"
            ]

            break;

        case "array":
            // implement item of array
            var arrItem = schema.items()
            if (arrItem === undefined || arrItem instanceof Array) {
                throw diagnostic(`items of array ${schema.id()} is either undefined or an array`)
            }

            if (isAnonymousSchema(arrItem)) {
                throw diagnostic(`item schema of array ${schema.id()} is missing x-parser-schema-id`)
            }

            var item = implementProperty(arrItem.id(), arrItem, true)
            dependentModels = dependentModels.concat(item.models)
            dependencies = dependencies.concat(item.dependencies)

            implementation = [
                `type ${toMixedCaps(schema.id())} []${item.name}`
            ]
            break;

        default:
            throw new Error("type should not exist")
    }

    return dependentModels.concat([
        {
            implementation: implementation.join("\n"),
            dependencies
        }
    ])
}

function implementProperty(
    name: string,
    schema: SchemaInterface, 
    required: boolean
): {name: string, lines: string[], dependencies: string[], models: Models} {
    var type = schema.type()
    if (typeof type !== "string") {
        throw diagnostic("type is not a string")
    }
    type = type.toLowerCase()

    var format = schema.format()?.toLowerCase()
    var goType: string
    var dependencies: string[] = []
    switch(type) {
        case "integer":
            switch(format) {
                case "int32":
                    goType = "int32"
                    break;
                case "int64":
                    goType = "int64"
                    break;
                default:
                    goType = "int"
                    break;
            }
            break;
        case "string":
            switch(format) {
                case "date":
                    dependencies = ["time"]
                    goType = "time.Time"
                    break
                case "date-time":
                    dependencies = ["time"]
                    goType = "time.Time"
                    break
                case "uuid":
                    dependencies = ["github.com/google/uuid"]
                    goType = " uuid.UUID"
                    break
                default:
                    goType = "string"
            }
            break;
        case "number":
            switch(format) {
                case "float":
                    goType = "float32"
                    break;
                case "double":
                    goType = "float64"
                    break;
                default:
                    goType = "float64"
                    break;
            }
            break;
        case "boolean":
            goType = "bool"
            break;
        case "object":
            // Ignore for now and handle later
            goType = ""
            break;
        case "array":
            // Ignore for now and handle later
            goType = ""
            break;
        default:
            throw diagnostic(`unknown type ${type}`)
    }

    var lines: string[]
    var models: Models = []
    switch(type) {
        case "object":
        case "array":
            models = implementModel(schema)
            // Our goType is the name of the schema
            goType = toMixedCaps(name)
            break;
    }
    
    lines = [
        ...schema.description() ? [`// ${schema.description()}`] : [],
        `${toMixedCaps(name)} ${required ? "" : "*"}${goType} \`json:"${name}${required ? "" : ",omitempty"}"\``
    ]

    return {
        name: toMixedCaps(name),
        lines,
        dependencies,
        models
    }
}

function isAnonymousSchema(schema: SchemaInterface): boolean {
    return schema.id() === undefined ? true : schema.id().includes("anonymous-schema")
}