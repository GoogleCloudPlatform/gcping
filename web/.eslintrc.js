module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: ["eslint:recommended", "google", "plugin:prettier/recommended"],
  plugins: [],
  rules: {
    "no-implicit-globals": "error",
  },
  parserOptions: {
    ecmaVersion: 11,
    sourceType: "module",
  },
};
