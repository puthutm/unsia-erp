/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  reactStrictMode: true,
  transpilePackages: [],
  env: {
    NEXT_PUBLIC_API_CORE_URL: process.env.NEXT_PUBLIC_API_CORE_URL || 'http://unsia-core-service:8001',
    NEXT_PUBLIC_API_REFERENCE_URL: process.env.NEXT_PUBLIC_API_REFERENCE_URL || 'http://unsia-reference-service:8002',
    NEXT_PUBLIC_API_HRIS_URL: process.env.NEXT_PUBLIC_API_HRIS_URL || 'http://unsia-hris-service:8008',
    NEXT_PUBLIC_API_FINANCE_URL: process.env.NEXT_PUBLIC_API_FINANCE_URL || 'http://unsia-finance-service:8005',
  },
  experimental: {
    webpackBuildWorker: false
  }
}

module.exports = nextConfig
