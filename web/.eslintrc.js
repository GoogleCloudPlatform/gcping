module.exports = {
    'env': {
        'browser': true,
        'es2021': true,
    },
    'extends': [
        'google',
        "plugin:prettier/recommended",
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