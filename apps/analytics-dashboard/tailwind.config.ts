//import type { Config } from 'tailwindcss';
import famgoPreset from '@famgo/ui-theme/tailwind-preset';
const config: Config = {
  darkMode: 'class',
  content: [
    './src/**/*.{ts,tsx}',
  ],
  presets: [famgoPreset],
  theme: {
    extend: {},
  },
  plugins: [],
};

export default config;
