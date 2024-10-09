module.exports = {
  'src/**/*.(ts|tsx|js|jsx)': [
    'prettier --list-different',
    'eslint',
    () =>
      'tsc --project tsconfig.json --noEmit --skipLibCheck --isolatedModules false',
  ],
}
