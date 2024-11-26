module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: [
    "airbnb-base",
    "plugin:vue/vue3-recommended",
  ],
  parserOptions: {
    ecmaVersion: 12,
    sourceType: "module",
  },
  plugins: ["vue"],
  rules: {
    "no-console": "off", 
    "no-debugger": "off", 
    "vue/multi-word-component-names": "off",
    "no-trailing-spaces": "off",
    quotes: "off",
  },
  globals: {
    google: "readonly", 
  },
};