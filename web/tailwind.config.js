module.exports = {
  mode: 'jit',
  content: ['./src/**/*.{vue,js,ts}'],
  plugins: [
    require('@tailwindcss/typography'),
    require('daisyui')
  ],
};
