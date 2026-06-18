/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // 主色调 - Vermillion / burnt-orange (Factory.ai signature). Anchor 500 ≈ #ef6f2e
        primary: {
          50: '#fdf4ef',
          100: '#f9e1d3',
          200: '#f3c3a6',
          300: '#ec9e72',
          400: '#e67f47',
          500: '#ef6f2e',
          600: '#d8581a',
          700: '#b3430c',
          800: '#8f3710',
          900: '#743012',
          950: '#411706'
        },
        // 辅助色 - repoint to warm neutral (same scale as gray). Used in .text-gradient, sidebar, etc.
        accent: {
          50: '#f7f5f3',
          100: '#ece7e1',
          200: '#d9d1c7',
          300: '#bfb3a5',
          400: '#a09486',
          500: '#7d7264',
          600: '#5e554a',
          700: '#4a4339',
          800: '#2b2722',
          900: '#1a1714',
          950: '#0c0a08'
        },
        // OVERRIDE default cool gray → warm neutral (warms the whole UI at once)
        gray: {
          50: '#f7f5f3',
          100: '#ece7e1',
          200: '#d9d1c7',
          300: '#bfb3a5',
          400: '#a09486',
          500: '#7d7264',
          600: '#5e554a',
          700: '#4a4339',
          800: '#2b2722',
          900: '#1a1714',
          950: '#0c0a08'
        },
        // 深色模式背景 - warm near-black surface scale (dark-theme surfaces)
        dark: {
          50: '#f7f5f3',
          100: '#e8e3dd',
          200: '#c9c0b5',
          300: '#9a8f80',
          400: '#6f655b',
          500: '#4d4947',
          600: '#3d3a39',
          700: '#2e2c2b',
          800: '#1a1816',
          900: '#0a0908',
          950: '#020202'
        }
      },
      fontFamily: {
        sans: [
          'Geist Variable',
          'Geist Fallback',
          'ui-sans-serif',
          'system-ui',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'PingFang SC',
          'Hiragino Sans GB',
          'Microsoft YaHei',
          'sans-serif'
        ],
        mono: [
          'Geist Mono Variable',
          'Geist Mono Fallback',
          'ui-monospace',
          'SFMono-Regular',
          'Menlo',
          'Monaco',
          'Consolas',
          'monospace'
        ]
      },
      // Keep NAMES, set to 'none' — elevation = hairline borders + tone separation
      boxShadow: {
        glass: 'none',
        'glass-sm': 'none',
        glow: 'none',
        'glow-lg': 'none',
        card: 'none',
        'card-hover': 'none',
        'inner-glow': 'none'
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-primary': 'linear-gradient(135deg, #ef6f2e 0%, #d15010 100%)',
        'gradient-dark': 'linear-gradient(135deg, #1a1714 0%, #020202 100%)',
        'gradient-glass':
          'linear-gradient(135deg, rgba(255,255,255,0.06) 0%, rgba(255,255,255,0.02) 100%)',
        'mesh-gradient':
          'radial-gradient(at 40% 20%, rgba(239,111,46,0.10) 0px, transparent 50%), radial-gradient(at 80% 0%, rgba(209,80,16,0.06) 0px, transparent 50%), radial-gradient(at 0% 50%, rgba(239,111,46,0.06) 0px, transparent 50%)'
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'slide-in-right': 'slideInRight 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        shimmer: 'shimmer 2s linear infinite',
        glow: 'glow 2s ease-in-out infinite alternate'
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' }
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideInRight: {
          '0%': { opacity: '0', transform: 'translateX(20px)' },
          '100%': { opacity: '1', transform: 'translateX(0)' }
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.95)' },
          '100%': { opacity: '1', transform: 'scale(1)' }
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' }
        },
        glow: {
          '0%': { boxShadow: '0 0 20px rgba(239, 111, 46, 0.25)' },
          '100%': { boxShadow: '0 0 30px rgba(239, 111, 46, 0.4)' }
        }
      },
      backdropBlur: {
        xs: '2px'
      },
      borderRadius: {
        '4xl': '2rem'
      }
    }
  },
  plugins: []
}
