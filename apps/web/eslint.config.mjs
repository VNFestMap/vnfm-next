import withNuxt from './.nuxt/eslint.config.mjs'
import a11y from 'eslint-plugin-vuejs-accessibility'

export default withNuxt(
  { ignores: ['app/assets/kun-icons.ts'] },
  a11y.configs['flat/recommended'],
  {
    rules: {
      'no-console': 'off',
      camelcase: 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      'vue/multi-word-component-names': 'off',
      'vue/singleline-html-element-content-newline': 'off',
      'vue/html-self-closing': 'off',
      'vue/attributes-order': 'off',
      'vue/no-multiple-template-root': 'off',
      'vuejs-accessibility/no-autofocus': 'off',
      'vuejs-accessibility/click-events-have-key-events': 'warn',
      'vuejs-accessibility/no-static-element-interactions': 'warn'
    }
  }
)
