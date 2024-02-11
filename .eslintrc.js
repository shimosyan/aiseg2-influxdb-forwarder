const prettierRule = require('./prettier.config');

const tsFileRules = {
  '@typescript-eslint/no-explicit-any': 'off',
  '@typescript-eslint/explicit-module-boundary-types': 'off',
  complexity: ['error', 10],
};

const jsFileRules = {
  'prettier/prettier': ['warn', prettierRule],
  'jest/expect-expect': ['warn', { assertFunctionNames: ['expect'] }],
};

module.exports = {
  extends: ['eslint:recommended', 'plugin:prettier/recommended', 'plugin:jest/recommended'],
  plugins: ['@typescript-eslint', 'jest'],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 2021,
    sourceType: 'module',
  },
  env: {
    browser: true,
    node: true,
    es6: true,
  },
  rules: jsFileRules,
  overrides: [
    {
      files: ['*.ts', '*.tsx'],
      extends: ['plugin:@typescript-eslint/recommended'],
      rules: tsFileRules,
    },
  ],
};
