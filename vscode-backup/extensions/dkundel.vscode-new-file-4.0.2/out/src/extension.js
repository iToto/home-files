"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const Debug = require("debug");
const vscode_1 = require("vscode");
const file_controller_1 = require("./file-controller");
const debug = Debug('vscode-new-file');
function activate(context) {
    debug('Your extension "vscode-new-file" is now active!');
    const disposable = vscode_1.commands.registerCommand('newFile.createNewFile', () => __awaiter(this, void 0, void 0, function* () {
        const File = new file_controller_1.FileController().readSettings();
        try {
            const root = yield File.determineRoot();
            const defaultFileName = yield File.getDefaultFileValue(root);
            const userFilePath = yield File.showFileNameDialog(defaultFileName);
            const createdFiles = yield File.createFiles(userFilePath);
            yield File.openFilesInEditor(createdFiles);
        }
        catch (err) {
            if (err && err.message) {
                vscode_1.window.showErrorMessage(err.message);
            }
        }
    }));
    context.subscriptions.push(disposable);
    const disposableExplorerEntry = vscode_1.commands.registerCommand('newFile.createFromExplorer', (file) => __awaiter(this, void 0, void 0, function* () {
        if (!file || !file.path) {
            return;
        }
        const File = new file_controller_1.FileController().readSettings(file);
        try {
            const root = yield File.getRootFromExplorerPath(file.fsPath);
            const defaultFileName = yield File.getDefaultFileValue(root);
            const userFilePath = yield File.showFileNameDialog(defaultFileName, true);
            const createdFiles = yield File.createFiles(userFilePath);
            yield File.openFilesInEditor(createdFiles);
        }
        catch (err) {
            if (err && err.message) {
                vscode_1.window.showErrorMessage(err.message);
            }
        }
    }));
}
exports.activate = activate;
//# sourceMappingURL=extension.js.map