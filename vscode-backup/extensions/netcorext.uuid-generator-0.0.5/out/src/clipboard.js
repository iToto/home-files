"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.copy = void 0;
const child_process_1 = require("child_process");
let _copy;
switch (process.platform) {
    case "darwin":
        _copy = { command: 'pbcopy', args: [] };
        break;
    case "win32":
        _copy = { command: 'clip', args: [] };
        break;
    case "linux":
        _copy = { command: 'xclip', args: ["-selection", "clipboard"] };
        break;
    default:
        throw new Error(`Unknown platform: '${process.platform}'.`);
}
const copy = function (text, callback) {
    const child = (0, child_process_1.spawn)(_copy.command, _copy.args);
    child.stdin.end(text);
    callback();
};
exports.copy = copy;
//# sourceMappingURL=clipboard.js.map