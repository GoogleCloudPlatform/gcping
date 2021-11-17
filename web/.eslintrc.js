module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: ["google", "plugin:prettier/recommended"],
  parserOptions: {},
  plugins: [],
  rules: {
    "require-jsdoc": [
      "error",
      {
        require: {
          FunctionDeclaration: true,
          MethodDefinition: false,
          ClassDeclaration: false,
          ArrowFunctionExpression: false,
          FunctionExpression: false,
        },
      },
    ],
  },
};
