// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
var vscode = require('vscode');
// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
function activate(context) {
    // Use the console to output diagnostic information (console.log) and errors (console.error)
    // This line of code will only be executed once when your extension is activated
    console.log('"vscode-multi-line-tricks" is now active!');
    /**
     * Selects the cursor line and moves the cursor to the next line. If the
     * cursor is at the last line of the document, selects the line to its end.
     */
    context.subscriptions.push(vscode.commands.registerCommand('bcjti.selectLine', function () {
        var textEditor = vscode.window.activeTextEditor;
        var posCur = textEditor.selection.active;
        var startLinePos = posCur.with(posCur.line, 0);
        var endLinePos = startLinePos.with(startLinePos.line + 1, 0);
        if (textEditor.document.lineCount === startLinePos.line + 1) {
            var line = textEditor.document.lineAt(startLinePos.line);
            console.log(line);
            endLinePos = startLinePos.with(startLinePos.line, line.range.end.character);
        }
        if (textEditor.selection.isEmpty) {
            textEditor.selection = new vscode.Selection(startLinePos, endLinePos);
        }
        else {
            var r = textEditor.selection.union(new vscode.Range(startLinePos, endLinePos));
            textEditor.selection = new vscode.Selection(r.start, r.end);
        }
    }));
    /**
     * Breaks the selection into one curso per line selected at the end of the line.
     */
    context.subscriptions.push(vscode.commands.registerCommand('bcjti.breakSelection', function () {
        var textEditor = vscode.window.activeTextEditor;
        var sel = textEditor.selection;
        if (!sel.isEmpty) {
            var doc = textEditor.document;
            var sels = new Array();
            for (var i = sel.start.line; i <= sel.end.line; i++) {
                if (i !== sel.end.line) {
                    var pos = new vscode.Position(i, doc.lineAt(i).range.end.character);
                    sels.push(new vscode.Selection(pos, pos));
                }
                else if (sel.end.character > 0) {
                    sels.push(new vscode.Selection(sel.end, sel.end));
                }
            }
            textEditor.selections = sels;
        }
    }));
}
exports.activate = activate;
//# sourceMappingURL=extension.js.map