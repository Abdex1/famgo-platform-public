const js = require("@eslint/js");
const tseslint = require("typescript-eslint");

module.exports = [
  js.configs.recommended,

  ...tseslint.configs.recommended,

  {
    ignores: ["dist/**", "build/**", ".next/**", "coverage/**", "node_modules/**"],
  },
];
/*
module.exports = {
  root: true,

  env: {
    node: true,
    es2023: true,
  },

  parser: "@typescript-eslint/parser",

  plugins: ["@typescript-eslint"],

  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "prettier",
  ],

  ignorePatterns: [
    "dist",
    "build",
    ".next",
    "coverage",
    "node_modules",
  ],
};
*/
