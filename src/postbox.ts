import { TemplateResult, catchDiagnostics, template } from "./templates"

export type PostboxOptions = {
    package?: string,
    version?: string
}

export type PostboxData = {
    package: string
    version: string
}

export default function postboxTemplate(opts?: PostboxOptions): TemplateResult<PostboxData> {
    try {
        var data: PostboxData = {
            package: opts?.package ?? "asyncapi",
            version: opts?.version ?? "v0.0.0",
        }

        return template("./postbox", data)
    } catch(e) {
        return catchDiagnostics(e)
    }
}