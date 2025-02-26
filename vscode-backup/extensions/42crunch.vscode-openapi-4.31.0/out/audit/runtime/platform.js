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
exports.runPlatformAudit = runPlatformAudit;
/*
 Copyright (c) 42Crunch Ltd. All rights reserved.
 Licensed under the GNU Affero General Public License version 3. See LICENSE.txt in the project root for license information.
*/
const vscode = __importStar(require("vscode"));
const platform_store_1 = require("../../platform/stores/platform-store");
const audit_1 = require("../audit");
const util_1 = require("../../platform/util");
async function runPlatformAudit(document, oas, mapping, cache, store, memento) {
    try {
        const tmpApi = await store.createTempApi(oas, (0, platform_store_1.getTagDataEntry)(memento, document.uri.fsPath));
        const report = await store.getAuditReport(tmpApi.apiId);
        const compliance = await store.readAuditCompliance(report.tid);
        const todoReport = await store.readAuditReportSqgTodo(report.tid);
        await store.clearTempApi(tmpApi);
        const audit = await (0, audit_1.parseAuditReport)(cache, document, report.data, mapping);
        const { issues: todo } = await (0, audit_1.parseAuditReport)(cache, document, todoReport.data, mapping);
        audit.compliance = compliance;
        audit.todo = todo;
        return audit;
    }
    catch (ex) {
        if (ex?.response?.statusCode === 409 &&
            ex?.response?.body?.code === 109 &&
            ex?.response?.body?.message === "limit reached") {
            vscode.window.showErrorMessage("You have reached your maximum number of APIs. Please contact support@42crunch.com to upgrade your account.");
        }
        else {
            vscode.window.showErrorMessage((0, util_1.formatException)("Unexpected error when trying to audit API using the platform:", ex));
        }
    }
}
//# sourceMappingURL=platform.js.map