"use strict";
/*
 Copyright (c) 42Crunch Ltd. All rights reserved.
 Licensed under the GNU Affero General Public License version 3. See LICENSE.txt in the project root for license information.
*/
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
exports.registerSecurityAudit = registerSecurityAudit;
exports.registerSingleOperationAudit = registerSingleOperationAudit;
exports.registerOutlineSingleOperationAudit = registerOutlineSingleOperationAudit;
exports.registerFocusSecurityAudit = registerFocusSecurityAudit;
exports.registerFocusSecurityAuditById = registerFocusSecurityAuditById;
const vscode = __importStar(require("vscode"));
const preserving_json_yaml_parser_1 = require("@xliic/preserving-json-yaml-parser");
const configuration_1 = require("../configuration");
const credentials_1 = require("../credentials");
const util_1 = require("../outlines/util");
const cli_ast_1 = require("../platform/cli-ast");
const extract_1 = require("../util/extract");
const decoration_1 = require("./decoration");
const cli_1 = require("./runtime/cli");
const platform_1 = require("./runtime/platform");
const service_1 = require("./service");
const config_1 = require("../util/config");
function registerSecurityAudit(context, cache, auditContext, pendingAudits, reportWebView, store, signUpWebView) {
    return vscode.commands.registerTextEditorCommand("openapi.securityAudit", async (textEditor, edit) => {
        await securityAudit(signUpWebView, context.workspaceState, context.secrets, cache, auditContext, pendingAudits, reportWebView, store, textEditor);
    });
}
function registerSingleOperationAudit(context, cache, auditContext, pendingAudits, reportWebView, store, signUpWebView) {
    return vscode.commands.registerTextEditorCommand("openapi.editorSingleOperationAudit", async (textEditor, edit, path, method) => {
        await securityAudit(signUpWebView, context.workspaceState, context.secrets, cache, auditContext, pendingAudits, reportWebView, store, textEditor, path, method);
    });
}
function registerOutlineSingleOperationAudit(context, cache, auditContext, pendingAudits, reportWebView, store, signUpWebView) {
    return vscode.commands.registerCommand("openapi.outlineSingleOperationAudit", async (node) => {
        if (!vscode.window.activeTextEditor) {
            vscode.window.showErrorMessage("No active editor");
            return;
        }
        const { path, method } = (0, util_1.getPathAndMethod)(node);
        await securityAudit(signUpWebView, context.workspaceState, context.secrets, cache, auditContext, pendingAudits, reportWebView, store, vscode.window.activeTextEditor, path, method);
    });
}
function registerFocusSecurityAudit(context, cache, auditContext, reportWebView) {
    return vscode.commands.registerCommand("openapi.focusSecurityAudit", async (documentUri) => {
        try {
            const audit = auditContext.auditsByMainDocument[documentUri];
            if (audit) {
                reportWebView.showReport(audit);
            }
        }
        catch (e) {
            vscode.window.showErrorMessage(`Unexpected error: ${e}`);
        }
    });
}
function registerFocusSecurityAuditById(context, auditContext, reportWebView) {
    return vscode.commands.registerTextEditorCommand("openapi.focusSecurityAuditById", async (textEditor, edit, params) => {
        try {
            const documentUri = textEditor.document.uri.toString();
            const uri = Buffer.from(params.uri, "base64").toString("utf8");
            const audit = auditContext.auditsByMainDocument[uri];
            if (audit && audit.issues[documentUri]) {
                reportWebView.showIds(audit, documentUri, params.ids);
            }
        }
        catch (e) {
            vscode.window.showErrorMessage(`Unexpected error: ${e}`);
        }
    });
}
async function securityAudit(signUpWebView, memento, secrets, cache, auditContext, pendingAudits, reportWebView, store, editor, path, method) {
    if (!(await (0, credentials_1.ensureHasCredentials)(signUpWebView, configuration_1.configuration, secrets))) {
        return;
    }
    if (!(await offerDataDictionaryUpdateAndContinue(editor.document.uri))) {
        return;
    }
    const uri = editor.document.uri.toString();
    if (pendingAudits[uri]) {
        vscode.window.showErrorMessage(`Audit for "${uri}" is already in progress`);
        return;
    }
    delete auditContext.auditsByMainDocument[uri];
    pendingAudits[uri] = true;
    try {
        reportWebView.prefetchKdb();
        await reportWebView.sendStartAudit();
        const audit = await vscode.window.withProgress({
            location: vscode.ProgressLocation.Notification,
            title: "Running API Security Audit...",
            cancellable: false,
        }, async (progress, cancellationToken) => {
            const isFullAudit = path === undefined || method === undefined;
            const { value, mapping } = await bundleOrThrow(cache, editor.document);
            const oas = isFullAudit
                ? (0, preserving_json_yaml_parser_1.stringify)(value)
                : (0, preserving_json_yaml_parser_1.stringify)((0, extract_1.extractSingleOperation)(method, path, value));
            if ((await chooseAuditRuntime(configuration_1.configuration, secrets)) === "platform") {
                return (0, platform_1.runPlatformAudit)(editor.document, oas, mapping, cache, store, memento);
            }
            else {
                // use CLI
                if (await (0, cli_ast_1.ensureCliDownloaded)(configuration_1.configuration, secrets)) {
                    const tags = store.isConnected()
                        ? await store.getTagsForDocument(editor.document, memento)
                        : [];
                    return (0, cli_1.runCliAudit)(editor.document, oas, mapping, tags, cache, secrets, configuration_1.configuration, progress, isFullAudit);
                }
                else {
                    // cli is not available and user chose to cancel download
                    vscode.window.showErrorMessage("42Crunch API Security Testing Binary is required to run Audit.");
                    return;
                }
            }
        });
        if (audit) {
            (0, service_1.setAudit)(auditContext, uri, audit);
            (0, decoration_1.setDecorations)(editor, auditContext);
            await reportWebView.showReport(audit);
        }
        else {
            await reportWebView.sendCancelAudit();
        }
        delete pendingAudits[uri];
    }
    catch (e) {
        delete pendingAudits[uri];
        vscode.window.showErrorMessage(`Failed to audit: ${e}`);
    }
}
async function bundleOrThrow(cache, document) {
    const bundle = await cache.getDocumentBundle(document, { rebundle: true });
    if (!bundle || "errors" in bundle) {
        vscode.commands.executeCommand("workbench.action.problems.focus");
        throw new Error("Failed to bundle for audit, check OpenAPI file for errors.");
    }
    return bundle;
}
async function offerDataDictionaryUpdateAndContinue(documentUri) {
    const proceed = await vscode.commands.executeCommand("openapi.platform.dataDictionaryPreAuditBulkUpdateProperties", documentUri);
    return proceed === true;
}
async function chooseAuditRuntime(configuration, secrets) {
    const config = await (0, config_1.loadConfig)(configuration, secrets);
    // paid users are allowed to choose the runtime, freemium users always use the cli
    if (config.platformAuthType === "api-token") {
        return config.auditRuntime;
    }
    else {
        return "cli";
    }
}
//# sourceMappingURL=commands.js.map