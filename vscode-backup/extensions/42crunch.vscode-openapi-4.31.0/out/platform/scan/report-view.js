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
exports.ScanReportWebView = void 0;
const vscode = __importStar(require("vscode"));
const audit_1 = require("../../audit/audit");
const service_1 = require("../../audit/service");
const util_1 = require("../../audit/util");
const web_view_1 = require("../../webapps/web-view");
const http_handler_1 = require("../../webapps/http-handler");
const cli_ast_1 = require("../cli-ast");
class ScanReportWebView extends web_view_1.WebView {
    constructor(title, extensionPath, cache, configuration, secrets, store, envStore, prefs, auditView, auditContext) {
        super(extensionPath, "scan", title, vscode.ViewColumn.One, "eye");
        this.cache = cache;
        this.configuration = configuration;
        this.secrets = secrets;
        this.store = store;
        this.envStore = envStore;
        this.prefs = prefs;
        this.auditView = auditView;
        this.auditContext = auditContext;
        this.hostHandlers = {
            sendHttpRequest: ({ id, request, config }) => (0, http_handler_1.executeHttpRequest)(id, request, config),
            sendCurlRequest: async (curl) => {
                return copyCurl(curl);
            },
            savePrefs: async (prefs) => {
                if (this.document) {
                    const uri = this.document.uri.toString();
                    this.prefs[uri] = {
                        ...this.prefs[uri],
                        ...prefs,
                    };
                }
            },
            showEnvWindow: async () => {
                vscode.commands.executeCommand("openapi.showEnvironment");
            },
            showJsonPointer: async (payload) => {
                if (this.document) {
                    let editor = undefined;
                    // check if document is already open
                    for (const visibleEditor of vscode.window.visibleTextEditors) {
                        if (visibleEditor.document.uri.toString() === this.document.uri.toString()) {
                            editor = visibleEditor;
                        }
                    }
                    if (!editor) {
                        editor = await vscode.window.showTextDocument(this.document, vscode.ViewColumn.One);
                    }
                    const root = this.cache.getParsedDocument(editor.document);
                    const lineNo = (0, util_1.getLocationByPointer)(editor.document, root, payload)[0];
                    const textLine = editor.document.lineAt(lineNo);
                    editor.selection = new vscode.Selection(lineNo, 0, lineNo, 0);
                    editor.revealRange(textLine.range, vscode.TextEditorRevealType.AtTop);
                }
            },
            showAuditReport: async () => {
                const uri = this.document.uri.toString();
                const audit = await (0, audit_1.parseAuditReport)(this.cache, this.document, this.auditReport.report, this.auditReport.mapping);
                (0, service_1.setAudit)(this.auditContext, uri, audit);
                await this.auditView.showReport(audit);
            },
        };
        envStore.onEnvironmentDidChange((env) => {
            if (this.isActive()) {
                this.sendRequest({
                    command: "loadEnv",
                    payload: { default: undefined, secrets: undefined, [env.name]: env.environment },
                });
            }
        });
        vscode.window.onDidChangeActiveColorTheme((e) => {
            if (this.isActive()) {
                this.sendColorTheme(e);
            }
        });
    }
    async onStart() {
        await this.sendColorTheme(vscode.window.activeColorTheme);
    }
    async onDispose() {
        this.document = undefined;
        if (this.temporaryReportDirectory !== undefined) {
            await (0, cli_ast_1.cleanupTempScanDirectory)(this.temporaryReportDirectory);
        }
        await super.onDispose();
    }
    async sendStartScan(document) {
        this.document = document;
        this.auditReport = undefined;
        await this.show();
        return this.sendRequest({ command: "startScan", payload: undefined });
    }
    async sendAuditError(document, report, mapping) {
        this.document = document;
        this.auditReport = {
            report,
            mapping,
        };
        return this.sendRequest({
            command: "showGeneralError",
            payload: {
                message: "OpenAPI has failed Security Audit. Please run API Security Audit, fix the issues and try running the Scan again.",
                code: "audit-error",
            },
        });
    }
    setTemporaryReportDirectory(dir) {
        this.temporaryReportDirectory = dir;
    }
    async sendLogMessage(message, level) {
        this.sendRequest({
            command: "showLogMessage",
            payload: { message, level, timestamp: new Date().toISOString() },
        });
    }
    async showGeneralError(error) {
        this.sendRequest({
            command: "showGeneralError",
            payload: error,
        });
    }
    async showScanReport(path, method, report, oas) {
        await this.sendRequest({
            command: "showScanReport",
            // FIXME path and method are ignored by the UI, fix message to make 'em optionals
            payload: {
                path,
                method,
                report: report,
                security: undefined,
                oas,
            },
        });
    }
    async showFullScanReport(report, oas) {
        await this.sendRequest({
            command: "showFullScanReport",
            // FIXME path and method are ignored by the UI, fix message to make 'em optionals
            payload: {
                report: report,
                security: undefined,
                oas,
            },
        });
    }
}
exports.ScanReportWebView = ScanReportWebView;
async function copyCurl(curl) {
    vscode.env.clipboard.writeText(curl);
    const disposable = vscode.window.setStatusBarMessage(`Curl command copied to the clipboard`);
    setTimeout(() => disposable.dispose(), 1000);
}
//# sourceMappingURL=report-view.js.map