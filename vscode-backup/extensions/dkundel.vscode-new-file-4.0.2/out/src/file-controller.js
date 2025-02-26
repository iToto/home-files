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
const vscode_1 = require("vscode");
const braces = require("braces");
const Debug = require("debug");
const fs = require("fs");
const mkdirp = require("mkdirp");
const path = require("path");
const q_1 = require("q");
const utils_1 = require("./utils");
const debug = Debug('vscode-new-file');
const mkdir = q_1.denodeify(mkdirp);
const readFile = q_1.denodeify(fs.readFile);
const appendFile = q_1.denodeify(fs.appendFile);
const fsStat = q_1.denodeify(fs.stat);
class FileController {
    readSettings(currentUri) {
        this.currentUri = currentUri || this.getUriOfCurrentFile();
        const config = vscode_1.workspace.getConfiguration('newFile', this.currentUri);
        this.settings = {
            defaultBaseFileName: config.get('defaultBaseFileName', 'newFile'),
            defaultFileExtension: config.get('defaultFileExtension', '.ts'),
            expandBraces: config.get('expandBraces', false),
            fileTemplates: config.get('fileTemplates', {}),
            relativeTo: config.get('relativeTo', 'file'),
            rootDirectory: config.get('rootDirectory', this.homedir()),
            showPathRelativeTo: config.get('showPathRelativeTo', 'root'),
            useFileTemplates: config.get('useFileTemplates', true),
        };
        const showFullPath = config.get('showFullPath');
        if (showFullPath !== undefined) {
            vscode_1.window.showInformationMessage('You are using a deprecated option "showFullPath". Switch instead to "showPathRelativeTo"');
            this.settings.showPathRelativeTo = 'root';
        }
        return this;
    }
    getRootFromExplorerPath(filePath) {
        return __awaiter(this, void 0, void 0, function* () {
            let dir = path.dirname(filePath);
            const stats = (yield fsStat(dir));
            if (!stats.isDirectory()) {
                dir = path.resolve(dir, '..');
            }
            this.rootPath = dir;
            return dir;
        });
    }
    determineRoot() {
        return __awaiter(this, void 0, void 0, function* () {
            let root;
            if (this.settings.relativeTo === 'project') {
                this.workspaceRoot = yield this.determineWorkspaceRoot();
                root = this.workspaceRoot;
            }
            else if (this.settings.relativeTo === 'file') {
                if (vscode_1.window.activeTextEditor) {
                    root = path.dirname(vscode_1.window.activeTextEditor.document.fileName);
                }
                else {
                    this.workspaceRoot = yield this.determineWorkspaceRoot();
                    root = this.workspaceRoot;
                }
            }
            if (!root) {
                root = this.settings.rootDirectory;
                if (root.indexOf('~') === 0) {
                    root = path.join(this.homedir(), root.substr(1));
                }
            }
            this.rootPath = root;
            return root;
        });
    }
    getDefaultFileValue(root) {
        return __awaiter(this, void 0, void 0, function* () {
            const newFileName = this.settings.defaultBaseFileName;
            const defaultExtension = this.settings.defaultFileExtension;
            const currentFileName = vscode_1.window.activeTextEditor
                ? vscode_1.window.activeTextEditor.document.fileName
                : '';
            const ext = path.extname(currentFileName) || defaultExtension;
            if (this.settings.showPathRelativeTo !== 'none') {
                const fullPath = path.join(root, `${newFileName}${ext}`);
                if (this.settings.showPathRelativeTo === 'project') {
                    if (!this.workspaceRoot) {
                        this.workspaceRoot = yield this.determineWorkspaceRoot();
                    }
                    return fullPath.replace(this.workspaceRoot + path.sep, '');
                }
                return fullPath;
            }
            else {
                return `${newFileName}${ext}`;
            }
        });
    }
    showFileNameDialog(defaultFileValue, fromExplorer = false) {
        return __awaiter(this, void 0, void 0, function* () {
            let question = `What's the path and name of the new file?`;
            if (fromExplorer) {
                question += ' (Relative to selected file)';
            }
            else if (this.settings.showPathRelativeTo === 'none') {
                if (this.settings.relativeTo === 'project') {
                    question += ' (Relative to project root)';
                }
                else if (this.settings.relativeTo === 'file') {
                    question += ' (Relative to current file)';
                }
            }
            else if (this.settings.showPathRelativeTo === 'project') {
                question += ' (Relative to project root)';
            }
            const selectionBoundsForFileName = [
                defaultFileValue.lastIndexOf(path.sep) + 1,
                defaultFileValue.lastIndexOf('.'),
            ];
            let selectedFilePath = yield vscode_1.window.showInputBox({
                prompt: question,
                value: defaultFileValue,
                valueSelection: selectionBoundsForFileName,
            });
            if (selectedFilePath === null || typeof selectedFilePath === 'undefined') {
                throw undefined;
            }
            selectedFilePath = selectedFilePath || defaultFileValue;
            if (selectedFilePath) {
                if (selectedFilePath.startsWith('./')) {
                    return this.normalizeDotPath(selectedFilePath);
                }
                else {
                    if (this.settings.showPathRelativeTo !== 'none') {
                        if (this.settings.showPathRelativeTo === 'project') {
                            if (!this.workspaceRoot) {
                                this.workspaceRoot = yield this.determineWorkspaceRoot();
                            }
                            selectedFilePath = path.resolve(this.workspaceRoot, selectedFilePath);
                        }
                        return selectedFilePath;
                    }
                    else {
                        return this.getFullPath(this.rootPath, selectedFilePath);
                    }
                }
            }
        });
    }
    createFiles(userEntry) {
        return __awaiter(this, void 0, void 0, function* () {
            if (!this.settings.expandBraces) {
                return Promise.all([this.createFile(userEntry)]);
            }
            const newFileNames = braces.expand(userEntry);
            const fileCreationPromises = newFileNames.map(fileName => this.createFile(fileName));
            return Promise.all(fileCreationPromises);
        });
    }
    createFile(newFileName) {
        return __awaiter(this, void 0, void 0, function* () {
            const dirname = path.dirname(newFileName);
            const extension = path.extname(newFileName);
            const doesFileExist = yield utils_1.fileExists(newFileName);
            if (!doesFileExist) {
                yield mkdir(dirname);
                let content = '';
                const templatePath = this.settings.fileTemplates[extension];
                if (this.settings.useFileTemplates && templatePath !== undefined) {
                    content = (yield readFile(path.resolve(this.settings.rootDirectory, templatePath), 'utf8'));
                }
                yield appendFile(newFileName, content);
            }
            return newFileName;
        });
    }
    openFilesInEditor(fileNames) {
        return fileNames.map((fileName) => __awaiter(this, void 0, void 0, function* () {
            const stats = (yield fsStat(fileName));
            if (stats.isDirectory()) {
                vscode_1.window.showInformationMessage('This file is already a directory. Try a different name.');
                return;
            }
            const textDocument = yield vscode_1.workspace.openTextDocument(fileName);
            if (!textDocument) {
                throw new Error('Could not open file!');
            }
            const editor = vscode_1.window.showTextDocument(textDocument);
            if (!editor) {
                throw new Error('Could not show document!');
            }
            return editor;
        }));
    }
    normalizeDotPath(filePath) {
        return __awaiter(this, void 0, void 0, function* () {
            const currentFileName = vscode_1.window.activeTextEditor
                ? vscode_1.window.activeTextEditor.document.fileName
                : '';
            if (!this.workspaceRoot) {
                this.workspaceRoot = yield this.determineWorkspaceRoot();
            }
            const directory = currentFileName.length > 0
                ? path.dirname(currentFileName)
                : this.workspaceRoot;
            return path.resolve(directory, filePath);
        });
    }
    getFullPath(root, filePath) {
        if (filePath.indexOf('/') === 0) {
            return filePath;
        }
        if (filePath.indexOf('~') === 0) {
            return path.join(this.homedir(), filePath.substr(1));
        }
        return path.resolve(root, filePath);
    }
    getUriOfCurrentFile() {
        const editor = vscode_1.window.activeTextEditor;
        return editor ? editor.document.uri : undefined;
    }
    determineWorkspaceFolder(currentUri) {
        return __awaiter(this, void 0, void 0, function* () {
            if (currentUri) {
                return vscode_1.workspace.getWorkspaceFolder(currentUri);
            }
            const selectedWorkspaceFolder = yield vscode_1.window.showWorkspaceFolderPick();
            if (selectedWorkspaceFolder !== undefined) {
                this.readSettings(selectedWorkspaceFolder.uri);
            }
            return selectedWorkspaceFolder;
        });
    }
    getRootPathFromWorkspace(currentWorkspace) {
        if (typeof currentWorkspace === 'undefined') {
            return undefined;
        }
        if (currentWorkspace.uri.scheme !== 'file') {
            return null;
        }
        return currentWorkspace.uri.fsPath;
    }
    determineWorkspaceRoot() {
        return __awaiter(this, void 0, void 0, function* () {
            const currentWorkspace = yield this.determineWorkspaceFolder(this.currentUri);
            const workspaceRoot = this.getRootPathFromWorkspace(currentWorkspace);
            if (workspaceRoot === null) {
                throw new Error('This extension currently only support file system workspaces.');
            }
            return workspaceRoot;
        });
    }
    homedir() {
        return process.env[process.platform === 'win32' ? 'USERPROFILE' : 'HOME'];
    }
}
exports.FileController = FileController;
//# sourceMappingURL=file-controller.js.map