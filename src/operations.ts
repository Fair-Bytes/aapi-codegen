import { ChannelInterface, MessageInterface, MessagesInterface, OperationInterface, OperationsInterface, SchemaInterface } from "@asyncapi/parser"
import { TemplateResult, catchDiagnostics, diagnostic, template } from "./templates"
import { channelAddressParameters, messageId, toMixedCaps } from "./utils"
import { uniq } from "underscore"

export type OperationsOptions = {
    package: string
}

export type Receivers = {
    id: string
    channel: {
        id: string
        address: string
    }
    parameters: string[]
    messages: {
        id: string
        name: string
        channel: string
    }[]
}[]

export type Senders = {
    id: string
    channel: {
        id: string
        address: string
    }
    parameters: string[]
    messages: {
        id: string
        name: string
        channel: string
    }[]
}[]

export type OperationsData = {
    package: string

    // contains receive operations
    receivers: Receivers

    // contains send operations
    senders: Senders

    // channels
    channels: {
        id: string
        parameters: {
            id: string
            name: string
        }[]
    }[]
}

export default function operationsTemplate(operations: OperationsInterface, opts?: OperationsOptions): TemplateResult<OperationsData> {
    try {

        var data: OperationsData = {
            package: opts?.package ?? "asyncapi",
            receivers: [],
            senders: [],
            channels: []
        }
    
        // render receive operations
        data.receivers = receivers(operations.filterByReceive())
    
        // render send operations
        data.senders = senders(operations.filterBySend())

        // render channels used by sender and receive operations
        data.channels = uniq(data.receivers.concat(data.senders)
            .map(op => {
                return {
                    id: toMixedCaps(op.channel.id),
                    parameters: op.parameters.map(paramId => {
                        return {
                            id: paramId,
                            name: toMixedCaps(paramId)
                        }
                    })
                }
            })
        , false, ch => ch.id)
    
        return template("./operations", data)

    } catch (e) {
        return catchDiagnostics(e)
    }   
}

function receivers(receiveOperations: OperationInterface[]): Receivers {
    return receiveOperations.map(receiver => {
        var id = receiver.id()
        if (id === undefined) {
            throw diagnostic("receive operation without id")
        }

        var channel = getChannel(receiver)

        var messages = receiver.messages()
            .all()
            .map(message => {
                var msgId = messageId(message)
                return {
                    id: msgId,
                    name: toMixedCaps(msgId),
                    channel: toMixedCaps(channel.channel.id())
                }
            })

        assertMessagesUseDefaultHeaders(receiver.messages())
        
        return {
            id: toMixedCaps(id),
            channel: {
                id: toMixedCaps(channel.channel.id()),
                address: channel.address
            },
            parameters: channel.parameters,
            messages
        }
    })
}

function senders(sendOperations: OperationInterface[]): Senders {
    return sendOperations.map(sender => {
        var id = sender.id()
        if (id === undefined) {
            throw diagnostic("send operation without id")
        }

        var channel = getChannel(sender)

        var messages = sender.messages()
            .all()
            .map(message => {
                var msgId = messageId(message)
                return {
                    id: msgId,
                    name: toMixedCaps(msgId),
                    channel: toMixedCaps(channel.channel.id())
                }
            })
        
        assertMessagesUseDefaultHeaders(sender.messages())
        
        return {
            id: toMixedCaps(id),
            channel: {
                id: toMixedCaps(channel.channel.id()),
                address: channel.address
            },
            parameters: channel.parameters,
            messages,
        }
    })
}

function getChannel(operation: OperationInterface): {channel: ChannelInterface, address: string, parameters: string[]} {
    var channels = operation.channels().all()
    if (channels.length != 1) {
        throw diagnostic(`expected exactly one channel for operation ${operation.id()}, but got ${channels.length}`)
    }
    
    var channel = channels.at(0)!
    var addr = channel.address()
    if (typeof addr !== "string") {
        throw diagnostic(`channel ${channel.id()} is missing address`)
    }

    var parameters = channel.parameters()
    var addrParameters = channelAddressParameters(addr)

    if (parameters.length !== addrParameters.length) {
        throw diagnostic(`channel parameter and parameters in address mismatch in length (${parameters.length} to ${addrParameters.length})`)
    }

    // validate parameters
    addrParameters.forEach(paramId => {
        var param = parameters.get(paramId)
        if (param === undefined) {
            throw diagnostic(`channel address contains parameter ${paramId} which is missing in parameters`)
        }

        var paramSchema = param.schema()
        if (paramSchema === undefined) {
            throw diagnostic(`no schema for parameter ${paramId}`)
        }

        if (paramSchema.type() !== "string") {
            throw diagnostic(`expected parameter schema type to be string, got ${paramSchema.type()}`)
        }
    })

    return {
        channel,
        address: addr,
        parameters: addrParameters
    }
}

function assertMessagesUseDefaultHeaders(messages: MessagesInterface) {
    messages.all().forEach(message => assertMessageUseDefaultHeaders(message))
}

function assertMessageUseDefaultHeaders(c: MessageInterface) {
    var header = c.headers()
    if (header === undefined) {
        throw diagnostic(`expected header in message ${c.id()}`)
    }

    if (header.type() !== "object") {
        throw diagnostic(`expected header type to be object, but got ${header.type()}`)
    }

    var properties = header.properties()
    if (properties === undefined) {
        throw diagnostic(`expected properties in header, but was undefined`)
    }

    assertHeaderField(properties, "Message-Id", "string")
    assertHeaderField(properties, "Content-Type", "string")
    assertHeaderField(properties, "Spec-Version", "string")
}

function assertHeaderField(header: Record<string, SchemaInterface>, field: string, type: string) {
    var fieldSchema = header[field]
    if (fieldSchema === undefined) {
        throw diagnostic(`header '${field}' is missing`)
    } else if (fieldSchema.type() !== "string") {
        throw diagnostic(`expected header '${field}' to have type ${type}`)
    }
}