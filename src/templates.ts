export type Template<T> = {
    source: string,
    data: T
}

export type TemplateDiagnostic = {
    message: string
    stack?: string
}

export type TemplateResult<T> = {
    template: Template<T>
    diagnostics: undefined
} | {
    template: undefined
    diagnostics: TemplateDiagnostic[]
}

export function template<T>(source: string, data: T): TemplateResult<T> {
    return {
        template: {
            source,
            data
        },
        diagnostics: undefined
    }
}

export function diagnostic(message: string): TemplateDiagnostic {
    return {
        message,
        stack: new Error().stack
    }
}

export function fail<T>(message: string): TemplateResult<T> {
    return {
        template: undefined,
        diagnostics: [diagnostic(message)]
    }
}

export function error<T>(e: Error): TemplateResult<T> {
    return {
        template: undefined,
        diagnostics: [{
            message: e.message,
            stack: e.stack
        }]
    }
}

export function catchDiagnostics<T>(e: any): TemplateResult<T> {
    if (e instanceof Error) {
        return error(e)
    } else if (e instanceof Array) {
        return {
            template: undefined,
            diagnostics: e
        }
    } else {
        return {
            template: undefined,
            diagnostics: [e]
        }
    }
}