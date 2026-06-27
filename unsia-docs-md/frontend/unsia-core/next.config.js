/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  reactStrictMode: true,
  transpilePackages: ['@phosphor-icons/react'],
  images: {
    domains: ['localhost'],
  },
  env: {
    NEXT_PUBLIC_API_CORE_URL: process.env.NEXT_PUBLIC_API_CORE_URL,
    NEXT_PUBLIC_API_REFERENCE_URL: process.env.NEXT_PUBLIC_API_REFERENCE_URL,
  },
  experimental: {
    webpackBuildWorker: false
  }
}

module.exports = nextConfig
