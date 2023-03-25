"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const react_1 = __importDefault(require("react"));
const Dialog_1 = __importDefault(require("@material-ui/core/Dialog"));
const DialogActions_1 = __importDefault(require("@material-ui/core/DialogActions"));
const DialogContent_1 = __importDefault(require("@material-ui/core/DialogContent"));
const DialogContentText_1 = __importDefault(require("@material-ui/core/DialogContentText"));
const DialogTitle_1 = __importDefault(require("@material-ui/core/DialogTitle"));
const Button_1 = __importDefault(require("@material-ui/core/Button"));
const AlertDialog = ({ close, message, options }) => {
    return (react_1.default.createElement(Dialog_1.default, { open: true, onClose: close, keepMounted: true, "aria-labelledby": "alert-dialog-slide-title", "aria-describedby": "alert-dialog-slide-description" },
        react_1.default.createElement(DialogTitle_1.default, { id: "alert-dialog-slide-title" }, options.title),
        react_1.default.createElement(DialogContent_1.default, null,
            react_1.default.createElement(DialogContentText_1.default, { id: "alert-dialog-slide-description" }, message)),
        react_1.default.createElement(DialogActions_1.default, null,
            options.actions &&
                options.actions.map((action, index) => (react_1.default.createElement(Button_1.default, { onClick: () => {
                        action.onClick();
                        close();
                    }, color: "primary", key: index }, action.copy))),
            react_1.default.createElement(Button_1.default, { onClick: close, color: "primary" }, options.closeCopy || "Okay"))));
};
exports.default = AlertDialog;
//# sourceMappingURL=index.js.map