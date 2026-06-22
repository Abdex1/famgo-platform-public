
/*
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here 
};

export default nextConfig;
*/
import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  output: 'standalone',

  transpilePackages: [
    '@famgo/api-client',
    '@famgo/types',
    '@famgo/ui-theme',
    '@famgo/i18n',
    '@famgo/utils',
  ],

  experimental: {
    optimizePackageImports: [
      '@famgo/ui-theme',
      '@famgo/utils',
    ],
  },
};

export default nextConfig;