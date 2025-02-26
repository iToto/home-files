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
exports.ConfigWebView = void 0;
const vscode = __importStar(require("vscode"));
const http2 = __importStar(require("http2"));
const web_view_1 = require("../../web-view");
const scandManagerApi = __importStar(require("../../../platform/api-scand-manager"));
const config_1 = require("../../../util/config");
const cli_ast_1 = require("../../../platform/cli-ast");
const utils_gen_1 = require("./utils-gen");
const cli_ast_update_1 = require("../../../platform/cli-ast-update");
const http_handler_1 = require("../../http-handler");
class ConfigWebView extends web_view_1.WebView {
    constructor(extensionPath, configuration, secrets, platform, logger) {
        super(extensionPath, "config", "Settings", vscode.ViewColumn.One);
        this.configuration = configuration;
        this.secrets = secrets;
        this.platform = platform;
        this.logger = logger;
        this.hostHandlers = {
            saveConfig: async (config) => {
                await (0, config_1.saveConfig)(config, this.configuration, this.secrets);
                this.config = await (0, config_1.loadConfig)(this.configuration, this.secrets);
                return {
                    command: "loadConfig",
                    payload: this.config,
                };
            },
            testOverlordConnection: async () => {
                const services = this.config?.platformServices.source === "auto"
                    ? this.config?.platformServices.auto
                    : this.config?.platformServices.manual;
                if (services === undefined || services === "") {
                    return {
                        command: "showOverlordConnectionTest",
                        payload: { success: false, message: "Services host is not configured" },
                    };
                }
                const result = await http2Ping(`https://${services}`);
                return {
                    command: "showOverlordConnectionTest",
                    payload: result,
                };
            },
            testPlatformConnection: async () => {
                if (this.config === undefined) {
                    return {
                        command: "showPlatformConnectionTest",
                        payload: { success: false, message: "no credentials" },
                    };
                }
                const credentials = {
                    platformUrl: this.config.platformUrl,
                    apiToken: this.config.platformApiToken,
                    services: "",
                };
                const result = await this.platform.testConnection(credentials);
                return { command: "showPlatformConnectionTest", payload: result };
            },
            testScandManagerConnection: async () => {
                const scandManager = this.config?.scandManager;
                if (scandManager === undefined || scandManager.url === "") {
                    return {
                        command: "showScandManagerConnectionTest",
                        payload: { success: false, message: "no scand manager confguration" },
                    };
                }
                const result = await scandManagerApi.testConnection(scandManager, this.logger);
                return {
                    command: "showScandManagerConnectionTest",
                    payload: result,
                };
            },
            testCli: async () => {
                const result = await (0, cli_ast_1.testCli)(this.config.cliDirectoryOverride);
                // if the binary was found, check for updates
                // otherwise the download button will be shown in the web UI
                if (result.success) {
                    (0, cli_ast_1.checkForCliUpdate)(this.config.repository, this.config.cliDirectoryOverride);
                }
                return {
                    command: "showCliTest",
                    payload: result,
                };
            },
            downloadCli: () => downloadCliHandler(this.config.repository, this.config.cliDirectoryOverride),
            openLink: async (url) => {
                // @ts-ignore
                // workaround for vscode https://github.com/microsoft/vscode/issues/85930
                vscode.env.openExternal(url);
            },
            sendHttpRequest: ({ id, request, config }) => (0, http_handler_1.executeHttpRequest)(id, request, config),
        };
        vscode.window.onDidChangeActiveColorTheme((e) => {
            if (this.isActive()) {
                this.sendColorTheme(e);
            }
        });
    }
    async onStart() {
        await this.sendColorTheme(vscode.window.activeColorTheme);
        this.config = await (0, config_1.loadConfig)(this.configuration, this.secrets);
        if (this.platform.isConnected()) {
            try {
                // this could throw if the token has become invalid since the start
                const convention = await this.platform.getCollectionNamingConvention();
                if (convention.pattern !== "") {
                    this.config.platformCollectionNamingConvention = convention;
                }
            }
            catch (ex) {
                // can't get naming convention if the token is invalid
            }
        }
        await this.sendRequest({
            command: "loadConfig",
            payload: this.config,
        });
    }
    async showConfig() {
        await this.show();
    }
}
exports.ConfigWebView = ConfigWebView;
async function* downloadCliHandler(repository, cliDirectoryOverride) {
    try {
        if (repository === undefined || repository === "") {
            throw new Error("Repository URL is not set");
        }
        const manifest = await (0, cli_ast_update_1.getCliUpdate)(repository, "0.0.0");
        if (manifest === undefined) {
            throw new Error("Failed to download 42Crunch API Security Testing Binary, manifest not found");
        }
        const location = yield* (0, utils_gen_1.transformValues)((0, cli_ast_1.downloadCli)(manifest, cliDirectoryOverride), (progress) => ({
            command: "showCliDownload",
            payload: { completed: false, progress },
        }));
        yield {
            command: "showCliDownload",
            payload: {
                completed: true,
                success: true,
                location,
            },
        };
    }
    catch (ex) {
        yield {
            command: "showCliDownload",
            payload: {
                completed: true,
                success: false,
                error: `Failed to download: ${ex}`,
            },
        };
    }
}
function http2Ping(url) {
    const timeout = 5000;
    return new Promise((resolve, reject) => {
        try {
            const client = http2.connect(url);
            client.setTimeout(timeout);
            client.on("error", (err) => {
                client.close();
                resolve({
                    success: false,
                    message: err.message,
                });
            });
            client.on("timeout", (err) => {
                client.close();
                resolve({
                    success: false,
                    message: `Timed out wating to connect after ${timeout}ms`,
                });
            });
            client.on("connect", () => {
                client.close();
                resolve({
                    success: true,
                });
            });
        }
        catch (ex) {
            resolve({
                success: false,
                message: `Failed to create connection: ${ex}`,
            });
        }
    });
}
//# sourceMappingURL=view.js.map