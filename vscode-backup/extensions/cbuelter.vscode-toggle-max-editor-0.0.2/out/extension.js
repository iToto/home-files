"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const vscode = require("vscode");
// Function taken from: https://github.com/DavidBabel/clever-vscode/blob/master/src/controlers/maximize.ts
let expanded = false;
function toggleMaximizePane(withSideBar = false) {
    return __awaiter(this, void 0, void 0, function* () {
        if (expanded) {
            if (withSideBar) {
                yield vscode.commands.executeCommand('workbench.action.toggleSidebarVisibility');
            }
            yield vscode.commands.executeCommand('workbench.action.evenEditorWidths');
            yield vscode.commands.executeCommand('workbench.action.focusActiveEditorGroup');
        }
        else {
            yield vscode.commands.executeCommand('workbench.action.maximizeEditor');
            yield vscode.commands.executeCommand('workbench.action.maximizeEditor');
        }
        expanded = !expanded;
    });
}
exports.toggleMaximizePane = toggleMaximizePane;
function activate(context) {
    let disposable = vscode.commands.registerCommand('togglemaxeditor.toggleMaximizeEditor', () => {
        toggleMaximizePane();
    });
    context.subscriptions.push(disposable);
}
exports.activate = activate;
function deactivate() { }
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map