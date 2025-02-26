"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const fs = require("fs");
function fileExists(path) {
    return new Promise((resolve, reject) => {
        fs.exists(path, exists => {
            resolve(exists);
        });
    });
}
exports.fileExists = fileExists;
//# sourceMappingURL=utils.js.map