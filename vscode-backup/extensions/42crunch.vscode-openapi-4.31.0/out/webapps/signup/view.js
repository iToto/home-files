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
exports.SignUpWebView = void 0;
const vscode = __importStar(require("vscode"));
const web_view_1 = require("../web-view");
const client_1 = require("../../audit/client");
const time_util_1 = require("../../time-util");
class SignUpWebView extends web_view_1.WebView {
    constructor(extensionPath, configuration, secrets, platform, logger) {
        super(extensionPath, "signup", "Sign Up", vscode.ViewColumn.One);
        this.configuration = configuration;
        this.secrets = secrets;
        this.platform = platform;
        this.logger = logger;
        this.hostHandlers = {
            requestAnondTokenByEmail: async (email) => {
                try {
                    await (0, client_1.requestToken)(email);
                    this.sendRequest({
                        command: "showAnondTokenResponse",
                        payload: {
                            success: true,
                        },
                    });
                }
                catch (e) {
                    this.sendRequest({
                        command: "showAnondTokenResponse",
                        payload: {
                            success: false,
                            message: "Unexpected error when trying to request token: " + e,
                        },
                    });
                }
            },
            anondSignUpComplete: async (anondCredentials) => {
                await this.configuration.update("securityAuditToken", anondCredentials.anondToken, vscode.ConfigurationTarget.Global);
                await this.configuration.update("platformAuthType", "anond-token", vscode.ConfigurationTarget.Global);
                await (0, time_util_1.delay)(3000);
                this.close("anond-token");
            },
            platformSignUpComplete: async (platformCredentials) => {
                const credentials = {
                    platformUrl: platformCredentials.platformUrl,
                    apiToken: platformCredentials.platformApiToken,
                    services: "",
                };
                const result = await this.platform.testConnection(credentials);
                if (result.success) {
                    await this.configuration.update("platformUrl", platformCredentials.platformUrl, vscode.ConfigurationTarget.Global);
                    await this.secrets.store("platformApiToken", platformCredentials.platformApiToken);
                    await this.configuration.update("platformAuthType", "api-token", vscode.ConfigurationTarget.Global);
                    await (0, time_util_1.delay)(3000);
                    this.close("api-token");
                }
                else {
                    this.sendRequest({
                        command: "showPlatformConnectionTestError",
                        payload: {
                            error: result.message,
                        },
                    });
                }
            },
            openLink: async (url) => {
                vscode.env.openExternal(vscode.Uri.parse(url));
            },
        };
        vscode.window.onDidChangeActiveColorTheme((e) => {
            if (this.isActive()) {
                this.sendColorTheme(e);
            }
        });
    }
    async onStart() {
        await this.sendColorTheme(vscode.window.activeColorTheme);
    }
    async showSignUp(resolve) {
        this.resolve = resolve;
        await this.show();
    }
    close(value) {
        if (this.resolve) {
            this.resolve(value);
        }
        this.resolve = undefined;
        // Event onDispose will be fired
        this.dispose();
    }
    async onDispose() {
        if (this.resolve) {
            // User closed this panel manually
            this.resolve(undefined);
        }
        this.resolve = undefined;
        await super.onDispose();
    }
}
exports.SignUpWebView = SignUpWebView;
//# sourceMappingURL=view.js.map