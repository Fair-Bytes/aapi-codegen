#!/usr/bin/env node

import { Parser, fromFile } from "@asyncapi/parser"
import { Eta } from "eta"
import path from "path"
import operationsTemplate from "./operations"
import modelsTemplate from "./models"
import { writeFileSync } from "fs"
import { Command } from 'commander';
import { exec } from 'child_process';

const program = new Command();

const parser = new Parser();
const eta = new Eta({ 
    views: path.join(__dirname, "templates"),
    autoTrim: false,
    autoEscape: false
});

type Options = {
    spec: string,
    package: string,
    models: string | undefined,
    operations: string | undefined
}

function proccessArgs(): Options {
    program
        .requiredOption('-p, --package <package>', 'golang package')
        .requiredOption('-g, --generate <generate>', 'what to generate', 'models,operations')
        .option('-m, --models <models>', 'models file')
        .option('-o, --operations <operations>', 'operations file')
        .argument('<asyncapi>')
        .parse()

    var opts: Options = {
        spec: program.args[0],
        package: program.opts().package as string,
        operations: program.opts().generate.split(',').includes('operations') ? program.opts().operations as string : undefined,
        models: program.opts().generate.split(',').includes('models') ? program.opts().models as string : undefined
    }

    if (program.args.length !== 1) {
        program.error(`only one asyncapi file expected, got ${program.args.length}`)
    }

    if (program.opts().generate.split(',').includes('operations') && opts.operations === undefined) {
        program.error("expected '-o, --operations' option")
    }

    if (program.opts().generate.split(',').includes('models') && opts.models === undefined) {
        program.error("expected '-m, --models' option")
    }

    return opts
}

(async () => {
    var opts = proccessArgs()

    var {document, diagnostics} = await fromFile(parser, opts.spec).parse()

    if (document == undefined) {
        console.error(diagnostics)
        process.exit(1)
    }

    // get version
    var version = document.info().version()
    if (version.charAt(0) !== "v") {
        version = "v" + version
    }

    // Calculate templates
    if (!!opts.models) {
        var models = modelsTemplate(document.allOperations(), { package: opts.package, version })
        if (models.diagnostics !== undefined) {
            console.error(models.diagnostics)
            process.exit(1)
        }
        
        var modelContent = eta.render(models.template.source, models.template.data)
        writeFileSync(opts.models, modelContent)

        // Post processing
        exec(`gofmt -w ${opts.models}`, (error, stdout, stderr) => {
            if (error != null) {
                console.error(`gofmt failed (${error.code}): ${error.message}`)
                console.error(stderr)
            }
        })
    }
    
    if (!!opts.operations) {
        var operations = operationsTemplate(document.allOperations(), { package: opts.package })
        if (operations.diagnostics !== undefined) {
            console.error(operations.diagnostics)
            process.exit(1)
        }
        
        var operationsContent = eta.render(operations.template.source, operations.template.data)
        writeFileSync(opts.operations, operationsContent)

        // Post processing
        exec(`gofmt -w ${opts.operations}`, (error, stdout, stderr) => {
            if (error != null) {
                console.error(`gofmt failed (${error.code}): ${error.message}`)
                console.error(stderr)
            }
        })
    }
})()