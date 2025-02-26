"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.provideYamlSchemas = provideYamlSchemas;
const fs = __importStar(require("fs"));
const path = __importStar(require("path"));
const vscode = __importStar(require("vscode"));
const types_1 = require("./types");
function activate(context, cache, configuration) {
    const yamlExtension = vscode.extensions.getExtension("redhat.vscode-yaml");
    if (yamlExtension) {
        provideYamlSchemas(context, cache, configuration, yamlExtension);
    }
    else {
        // TODO log
    }
}
async function provideYamlSchemas(context, cache, configuration, yamlExtension) {
    if (!yamlExtension.isActive) {
        await yamlExtension.activate();
    }
    let disabled = configuration.get("advanced.disableYamlSchemaContribution");
    configuration.onDidChange((e) => {
        if (configuration.changed(e, "advanced.disableYamlSchemaContribution")) {
            disabled = configuration.get("advanced.disableYamlSchemaContribution");
        }
    });
    function requestSchema(uri) {
        if (disabled) {
            return null;
        }
        for (const document of vscode.workspace.textDocuments) {
            if (document.uri.toString() === uri) {
                const version = cache.getDocumentVersion(document);
                if (version === types_1.OpenApiVersion.V2) {
                    return "openapi:v2";
                }
                else if (version === types_1.OpenApiVersion.V3) {
                    return "openapi:v3";
                }
                break;
            }
        }
        return null;
    }
    function requestSchemaContent(uri) {
        if (!disabled && uri === "openapi:v2") {
            const filename = path.join(context.extensionPath, "schema/generated", "openapi-2.0.json");
            return fs.readFileSync(filename, { encoding: "utf8" });
        }
        else if (!disabled && uri === "openapi:v3") {
            const filename = path.join(context.extensionPath, "schema/generated", "openapi-3.0-2019-04-02.json");
            return fs.readFileSync(filename, { encoding: "utf8" });
        }
        return null;
    }
    const schemaContributor = yamlExtension.exports;
    schemaContributor.registerContributor("openapi", requestSchema, requestSchemaContent);
}
//# sourceMappingURL=yaml-schema-contributor.js.map