module.exports = {
    'env': {
        'browser': true,
        'es2021': true,
    },
    'extends': [
        // "react-app",
        // 'plugin:react/recommended',
        'google',
        "plugin:prettier/recommended",
        // "prettier/@typescript-eslint",
        // "prettier/react"
    ],
    'parserOptions': {
        'ecmaFeatures': {
            'jsx': true,
        },
        'ecmaVersion': 13,
        'sourceType': 'module',
    },
    'plugins': [
        // 'react',
    ],
    'rules': {
        "require-jsdoc": ["error", {
            "require": {
                "FunctionDeclaration": true,
                "MethodDefinition": false,
                "ClassDeclaration": false,
                "ArrowFunctionExpression": false,
                "FunctionExpression": false
            }
        }]
    },
};