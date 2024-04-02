import { MessageInterface } from "@asyncapi/parser"

export function stringToUpperCaseStart(s: string): string {
    if (s.length === 0) {
        return ""
        //throw new Error("expected a string with at least one charracter")
    }

    return s.at(0)!.toUpperCase() + s.slice(1)
}

export function toMixedCaps(s: string): string {
    return s
        .split(/[^a-zA-z0-9]{1,}/)
        .map(word => stringToUpperCaseStart(word))
        .join('')
}

export function channelAddressParameters(addr: string): string[] {
    var parameters: string[] = []
    for (var r of addr.matchAll(/{(?<parameter>[^{}]*)}/g)) {
        if (r.groups !== undefined) {
            parameters.push(r.groups["parameter"])
        }
    }
    return parameters
}

export function messageId(message: MessageInterface): string {
    return message.name() ?? message.id()
}
