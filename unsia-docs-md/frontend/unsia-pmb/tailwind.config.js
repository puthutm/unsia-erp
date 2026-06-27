/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#FEF7E6',
          100: '#FDEFC4',
          200: '#FBDF89',
          300: '#F9CF4D',
          400: '#F7BF12',
          500: '#FED524', // Main brand yellow
          600: '#D9B81E',
          700: '#A68918',
          800: '#735B12',
          900: '#402D0C',
        },
        primary: {
          50: '#E6F0F7',
          100: '#CCE0EF',
          200: '#99C1DF',
          300: '#66A2CF',
          400: '#3383BF',
          500: '#0F487B', // Main primary blue
          600: '#0D3D6B',
          700: '#0A3052',
          800: '#072339',
          900: '#041620',
        },
        surface: {
          subtle: '#F8FAFC',
          border: '#E2E8F0',
          hover: '#F1F5F9',
        },
        dark: {
          50: '#F8FAFC',
          100: '#F1F5F9',
          200: '#E2E8F0',
          300: '#CBD5E1',
          400: '#94A3B8',
          500: '#64748B',
          600: '#475569',
          700: '#334155',
          800: '#1E293B',
          900: '#0F172A',
        },
      },
      fontFamily: {
        display: ['var(--font-display)', 'system-ui', 'sans-serif'],
        body: ['var(--font-body)', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        'soft': '0 2px 8px -2px rgba(0, 0, 0, 0.05), 0 4px 16px -4px rgba(0, 0, 0, 0.05)',
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(8px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
  plugins: [],
}
